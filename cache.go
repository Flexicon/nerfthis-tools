package main

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

// Cache wrapper.
type Cache struct {
	cache *cache.Cache
}

// Get retrieves an item from cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

// Set adds an item to the cache, replacing any existing items for the given key.
func (c *Cache) Set(key string, val interface{}, ttl ...time.Duration) {
	d := cache.DefaultExpiration
	if len(ttl) >= 1 {
		d = ttl[0]
	}

	c.cache.Set(key, val, d)
}

// AppCache instance.
var AppCache *Cache

// CacheInit sets up the app cache.
func CacheInit() error {
	cacheTTL := time.Duration(viper.GetInt64("cache.ttl")) * time.Second
	AppCache = &Cache{cache.New(cacheTTL, 10*time.Minute)}
	return nil
}
