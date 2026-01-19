package models

import "time"

type DependencyFolder struct {
	Path         string    `json:"path"`
	AbsolutePath string    `json:"absolute_path"`
	Size         int64     `json:"size"`
	ModTime      time.Time `json:"mod_time"`
	AccessTime   time.Time `json:"access_time"`
	Type         string    `json:"type"`
	Selected     bool      `json:"selected"`
}

type FailedOp struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type Config struct {
	ScanPaths      []string `mapstructure:"scan_paths" json:"scan_paths"`
	IgnorePaths    []string `mapstructure:"ignore_paths" json:"ignore_paths"`
	CachePath      string   `mapstructure:"cache_path" json:"cache_path"`
	LogPath        string   `mapstructure:"log_path" json:"log_path"`
	FollowSymlinks bool     `mapstructure:"follow_symlinks" json:"follow_symlinks"`
	MaxDepth       int      `mapstructure:"max_depth" json:"max_depth"`
	Workers        int      `mapstructure:"workers" json:"workers"`
}

// CacheEntry represents a cached folder information
type CacheEntry struct {
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"mod_time"`
	LastScan time.Time `json:"last_scan"`
	Hash     string    `json:"hash,omitempty"` // optional hash of folder contents
}

// CacheIndex represents the overall cache root structure
type CacheIndex struct {
	Version   string                `json:"version"`
	Entries   map[string]CacheEntry `json:"entries"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// ScanResult represents the result of a scan operation
type ScanResult struct {
	Folders     []DependencyFolder `json:"folders"`
	TotalSize   int64              `json:"total_size"`
	TotalCount  int                `json:"total_count"`
	ScanPath    string             `json:"scan_path"`
	ScanTime    time.Time          `json:"scan_time"`
	Duration    time.Duration      `json:"duration"`
	CacheHits   int                `json:"cache_hits"`
	CacheMisses int                `json:"cache_misses"`
}
