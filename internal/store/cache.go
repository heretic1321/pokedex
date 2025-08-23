package store

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu    sync.Mutex
	cache map[string]CacheEntry
}

func New(interval time.Duration) *Cache {
	cache := Cache{
		mu:    sync.Mutex{},
		cache: make(map[string]CacheEntry),
	}

	go cache.readLoop(interval)
	return &cache
}

func (c *Cache) readLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for time := range ticker.C {
		
		c.mu.Lock()
		i := 0
		for name, cache := range c.cache {
			if time.Sub(cache.createdAt) > interval {
				delete(c.cache, name)
				i += 1
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Add(key string, val []byte) {

	cacheEntry := CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	c.cache[key] = cacheEntry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {


	c.mu.Lock()
	e, ok := c.cache[key]
	c.mu.Unlock()
	if !ok {
		return nil, false
	}

	return e.val, true

}
