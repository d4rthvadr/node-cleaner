package cleaner

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/d4rthvadr/node-cleaner/pkg/models"
)

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}
type Cleaner struct {
	dryRun bool
	logger Logger
}

func NewCleaner(dryRun bool, logger Logger) *Cleaner {
	return &Cleaner{
		dryRun: dryRun,
		logger: logger,
	}
}

func (c *Cleaner) Clean(ctx context.Context, folders []models.DependencyFolder) (*models.CleanResult, error) {

	result := &models.CleanResult{
		DryRun: c.dryRun,
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, folder := range folders {
		wg.Add(1)

		go func(f models.DependencyFolder) {
			defer wg.Done()

			if err := c.deleteFolder(ctx, f.Path); err != nil {
				mu.Lock()
				result.Failed = append(result.Failed, models.FailedOp{
					Path:   f.Path,
					Reason: err.Error(),
				})
				mu.Unlock()
			} else {
				mu.Lock()
				result.DeletedFolders = append(result.DeletedFolders, f.Path)
				result.SpaceReclaimed += f.Size
				mu.Unlock()
			}
		}(folder)
	}

	wg.Wait()

	return result, nil
}

func (c *Cleaner) deleteFolder(ctx context.Context, path string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// c.logger.Info("Folder does not exist, skipping", "path", path)
		return fmt.Errorf("path no longer exists")
	}

	if c.dryRun {
		// c.logger.Info("Dry run: skipping deletion", "path", path)
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return os.RemoveAll(path)

}
