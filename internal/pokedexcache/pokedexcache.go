package pokedexcache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

type Cache struct {
	Map   map[string]cacheEntry
	Mutex sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	cache := &Cache{
		Map: map[string]cacheEntry{},
	}
	go cache.readLoop(interval, ticker.C)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	newEntry := cacheEntry{
		createdAt: time.Now(),
		data:      val,
	}

	if entry, ok := c.Map[key]; ok {
		fmt.Printf("Value with key %s already exists: %q", key, entry.data)
		return
	}

	c.Map[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if entry, ok := c.Map[key]; ok {
		return entry.data, true
	}

	return []byte{}, false
}

func (c *Cache) readLoop(interval time.Duration, tickChan <-chan time.Time) {
	for range tickChan {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		for key, val := range c.Map {
			if time.Now().After(val.createdAt.Add(interval)) {
				delete(c.Map, key)
			}
		}
	}
}
