package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

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
// Any modification marks the cache as dirty but stored in memory
// until Save() is called
func (c *Cache) Set(path string, entry *models.CacheEntry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.index.Entries[path] = *entry
	c.modified = true
	return nil
}

// IsValid checks if the cache entry for the given path is still valid
func (c *Cache) IsValid(path string, currentModTime time.Time) bool {

	entry, exists := c.Get(path)
	if !exists {
		return false
	}

	return entry.ModTime.Equal(currentModTime)

}

// Save writes data from memory cache to the disk file
// if there are modifications
func (c *Cache) Save() error {

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.modified {
		return nil // no changes to save
	}

	c.index.UpdatedAt = time.Now()

	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// write to a temp file first for atomicity
	tempPath := c.path + ".tmp"
	f, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	err = encoder.Encode(c.index)

	f.Close()

	if err != nil {
		return err
	}

	// replace old cache file with the new one
	return os.Rename(tempPath, c.path)
}

// Clear removes all entries from the cache
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.index.Entries = make(map[string]models.CacheEntry)
	c.modified = true

	return c.Save()
}
