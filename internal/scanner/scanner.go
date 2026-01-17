package scanner

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"
	"time"

	"github.com/d4rthvadr/node-cleaner/internal/analyzer"
	"github.com/d4rthvadr/node-cleaner/pkg/models"
	"github.com/d4rthvadr/node-cleaner/pkg/utils"
)

type Scanner struct {
	config    *models.Config
	cache     CacheProvider
	results   chan models.DependencyFolder
	errors    chan error
	analyzer  *analyzer.Analyzer
	workQueue chan string
}

type CacheProvider interface {
	Get(path string) (*models.CacheEntry, bool)
	Set(path string, entry *models.CacheEntry) error
	IsValid(path string, modTime time.Time) bool
}

// NewScanner creates a new Scanner instance
func NewScanner(cfg *models.Config, cache CacheProvider) *Scanner {
	return &Scanner{
		analyzer: analyzer.NewAnalyzer(),
		config:   cfg,
		cache:    cache,
		results:  make(chan models.DependencyFolder),
		errors:   make(chan error),
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

		if utils.IsTargetDirectory(d.Name()) {

			info, _ := d.Info()

			if s.cache != nil && s.cache.IsValid(path, info.ModTime()) {

				// use cached data
				cached, _ := s.cache.Get(path)
				s.results <- models.DependencyFolder{
					Path:         path,
					AbsolutePath: path,
					Size:         cached.Size,
					ModTime:      cached.ModTime,
					Type:         utils.DetectType(d.Name()),
				}

			} else {
				s.enqueueAndAnalysis(path)
			}

			return fs.SkipDir // skip further traversal into this directory

		}

		// Continue walking
		return nil

	})

}

func (s *Scanner) enqueueAndAnalysis(path string) {

	// send path to worker pool
	select {
	case s.workQueue <- path:
	case <-time.After(100 * time.Millisecond):

		folder, err := s.analyzer.Analyze(path)
		if err != nil {
			s.errors <- fmt.Errorf("analyzing %s: %w", path, err)
			return
		}
		s.results <- *folder
	}

}

// worker processes directories and sends results/errors

func (s *Scanner) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-s.workQueue:
			if !ok {
				return // Channel closed
			}

			folder, err := s.analyzer.Analyze(path)
			if err != nil {
				s.errors <- fmt.Errorf("analyzing %s: %w", path, err)
				continue
			}

			// Cache the result if caching is enabled
			if s.cache != nil {
				s.cache.Set(path, &models.CacheEntry{
					Path:     path,
					Size:     folder.Size,
					ModTime:  folder.ModTime,
					LastScan: time.Now(),
				})
			}

			// context could cancelled while sending result
			select {
			case s.results <- *folder:
			case <-ctx.Done():
				return

			}

		}
	}
}
