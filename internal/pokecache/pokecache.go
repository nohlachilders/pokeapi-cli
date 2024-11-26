package pokecache
import (
    "sync"
    "time"
)

func NewCache(interval time.Duration) *Cache {
    cache := Cache{
        entries: map[string]cacheEntry{},
    }
    go cache.reapLoop(interval)
    return &cache
}

func (c *Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    for t := range ticker.C {
        c.mu.Lock()
        for key, value := range c.entries {
            if t.Sub(value.createdAt) > interval {
                delete(c.entries, key)
            }
        }
        c.mu.Unlock()
    }
}

func (c *Cache) Add(key string, val []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.entries[key] = cacheEntry{
        val: val,
        createdAt: time.Now(),
    }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    entry, ok := c.entries[key]
    switch ok{
    case true:
        return entry.val, ok
    case false:
        return []byte{}, ok
    }
    return []byte{}, ok
}

type Cache struct {
    entries map[string]cacheEntry
    mu sync.Mutex
}

type cacheEntry struct {
    createdAt time.Time
    val []byte
}
