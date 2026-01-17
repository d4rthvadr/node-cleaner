package analyzer

import (
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/d4rthvadr/node-cleaner/pkg/models"
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// Analyze inspects the given path and returns a DependencyFolder with its details
func (a *Analyzer) Analyze(path string) (*models.DependencyFolder, error) {

	info, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	folder := &models.DependencyFolder{
		Path:         path,
		AbsolutePath: path,
		ModTime:      info.ModTime(),
		Type:         a.detectType(info.Name()),
	}

	// Calculate size recursively
	// using os.Stat on file returns inode size only
	// this is not accurate for folder size
	// we need to walk the folder and sum up file sizes
	size, err := a.calculateSize(path)
	if err != nil {
		return nil, err
	}
	folder.Size = size
	folder.AccessTime = a.getAccessTime(info)

	return folder, nil
}

func (a *Analyzer) detectType(folderName string) string {
	switch folderName {
	case "node_modules", "node_modules_cache":
		return "Node.js"
	case "vendor":
		return "Go/PHP"
	case ".venv", "venv", "__pycache__":
		return "Python"
	case "target":
		return "Rust"
	default:
		return "Unknown"
	}
}

// getAccessTime uses platform-specific syscall to get the last access time of the file/folder
func (a *Analyzer) getAccessTime(info os.FileInfo) (atime time.Time) {

	statT, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return info.ModTime()
	}

	return time.Unix(int64(statT.Atimespec.Sec), int64(statT.Atimespec.Nsec))
}

// calculateSize computes the total size of all files within the specified path
func (a *Analyzer) calculateSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {

		if err != nil {

			return nil // returning nil to continue walking despite the error
		}

		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err

}
