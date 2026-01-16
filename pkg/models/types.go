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
	ScanPaths      []string `json:"scan_paths"`
	IgnorePaths    []string `json:"ignore_paths"`
	CachePath      string   `json:"-"`
	LogPath        string   `json:"log_path"`
	FollowSymlinks bool     `json:"follow_symlinks"`
	MaxDepth       int      `json:"max_depth"`
	Workers        int      `json:"workers"`
}
