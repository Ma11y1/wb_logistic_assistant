package request

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"wb_logistic_assistant/external/wb_logistic_api/transport"
)

type APIRequest[T any] interface {
	Do(ctx context.Context) (T, error)
	Get(ctx context.Context) (*http.Response, error)
	GetUnmarshal(ctx context.Context, target interface{}) error
	Post(ctx context.Context) (*http.Response, error)
	PostUnmarshal(ctx context.Context, target interface{}) error
	PostData(ctx context.Context, data io.Reader) (*http.Response, error)
	Params() *Parameters
	SetParams(params *Parameters)
	QueryParams() *Parameters
	SetQueryParams(params *Parameters)
	SetHeader(h http.Header)
	SetHeaderValue(key, value string)
	AddHeaderValue(key, value string)
	GetHeaderValue(key string) string
	DeleteHeader(key string)
	HeaderCopy() http.Header
	URL() string
	SetURL(url string)
	SetAccessToken(token string)
	SetContentType(t string)
	SetAcceptEncoding(a string)
	IsUnauthorized() bool
	ClearParameters()
}

type BaseRequest struct {
	mtx             sync.RWMutex
	url             string
	client          transport.HTTPClient
	header          http.Header
	queryParameters *Parameters
	parameters      *Parameters
	isUnauthorized  bool
}

func NewRequest(client transport.HTTPClient, url string) *BaseRequest {
	return &BaseRequest{
		mtx:             sync.RWMutex{},
		client:          client,
		url:             url,
		header:          http.Header{},
		queryParameters: NewParameters(),
		parameters:      NewParameters(),
	}
}

func NewRequestToken(client transport.HTTPClient, url string, token string) *BaseRequest {
	req := NewRequest(client, url)
	req.SetAccessToken(token)
	return req
}

func (r *BaseRequest) Do(context.Context) (interface{}, error) {
	return nil, nil
}

func (r *BaseRequest) Get(ctx context.Context) (*http.Response, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	u := r.buildFullURL()

	httpResponse, err := r.client.Get(ctx, u, r.header)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusUnauthorized {
			r.isUnauthorized = true
		}
		return nil, fmt.Errorf("BaseRequest.Get(): Error GET request %s: %s", u, err)
	}

	return httpResponse, nil
}

func (r *BaseRequest) GetUnmarshal(ctx context.Context, target interface{}) error {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	u := r.buildFullURL()

	httpResponse, err := r.client.GetDecodeJSON(ctx, u, target, r.header)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusUnauthorized {
			r.isUnauthorized = true
		}
		return fmt.Errorf("BaseRequest.GetUnmarshal(): Error GET and unmarshal JSON request %s: %s", u, err)
	}
	return nil
}

func (r *BaseRequest) Post(ctx context.Context) (*http.Response, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	u := r.buildFullURL()

	body, err := r.parameters.BuildJSON()
	if err != nil {
		return nil, fmt.Errorf("BaseRequest.Post(): Error building body request %s: %s", u, err)
	}

	httpResponse, err := r.client.Post(ctx, u, bytes.NewBuffer(body), r.header)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusUnauthorized {
			r.isUnauthorized = true
		}
		return nil, fmt.Errorf("BaseRequest.Post(): Error POST request %s: %s", u, err)
	}

	return httpResponse, nil
}

func (r *BaseRequest) PostUnmarshal(ctx context.Context, target interface{}) error {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	u := r.buildFullURL()

	body, err := r.parameters.BuildJSON()
	if err != nil {
		return fmt.Errorf("BaseRequest.PostUnmarshal(): Error building body request %s: %s", u, err)
	}

	httpResponse, err := r.client.PostDecodeJSON(
		ctx,
		u,
		bytes.NewBuffer(body),
		target,
		r.header,
	)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusUnauthorized {
			r.isUnauthorized = true
		}
		return fmt.Errorf("BaseRequest.PostUnmarshal(): Error POST and unmarshal JSON request %s: %s", u, err)
	}
	return nil
}

func (r *BaseRequest) PostData(ctx context.Context, data io.Reader) (*http.Response, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	u := r.buildFullURL()

	httpResponse, err := r.client.Post(ctx, u, data, r.header)
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusUnauthorized {
			r.isUnauthorized = true
		}
		return nil, fmt.Errorf("BaseRequest.PostData(): Error POST request %s: %s", u, err)
	}

	return httpResponse, nil
}

func (r *BaseRequest) buildFullURL() string {
	u := r.url
	queryString := r.queryParameters.BuildURLValuesEncode()
	if queryString != "" {
		if strings.Contains(u, "?") {
			u += "&" + queryString
		} else {
			u += "?" + queryString
		}
	}
	return u
}

func (r *BaseRequest) Params() *Parameters {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.parameters
}

func (r *BaseRequest) SetParams(params *Parameters) {
	if params == nil {
		return
	}
	r.mtx.Lock()
	r.parameters = params
	r.mtx.Unlock()
}

func (r *BaseRequest) QueryParams() *Parameters {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.queryParameters
}

func (r *BaseRequest) SetQueryParams(params *Parameters) {
	if params == nil {
		return
	}
	r.mtx.Lock()
	r.queryParameters = params
	r.mtx.Unlock()
}

func (r *BaseRequest) SetHeader(h http.Header) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.header = make(http.Header, len(h))
	for k, v := range h {
		r.header[k] = append([]string(nil), v...)
	}
}

func (r *BaseRequest) SetHeaderValue(key, value string) {
	r.mtx.Lock()
	r.header.Set(key, value)
	r.mtx.Unlock()
}

func (r *BaseRequest) AddHeaderValue(key, value string) {
	r.mtx.Lock()
	r.header.Add(key, value)
	r.mtx.Unlock()
}

func (r *BaseRequest) GetHeaderValue(key string) string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.header.Get(key)
}

func (r *BaseRequest) DeleteHeader(key string) {
	r.mtx.Lock()
	r.header.Del(key)
	r.mtx.Unlock()
}

func (r *BaseRequest) HeaderCopy() http.Header {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	copy := make(http.Header, len(r.header))
	for k, v := range r.header {
		vcopy := make([]string, len(v))
		copy[k] = append(vcopy, v...)
	}
	return copy
}

func (r *BaseRequest) URL() string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.url
}

func (r *BaseRequest) SetURL(url string) {
	if url != "" {
		r.mtx.Lock()
		r.url = url
		r.mtx.Unlock()
	}
}

func (r *BaseRequest) SetAccessToken(token string) {
	if token != "" {
		r.mtx.Lock()
		r.header.Set("Authorization", "Bearer "+token)
		r.mtx.Unlock()
	}
}

func (r *BaseRequest) SetContentType(t string) {
	if t != "" {
		r.mtx.Lock()
		r.header.Set("Content-Type", t)
		r.mtx.Unlock()
	}
}

func (r *BaseRequest) SetAcceptEncoding(a string) {
	if a != "" {
		r.mtx.Lock()
		r.header.Set("Accept-Encoding", a)
		r.mtx.Unlock()
	}
}

func (r *BaseRequest) IsUnauthorized() bool {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return r.isUnauthorized
}

func (r *BaseRequest) ClearParameters() {
	r.mtx.Lock()
	r.queryParameters.Clear()
	r.parameters.Clear()
	r.mtx.Unlock()
}
