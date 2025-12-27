package storage

import "sync"

type FileCache struct {
	mtx  sync.RWMutex
	data map[string]string
}

func NewFileCache() *FileCache {
	return &FileCache{
		data: make(map[string]string),
	}
}

func (c *FileCache) Set(name, data string) {
	if name == "" {
		return
	}
	c.mtx.Lock()
	c.data[name] = data
	c.mtx.Unlock()
}

func (c *FileCache) Get(name string) string {
	if name == "" {
		return ""
	}
	c.mtx.RLock()
	data := c.data[name]
	c.mtx.RUnlock()
	return data
}

func (c *FileCache) Remove(name string) {
	if name == "" {
		return
	}
	c.mtx.Lock()
	delete(c.data, name)
	c.mtx.Unlock()
}

func (c *FileCache) Has(name string) bool {
	if name == "" {
		return false
	}
	c.mtx.RLock()
	_, ok := c.data[name]
	c.mtx.RUnlock()
	return ok
}

func (c *FileCache) Clear() {
	c.mtx.Lock()
	c.data = make(map[string]string)
	c.mtx.Unlock()
}
