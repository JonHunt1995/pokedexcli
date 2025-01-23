package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu    sync.Mutex
	cache map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		mu:    sync.Mutex{},
		cache: make(map[string]cacheEntry),
	}
	c.reapLoop(interval) // Purge the old caches
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ce := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.cache[key] = ce

}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.cache[key]
	if !ok {
		return []byte{}, false
	}
	return v.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.mu.Lock()
			for k, v := range c.cache {
				if time.Since(v.createdAt) > interval {
					delete(c.cache, k)
				}
			}
			c.mu.Unlock()
		}
	}()
}
