package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
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

var targetDirectories = []string{
	"node_modules",       // common Node.js dependencies folder
	"node_modules_cache", // alternative Node.js cache folder
	"vendor",             // Go/Php vendor folders
	".venv",              // Python virtual environment
	"__pycache__",        // Python cache folder
	"venv",               // Python virtual environment
	"target",             // Rust target folder
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

func (s *Scanner) isTargetDirectory(name string) bool {

	for _, target := range targetDirectories {
		if name == target {
			return true
		}
	}

	return false

}

func (s *Scanner) walkFileSystem(ctx context.Context, rootPath string, depth int) error {

	pathDepths := make(map[string]int)
	pathDepths[rootPath] = depth
	// WalkDir is a convenient way to traverse directories tho
	// it may need to be replaced with a custom implementation
	// if we need to follow symlinks
	return filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {

		// check for context cancellation
		select {
		case <-ctx.Done():
			return fs.SkipAll
		default:
			// continue processing
		}

		if err != nil {
			s.errors <- fmt.Errorf("accessing %s: %w", path, err)
			return fs.SkipDir // skip this directory on error but keep walking
		}

		if !d.IsDir() {
			return nil // we only care about directories
		}

		parentDepth := pathDepths[filepath.Dir(path)]
		currentDepth := parentDepth + 1
		pathDepths[path] = currentDepth

		// check max depth
		if s.config.MaxDepth > 0 && currentDepth > s.config.MaxDepth {
			return fs.SkipDir
		}

		// TODO: check ignore paths

		if s.isTargetDirectory(d.Name()) {

			// info, _ := d.Info()

			s.enqueueAndAnalysis(path)

			return fs.SkipDir // skip further traversal into this directory

		}

		return nil

	})

}

func (s *Scanner) enqueueAndAnalysis(path string) {

	// send directory for analysis

}
