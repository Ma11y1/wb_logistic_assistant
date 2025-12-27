package services

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Cache[T any] struct {
	mtx       sync.RWMutex
	value     T
	expiresAt time.Time
	ttl       time.Duration
	loader    func(ctx context.Context) (T, error)
}

func NewCache[T any](ttl time.Duration, loader func(ctx context.Context) (T, error)) *Cache[T] {
	return &Cache[T]{ttl: ttl, loader: loader}
}

func (c *Cache[T]) Get(ctx context.Context) (T, error) {
	c.mtx.RLock()
	if time.Now().Before(c.expiresAt) {
		val := c.value
		c.mtx.RUnlock()
		return val, nil
	}
	c.mtx.RUnlock()

	// load new
	newVal, err := c.loader(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	c.mtx.Lock()
	c.value = newVal
	c.expiresAt = time.Now().Add(c.ttl)
	c.mtx.Unlock()

	return newVal, nil
}

func (c *Cache[T]) Invalidate() {
	c.mtx.Lock()
	c.expiresAt = time.Now()
	c.mtx.Unlock()
}

type cachedItem[V any] struct {
	value     V
	expiresAt time.Time
}

type GenericMapCache[K comparable, V any] struct {
	mtx            sync.RWMutex
	data           map[K]cachedItem[V]
	ttl            time.Duration
	loader         func(ctx context.Context, key K) (V, error)
	cleanupCounter int32
}

func NewGenericMapCache[K comparable, V any](
	ttl time.Duration,
	loader func(ctx context.Context, key K) (V, error),
) *GenericMapCache[K, V] {
	return &GenericMapCache[K, V]{
		data:   make(map[K]cachedItem[V]),
		ttl:    ttl,
		loader: loader,
	}
}

func (c *GenericMapCache[K, V]) Get(ctx context.Context, key K) (V, error) {
	c.cleanup()

	c.mtx.RLock()
	item, ok := c.data[key]
	if ok && time.Now().Before(item.expiresAt) {
		val := item.value
		c.mtx.RUnlock()
		return val, nil
	}
	c.mtx.RUnlock()

	newVal, err := c.loader(ctx, key)
	if err != nil {
		var zero V
		return zero, err
	}

	c.mtx.Lock()
	c.data[key] = cachedItem[V]{
		value:     newVal,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mtx.Unlock()

	return newVal, nil
}

func (c *GenericMapCache[K, V]) Invalidate(key K) {
	c.mtx.Lock()
	delete(c.data, key)
	c.mtx.Unlock()
}

func (c *GenericMapCache[K, V]) InvalidateAll() {
	c.mtx.Lock()
	c.data = make(map[K]cachedItem[V])
	c.mtx.Unlock()
}

func (c *GenericMapCache[K, V]) cleanup() {
	if atomic.AddInt32(&c.cleanupCounter, 1)%5000 == 0 {
		atomic.StoreInt32(&c.cleanupCounter, 0)
		now := time.Now()
		c.mtx.Lock()
		for k, item := range c.data {
			if now.After(item.expiresAt) {
				delete(c.data, k)
			}
		}
		c.mtx.Unlock()
	}
}
