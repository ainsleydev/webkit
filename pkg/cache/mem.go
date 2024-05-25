package cache

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"
)

type (
	// MemCache represents an in-memory cache store.
	MemCache struct {
		cache             map[string]inMemCacheItem
		defaultExpiration time.Duration
		mutex             *sync.RWMutex
	}
	// inMemCacheItem represents an item stored in the cache.
	inMemCacheItem struct {
		value      any
		expiration time.Time
		tags       []string
		noExpiry   bool
	}
)

// NewInMemory creates a new instance of MemCache.
func NewInMemory(defaultExpiration time.Duration) *MemCache {
	return &MemCache{
		cache:             make(map[string]inMemCacheItem),
		defaultExpiration: defaultExpiration,
		mutex:             &sync.RWMutex{},
	}
}

func (c *MemCache) Ping(_ context.Context) error {
	// In-memory cache doesn't need to be pinged
	return nil
}

func (c *MemCache) Get(_ context.Context, key string, value interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.cache[key]
	if !found {
		return errors.New("key not found")
	}
	if !item.noExpiry && time.Now().After(item.expiration) {
		// Item has expired, delete it from the cache
		delete(c.cache, key)
		return errors.New("key expired")
	}

	// Check if value is a pointer
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return errors.New("value must be a pointer")
	}

	// Copy value to the provided interface
	reflect.ValueOf(value).Elem().Set(reflect.ValueOf(item.value))

	return nil
}

func (c *MemCache) Set(ctx context.Context, key string, value interface{}, opts Options) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if opts.Expiration == 0 {
		opts.Expiration = c.defaultExpiration
	}
	c.cache[key] = inMemCacheItem{
		value:      value,
		expiration: time.Now().Add(opts.Expiration),
		tags:       opts.Tags,
		noExpiry:   opts.Expiration == Forever,
	}
	return nil
}

func (c *MemCache) Delete(_ context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, key)
	return nil
}

func (c *MemCache) Invalidate(_ context.Context, tags []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, item := range c.cache {
		for _, tag := range tags {
			for _, itemTag := range item.tags {
				if tag == itemTag {
					delete(c.cache, key)
					break
				}
			}
		}
	}
}

func (c *MemCache) Flush(_ context.Context) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache = make(map[string]inMemCacheItem)
}

func (c *MemCache) Close() error {
	// No action needed for an in-memory cache store
	return nil
}
