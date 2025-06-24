package cache

import (
	"sync"
	"time"
)

type item struct {
	value      []byte
	expiration int64
}

type Cache struct {
	items map[string]item
	mu    sync.RWMutex
}

func New() *Cache {
	cache := &Cache{
		items: make(map[string]item),
	}
	go cache.cleanupLoop()
	return cache
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:      value,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.expiration {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		for k, v := range c.items {
			if time.Now().UnixNano() > v.expiration {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
