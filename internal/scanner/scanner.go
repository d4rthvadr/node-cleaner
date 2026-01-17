package scanner

import (
	"context"
	"sync"
	"time"

	"github.com/d4rthvadr/node-cleaner/pkg/models"
)

type Scanner struct {
	config  *models.Config
	cache   CacheProvider
	results chan models.DependencyFolder
	errors  chan error
}

type CacheProvider interface {
	Get(path string) (*models.CacheEntry, bool)
	Set(path string, entry *models.CacheEntry) error
}

func NewScanner(cfg *models.Config, cache CacheProvider) *Scanner {
	return &Scanner{
		config:  cfg,
		cache:   cache,
		results: make(chan models.DependencyFolder),
		errors:  make(chan error),
	}
}

// Scan initiates file traversal process
func (s *Scanner) Scan(ctx context.Context, rootPath string) (*models.ScanResult, error) {

	finalResult := &models.ScanResult{
		ScanPath: rootPath,
		ScanTime: time.Now(),
	}

	var wg sync.WaitGroup

	// start workers
	for i := 0; i < s.config.Workers; i++ {
		wg.Add(1)
		go s.worker(ctx, &wg)
	}

	// walk the file system starting from rootPath
	// and send directories to be processed by workers

	go func() {
		defer close(s.results)
		s.walkFileSystem(ctx, rootPath, 0)
	}()

	go func() {
		wg.Wait()
		close(s.errors)
	}()

	// aggregate results
	for r := range s.results {
		finalResult.Folders = append(finalResult.Folders, r)
		finalResult.TotalSize += r.Size
		finalResult.TotalCount++
	}

	finalResult.Duration = time.Since(finalResult.ScanTime)
	return finalResult, nil

}

// worker processes directories and sends results/errors
func (s *Scanner) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

}

func (s *Scanner) walkFileSystem(ctx context.Context, rootPath string, depth int) {
	// File system walking logic
}
