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

func isJSONResponse(contentType string, body []byte) bool {
	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if ct != "application/json" && ct != "text/json" {
		return false
	}
	body = bytes.TrimSpace(body)
	return len(body) > 0 && (body[0] == '{' || body[0] == '[')
}

const (
	socketTimeout             time.Duration = time.Second * 30
	maxIdleConnections        int           = 100
	maxIdleConnectionsPerHost int           = 100
	maxConnectionsPerHost     int           = 100
	idleConnTimeout           time.Duration = time.Second * 90
	netDialerTimeout          time.Duration = time.Second * 10
	netDialerKeepALive        time.Duration = time.Second * 10
	tlsHandshakeTimeout       time.Duration = time.Second * 5
	expectContinueTimeout     time.Duration = time.Second * 1
	forceHTTP2                bool          = true
	userAgent                 string        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	secUserAgent              string        = "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\""
	platform                  string        = "Windows"
	deviceIDLength            int           = 32
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
	UserAgent     string          `json:"user_agent,omitempty"`
	SecUserAgent  string          `json:"sec_user_agent,omitempty"`
	Platform      string          `json:"platform,omitempty"`
	DeviceID      string          `json:"device_id,omitempty"`
	RetryCount    int             `json:"retry_count,omitempty"`
	RetryInterval time.Duration   `json:"retry_interval,omitempty"`
	Transport     *http.Transport `json:"transport,omitempty"`
}

type BaseHTTPClient struct {
	http.Client
	mtx           sync.RWMutex
	userAgent     string
	secUserAgent  string
	platform      string
	deviceID      string
	retryCount    int
	retryInterval time.Duration
	isClose       bool
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
		userAgent:     userAgent,
		secUserAgent:  secUserAgent,
		platform:      platform,
		deviceID:      deviceID,
		retryCount:    3,
		retryInterval: 500 * time.Millisecond,
	}
}

func NewBaseHTTPClientWithParams(p *HTTPClientParameters) *BaseHTTPClient {
	client := NewBaseHTTPClient()

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

	if p.RetryCount > 0 {
		client.retryCount = p.RetryCount
	}

	if p.RetryInterval > 0 {
		client.retryInterval = p.RetryInterval
	}

	return client
}

func (c *BaseHTTPClient) request(ctx context.Context, req *http.Request) (*http.Response, error) {
	c.mtx.RLock()
	if c.isClose {
		c.mtx.RUnlock()
		return nil, fmt.Errorf("Transport.HttpClient.request(): client is already closed")
	}
	retryCount := c.retryCount
	retryInterval := c.retryInterval
	c.mtx.RUnlock()

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Sec-Ch-Ua", c.secUserAgent)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", c.platform)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("X-Device-Id", c.deviceID)

	var res *http.Response
	var err error

	for attempt := 0; attempt <= retryCount; attempt++ {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("Transport.HttpClient.request(): context canceled: %w", ctx.Err())
		}

		res, err = c.Do(req)
		if err != nil {
			// сетевые ошибки — не повторяем
			return nil, err
		}

		if res.StatusCode != http.StatusInternalServerError && res.StatusCode != http.StatusServiceUnavailable {
			break
		}

		if attempt == retryCount {
			break
		}

		res.Body.Close()

		select {
		case <-time.After(retryInterval):
		case <-ctx.Done():
			return nil, fmt.Errorf("Transport.HttpClient.request(): context canceled during retry wait: %w", ctx.Err())
		}
	}

	if res.StatusCode == http.StatusNotFound {
		return res, fmt.Errorf("Transport.HttpClient.request(): 404 page not found")
	}
	if res.StatusCode == http.StatusUnauthorized {
		return res, fmt.Errorf("Transport.HttpClient.request(): 401 user unauthorized")
	}

	if res.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(res.Body)
		if err != nil {
			res.Body.Close()
			return nil, fmt.Errorf("Transport.HttpClient.request(): error GZip decompressing data: %w", err)
		}

		res.Body = &ReadCloser{
			Origin: res.Body,
			Encode: gzipReader,
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
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return res, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): error reading response body: %w", err)
	}

	if !isJSONResponse(res.Header.Get("Content-Type"), bodyBytes) {
		return res, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): response is not valid JSON. Content-Type: %s\nBody: %s",
			res.Header.Get("Content-Type"), string(bodyBytes))
	}

	if err = json.Unmarshal(bodyBytes, target); err != nil {
		return res, fmt.Errorf("Transport.HttpClient.GetDecodeJSON(): error decoding JSON: %w\nBody: %s", err, string(bodyBytes))
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
		return res, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): error reading response body: %w", err)
	}

	if !isJSONResponse(res.Header.Get("Content-Type"), bodyBytes) {
		return res, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): response is not valid JSON. Content-Type: %s\nBody: %s",
			res.Header.Get("Content-Type"), string(bodyBytes))
	}

	if err = json.Unmarshal(bodyBytes, target); err != nil {
		return res, fmt.Errorf("Transport.HttpClient.PostDecodeJSON(): error decoding JSON: %w\nBody: %s", err, string(bodyBytes))
	}

	return res, err
}

func (c *BaseHTTPClient) SetUserAgent(ua string) {
	if ua != "" {
		c.mtx.Lock()
		c.userAgent = ua
		c.mtx.Unlock()
	}
}

func (c *BaseHTTPClient) GetUserAgent() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.userAgent
}

func (c *BaseHTTPClient) SetPlatform(p string) {
	if p != "" {
		c.mtx.Lock()
		c.platform = p
		c.mtx.Unlock()
	}
}

func (c *BaseHTTPClient) GetPlatform() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.platform
}

// SetDeviceID length id == 64
func (c *BaseHTTPClient) SetDeviceID(id string) error {
	if len(id) != 64 {
		return fmt.Errorf("Transport.HttpClient.SetDeviceID(): device id should be 64 characters long")
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
	if c.isClose {
		return fmt.Errorf("Transport.HttpClient.Close(): client is already closed")
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()

	if transport, ok := c.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
		transport.DisableKeepAlives = true
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
	bytes := make([]byte, deviceIDLength)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("errors generating device ID: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
