package pokecache
import (
    "sync"
    "time"
)

func NewCache(interval time.Duration) *Cache {
    // initialize a cache and start its timeout loop concurrently
    cache := Cache{
        entries: map[string]cacheEntry{},
    }
    go cache.reapLoop(interval)
    return &cache
}

func (c *Cache) reapLoop(interval time.Duration) {
    // clears cache entries after they reach a certain age
    // ran in a goroutine necessitating mutexes for read/write safety
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
    // add a url/response pair to the cache with mutex safety
    c.mu.Lock()
    defer c.mu.Unlock()
    c.entries[key] = cacheEntry{
        val: val,
        createdAt: time.Now(),
    }
}

func (c *Cache) Get(key string) ([]byte, bool) {
    // check for and return the response and existence for a url in the cache
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
    // cache struct which, when initialized with NewCache(), will clear old entries
    // used for url/response pairs
    entries map[string]cacheEntry
    mu sync.Mutex
}

type cacheEntry struct {
    createdAt time.Time
    val []byte
}
