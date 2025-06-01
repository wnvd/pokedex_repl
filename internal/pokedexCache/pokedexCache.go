package pokedexCache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt 	time.Time
	val 		[]byte
}

type Cache struct {
	cachedPayload 	map[string]cacheEntry
	mu 				sync.Mutex
}

func (c *Cache) Add(key string, payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cachedPayload[key] = cacheEntry {
		createdAt: time.Now(),
		val: payload,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cachedPayload, present := c.cachedPayload[key]
	if !present {
		return nil, false
	}

	return cachedPayload.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.cachedPayload {
			if time.Since(val.createdAt) > interval {
				delete(c.cachedPayload, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache {
		cachedPayload: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return &cache
}
