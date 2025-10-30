package request

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"sync"
)

type Parameters struct {
	mu sync.RWMutex
	m  map[string]interface{}
}

func NewParameters() *Parameters {
	return &Parameters{
		m: make(map[string]interface{}),
	}
}

func (p *Parameters) Set(key string, value interface{}) error {
	if key == "" {
		return errors.New("empty key")
	}
	p.mu.Lock()
	p.m[key] = value
	p.mu.Unlock()
	return nil
}

func (p *Parameters) SetIfNotEmpty(key string, value interface{}) error {
	if value != "" && value != nil {
		return p.Set(key, value)
	}
	return nil
}

func (p *Parameters) SetMap(m map[string]interface{}) error {
	if m == nil {
		return errors.New("empty map")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for k, v := range m {
		p.m[k] = v
	}
	return nil
}

// SetArr [key, value, key, value...]
func (p *Parameters) SetArr(pairs []interface{}) error {
	if len(pairs) == 0 || len(pairs)%2 != 0 {
		return errors.New("pairs length must be even")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := 0; i < len(pairs)-1; i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return errors.New("key must be a string")
		}
		p.m[key] = pairs[i+1]
	}
	return nil
}

func (p *Parameters) Get(key string) interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.m[key]
}

func (p *Parameters) GetString(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.getStringUnsafe(key)
}

func (p *Parameters) Remove(key string) {
	p.mu.Lock()
	delete(p.m, key)
	p.mu.Unlock()
}

func (p *Parameters) Has(key string) bool {
	p.mu.RLock()
	_, ok := p.m[key]
	p.mu.RUnlock()
	return ok
}

func (p *Parameters) Clear() {
	p.mu.Lock()
	p.m = make(map[string]interface{})
	p.mu.Unlock()
}

func (p *Parameters) BuildURLValues() *url.Values {
	values := &url.Values{}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for key, value := range p.m {
		if value != "" && value != nil {
			values.Add(key, p.getStringUnsafe(key))
		}
	}
	return values
}

func (p *Parameters) BuildURLValuesEncode() string {
	return p.BuildURLValues().Encode()
}

func (p *Parameters) BuildJSON() ([]byte, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return json.Marshal(p.m)
}

func (p *Parameters) getStringUnsafe(key string) string {
	if val, ok := p.m[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case int:
			return strconv.Itoa(v)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(v)
		default:
			return ""
		}
	}
	return ""
}
