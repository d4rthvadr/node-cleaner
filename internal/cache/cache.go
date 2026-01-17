package cache

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/d4rthvadr/node-cleaner/pkg/models"
)

type Cache struct {
	index    *models.CacheIndex
	path     string
	mu       sync.RWMutex
	modified bool
}

func NewCache(cachePath string) (*Cache, error) {

	c := &Cache{
		path: cachePath,
		index: &models.CacheIndex{
			Version: "1.0",
			Entries: make(map[string]models.CacheEntry),
		},
	}

	// load existing cache if available
	// if the cache file does not exist, it's not an error
	if err := c.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return c, nil
}

// load reads the cache from the specified file
func (c *Cache) load() error {
	f, err := os.Open(c.path)
	if err != nil {
		return err

	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&c.index)

}

// Get retrieves a cache entry for the given path
func (c *Cache) Get(path string) (*models.CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.index.Entries[path]
	if !exists {
		return nil, false
	}
	return &entry, true
}

// Set adds or updates a cache entry for the given path
func (c *Cache) Set(path string, entry *models.CacheEntry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.index.Entries[path] = *entry
	c.modified = true
	return nil
}
