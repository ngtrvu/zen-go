package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	Timestamp time.Time
}

type Cache struct {
	item  CacheItem
	mutex sync.RWMutex
	ttl   time.Duration
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		ttl: ttl,
	}
}

func (c *Cache) Get() (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if time.Since(c.item.Timestamp) < c.ttl {
		return c.item.Data, true
	}
	return nil, false
}

func (c *Cache) Set(data interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.item = CacheItem{
		Data:      data,
		Timestamp: time.Now(),
	}
}
