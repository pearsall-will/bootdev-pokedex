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
	mu      *sync.Mutex
	entries map[string]cacheEntry
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.entries[key] = cacheEntry{time.Now(), val}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, exists := cache.entries[key]
	if exists {
		return entry.val, true
	}
	return []byte{}, false
}

func (cache *Cache) Remove(key string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	delete(cache.entries, key)
}

func (cache *Cache) reapLoop(interval time.Duration) {
	for {
		deleteTheshold := time.Now()
		time.Sleep(interval)
		for key, entry := range cache.entries {
			needsDelete := entry.createdAt.Before(deleteTheshold)
			if needsDelete {
				cache.Remove(key)

			}
		}
	}
}

func NewCache(inteval time.Duration) Cache {
	cache := Cache{}
	go cache.reapLoop(inteval)
	return cache
}
