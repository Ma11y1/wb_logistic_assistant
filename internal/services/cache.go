package services

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

type Cache[T interface{}] struct {
	mtx       sync.RWMutex
	value     T
	expiresAt time.Time
	ttl       time.Duration
	loader    func(ctx context.Context) (T, error)
	group     singleflight.Group
}

func NewCache[T interface{}](ttl time.Duration, loader func(ctx context.Context) (T, error)) *Cache[T] {
	return &Cache[T]{ttl: ttl, loader: loader}
}

func (c *Cache[T]) Get(ctx context.Context) (T, error) {
	if ctx.Err() != nil {
		var zero T
		return zero, ctx.Err()
	}

	c.mtx.RLock()
	if time.Now().Before(c.expiresAt) {
		val := c.value
		c.mtx.RUnlock()
		return val, nil
	}
	c.mtx.RUnlock()

	val, err, _ := c.group.Do("load", func() (interface{}, error) {
		c.mtx.RLock()
		if time.Now().Before(c.expiresAt) {
			v := c.value
			c.mtx.RUnlock()
			return v, nil
		}
		c.mtx.RUnlock()

		newVal, err := c.loader(context.Background())
		if err != nil {
			return nil, err
		}

		c.mtx.Lock()
		c.value = newVal
		c.expiresAt = time.Now().Add(c.ttl)
		c.mtx.Unlock()
		return newVal, nil
	})

	if err != nil {
		var zero T
		return zero, err
	}
	return val.(T), nil
}

func (c *Cache[T]) Invalidate() {
	c.mtx.Lock()
	c.expiresAt = time.Time{}
	c.mtx.Unlock()
}

type cachedItem[V interface{}] struct {
	value     V
	expiresAt time.Time
}

type GenericMapCache[K comparable, V interface{}] struct {
	mtx            sync.RWMutex
	data           map[K]cachedItem[V]
	ttl            time.Duration
	loader         func(ctx context.Context, key K) (V, error)
	group          singleflight.Group
	cleanupCounter int32
}

func NewGenericMapCache[K comparable, V interface{}](
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
	if ctx.Err() != nil {
		var zero V
		return zero, ctx.Err()
	}

	c.cleanup()

	c.mtx.RLock()
	item, ok := c.data[key]
	if ok && time.Now().Before(item.expiresAt) {
		val := item.value
		c.mtx.RUnlock()
		return val, nil
	}
	c.mtx.RUnlock()

	stringKey := c.keyToString(key)
	val, err, _ := c.group.Do(stringKey, func() (interface{}, error) {
		c.mtx.RLock()
		if it, exists := c.data[key]; exists && time.Now().Before(it.expiresAt) {
			v := it.value
			c.mtx.RUnlock()
			return v, nil
		}
		c.mtx.RUnlock()

		newVal, err := c.loader(ctx, key)
		if err != nil {
			return nil, err
		}

		c.mtx.Lock()
		c.data[key] = cachedItem[V]{
			value:     newVal,
			expiresAt: time.Now().Add(c.ttl),
		}
		c.mtx.Unlock()
		return newVal, nil
	})

	if err != nil {
		var zero V
		return zero, err
	}
	return val.(V), nil
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
	count := atomic.AddInt32(&c.cleanupCounter, 1)

	if count >= 1000 {
		if atomic.CompareAndSwapInt32(&c.cleanupCounter, count, 0) {
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
}

func (c *GenericMapCache[K, V]) keyToString(key K) string {
	switch v := any(key).(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	default:
		// Для остальных типов (float, struct и т.д.) оставляем Sprint
		return fmt.Sprint(key)
	}
}
