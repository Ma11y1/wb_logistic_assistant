package transport

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

func isLikelyJSON(contentType string, body []byte) bool {
	if contentType == "application/json" {
		return true
	}
	if len(body) == 0 {
		return false
	}
	b := bytes.TrimSpace(body)
	if len(b) == 0 {
		return false
	}
	c := b[0]
	return c == '{' || c == '['
}

const (
	socketTimeout             = 30 * time.Second
	maxIdleConnections        = 150
	maxIdleConnectionsPerHost = 150
	maxConnectionsPerHost     = 150
	idleConnTimeout           = 120 * time.Second
	netDialerTimeout          = 10 * time.Second
	netDialerKeepALive        = 10 * time.Second
	tlsHandshakeTimeout       = 5 * time.Second
	expectContinueTimeout     = 1 * time.Second
	forceHTTP2                = true

	userAgent    string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	secUserAgent string = "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\""
	platform     string = "Windows"

	deviceIDLength int = 32
)

type HTTPClient interface {
	Get(ctx context.Context, url string, header http.Header) (*http.Response, error)
	GetDecodeJSON(ctx context.Context, url string, target interface{}, header http.Header) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, header http.Header) (*http.Response, error)
	PostDecodeJSON(ctx context.Context, url string, body io.Reader, target interface{}, header http.Header) (*http.Response, error)

	SetUserAgent(ua string)
	GetUserAgent() string

	SetPlatform(p string)
	GetPlatform() string

	SetDeviceID(id string) error
	GetDeviceID() string

	Close() error
	IsClose() bool
}

type HTTPClientParameters struct {
	UserAgent    string          `json:"user_agent,omitempty"`
	SecUserAgent string          `json:"sec_user_agent,omitempty"`
	Platform     string          `json:"platform,omitempty"`
	DeviceID     string          `json:"device_id,omitempty"`
	Transport    *http.Transport `json:"transport,omitempty"`
}

type BaseHTTPClient struct {
	http.Client

	mtx          sync.RWMutex
	userAgent    string
	secUserAgent string
	platform     string
	deviceID     string
	isClose      bool

	defaultHeaders http.Header
}

func NewBaseHTTPClient() *BaseHTTPClient {
	jar, _ := cookiejar.New(nil)

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   netDialerTimeout,
			KeepAlive: netDialerKeepALive,
		}).DialContext,
		MaxIdleConns:          maxIdleConnections,
		MaxConnsPerHost:       maxConnectionsPerHost,
		MaxIdleConnsPerHost:   maxIdleConnectionsPerHost,
		IdleConnTimeout:       idleConnTimeout,
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ExpectContinueTimeout: expectContinueTimeout,
		ForceAttemptHTTP2:     forceHTTP2,
	}

	deviceID, _ := generateDeviceID()

	return &BaseHTTPClient{
		Client: http.Client{
			Jar:       jar,
			Timeout:   socketTimeout,
			Transport: transport,
		},
		userAgent:    userAgent,
		secUserAgent: secUserAgent,
		platform:     platform,
		deviceID:     deviceID,
	}
}

func NewBaseHTTPClientWithParams(p *HTTPClientParameters) *BaseHTTPClient {
	client := NewBaseHTTPClient()

	if p == nil {
		return client
	}

	if p.Transport != nil {
		client.Transport = p.Transport
	}

	if p.UserAgent != "" {
		client.userAgent = p.UserAgent
	}

	if p.SecUserAgent != "" {
		client.secUserAgent = p.SecUserAgent
	}

	if p.Platform != "" {
		client.platform = p.Platform
	}

	if p.DeviceID != "" && len(p.DeviceID) == 64 {
		client.deviceID = p.DeviceID
	}

	return client
}

func (c *BaseHTTPClient) request(ctx context.Context, req *http.Request) (*http.Response, error) {
	c.mtx.RLock()
	if c.isClose {
		c.mtx.RUnlock()
		return nil, fmt.Errorf("Transport.HttpClient.request(): client is closed")
	}
	ua := c.userAgent
	secUA := c.secUserAgent
	plat := c.platform
	deviceID := c.deviceID
	c.mtx.RUnlock()

	if req.Header == nil {
		req.Header = make(http.Header)
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", ua)
	}
	if req.Header.Get("Sec-Ch-Ua") == "" {
		req.Header.Set("Sec-Ch-Ua", secUA)
	}
	if req.Header.Get("Sec-Ch-Ua-Mobile") == "" {
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	}
	if req.Header.Get("Sec-Ch-Ua-Platform") == "" {
		req.Header.Set("Sec-Ch-Ua-Platform", plat)
	}
	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json, text/plain, */*")
	}
	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")
	}
	if req.Header.Get("X-Device-Id") == "" {
		req.Header.Set("X-Device-Id", deviceID)
	}
	if req.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if ctx.Err() != nil {
		return nil, fmt.Errorf("Transport.HttpClient.request(): context canceled: %v", ctx.Err())
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Transport.HttpClient.request(): %v", err)
	}

	if strings.EqualFold(res.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(res.Body)
		if err != nil {
			res.Body.Close()
			return nil, fmt.Errorf("Transport.HttpClient.request(): GZIP error: %v", err)
		}
		res.Body = &ReadCloser{
			Origin: res.Body,
			Encode: gz,
		}
	}

	return res, nil
}

func (c *BaseHTTPClient) Get(ctx context.Context, url string, header http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Transport.HttpClient.Get(): error create HTTP request with context %s", err)
	}

	if header != nil {
		req.Header = header
	}

	return c.request(ctx, req)
}

func (c *BaseHTTPClient) GetDecodeJSON(ctx context.Context, url string, target interface{}, header http.Header) (*http.Response, error) {
	if target == nil {
		return nil, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): target cannot be nil")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): error create HTTP request with context %s", err)
	}

	if header != nil {
		req.Header = header
	}

	res, err := c.request(ctx, req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return res, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): error reading response body: %v", err)
	}

	if !isLikelyJSON(res.Header.Get("Content-Type"), body) {
		return res, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): response is not valid JSON. Content-Type: %s\nBody: %s",
			res.Header.Get("Content-Type"), string(body))
	}

	if err = json.Unmarshal(body, target); err != nil {
		return res, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): error decoding JSON: %v\nBody: %s", err, string(body))
	}
	return res, err
}

func (c *BaseHTTPClient) Post(ctx context.Context, url string, body io.Reader, header http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("Transport.HttpClient.Post(): error create HTTP request with context %s", err)
	}

	if header != nil {
		req.Header = header
	}

	return c.request(ctx, req)
}

func (c *BaseHTTPClient) PostDecodeJSON(ctx context.Context, url string, body io.Reader, target interface{}, header http.Header) (*http.Response, error) {
	if target == nil {
		return nil, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): target cannot be nil")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): error create HTTP request with context %s", err)
	}

	if header != nil {
		req.Header = header
	}

	res, err := c.request(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return res, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): error reading response body: %v", err)
	}

	if !isLikelyJSON(res.Header.Get("Content-Type"), bodyBytes) {
		return res, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): response is not valid JSON. Content-Type: %s\nBody: %s",
			res.Header.Get("Content-Type"), string(bodyBytes))
	}

	if err = json.Unmarshal(bodyBytes, target); err != nil {
		return res, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): error decoding JSON: %v\nBody: %s", err, string(bodyBytes))
	}

	return res, err
}

func (c *BaseHTTPClient) SetUserAgent(ua string) {
	c.mtx.Lock()
	c.userAgent = ua
	c.mtx.Unlock()
}

func (c *BaseHTTPClient) GetUserAgent() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.userAgent
}

func (c *BaseHTTPClient) SetPlatform(p string) {
	c.mtx.Lock()
	c.platform = p
	c.mtx.Unlock()
}

func (c *BaseHTTPClient) GetPlatform() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.platform
}

func (c *BaseHTTPClient) SetDeviceID(id string) error {
	if len(id) <= 0 || len(id) > deviceIDLength {
		return fmt.Errorf("device id must be > 0 and < %d chars", deviceIDLength)
	}
	c.mtx.Lock()
	c.deviceID = id
	c.mtx.Unlock()
	return nil
}

func (c *BaseHTTPClient) GetDeviceID() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.deviceID
}

func (c *BaseHTTPClient) Close() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if c.isClose {
		return fmt.Errorf("client already closed")
	}

	if tr, ok := c.Transport.(*http.Transport); ok {
		tr.CloseIdleConnections()
		tr.DisableKeepAlives = true
	}

	c.isClose = true
	return nil
}

func (c *BaseHTTPClient) IsClose() bool {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.isClose
}
func generateDeviceID() (string, error) {
	b := make([]byte, deviceIDLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("errors generating device ID: %v", err)
	}
	return hex.EncodeToString(b), nil
}
