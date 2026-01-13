# NodeCleaner Technical Specification

**Version:** 1.0  
**Last Updated:** January 13, 2026  
**Status:** Draft  
**Language:** Go 1.20+  
**Target Platforms:** macOS, Linux

---

## Table of Contents

1. [Technical Overview](#1-technical-overview)
2. [System Architecture](#2-system-architecture)
3. [Technology Stack](#3-technology-stack)
4. [Data Models](#4-data-models)
5. [Core Components](#5-core-components)
6. [Algorithms & Logic](#6-algorithms--logic)
7. [File System Operations](#7-file-system-operations)
8. [Caching Strategy](#8-caching-strategy)
9. [CLI Interface Design](#9-cli-interface-design)
10. [Error Handling](#10-error-handling)
11. [Testing Strategy](#11-testing-strategy)
12. [Performance Optimization](#12-performance-optimization)
13. [Security Considerations](#13-security-considerations)
14. [Build & Distribution](#14-build--distribution)
15. [Development Workflow](#15-development-workflow)

---

## 1. Technical Overview

### 1.1 Architecture Philosophy

NodeCleaner follows a **modular, pipeline-based architecture** where data flows through discrete stages:

```
Scan → Analyze → Present → Select → Execute → Report
```

**Key Principles:**

- **Single Responsibility**: Each package handles one concern
- **Testability**: Components are mockable and unit-testable
- **Performance**: Concurrent operations where safe
- **Safety**: Multiple validation checkpoints before destructive operations

### 1.2 Project Structure

```
nodecleaner/
├── cmd/
│   └── nodecleaner/
│       └── main.go              # Application entry point
├── internal/
│   ├── scanner/
│   │   ├── scanner.go           # Filesystem scanning logic
│   │   ├── walker.go            # Directory traversal
│   │   └── filter.go            # Path filtering/ignoring
│   ├── analyzer/
│   │   ├── analyzer.go          # Metadata collection & analysis
│   │   └── stats.go             # Size calculation utilities
│   ├── cache/
│   │   ├── cache.go             # Cache management
│   │   ├── store.go             # Persistence layer
│   │   └── invalidation.go      # Cache invalidation logic
│   ├── ui/
│   │   ├── presenter.go         # Results display
│   │   ├── selector.go          # Interactive selection
│   │   └── formatter.go         # Output formatting
│   ├── cleaner/
│   │   ├── cleaner.go           # Deletion operations
│   │   └── validator.go         # Pre-deletion validation
│   ├── config/
│   │   ├── config.go            # Configuration management
│   │   └── defaults.go          # Default settings
│   └── logger/
│       └── logger.go            # Logging infrastructure
├── pkg/
│   └── models/
│       └── types.go             # Shared data types
├── test/
│   ├── integration/
│   └── fixtures/
├── scripts/
│   ├── build.sh
│   └── release.sh
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 2. System Architecture

### 2.1 Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  (Cobra commands, flag parsing, user interaction)           │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                     Application Core                         │
│                                                               │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐    │
│  │   Scanner   │──▶│   Analyzer   │──▶│   Presenter    │    │
│  └─────────────┘  └──────────────┘  └─────────────────┘    │
│         │                                      │              │
│         │                                      │              │
│  ┌──────▼──────┐                      ┌───────▼────────┐    │
│  │    Cache    │                      │    Selector    │    │
│  └─────────────┘                      └───────┬────────┘    │
│                                               │              │
│                                       ┌───────▼────────┐    │
│                                       │    Cleaner     │    │
│                                       └────────────────┘    │
└───────────────────────────────────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────────┐
│                   Infrastructure Layer                       │
│  (Filesystem, JSON storage, logging)                        │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Data Flow

**Scan Operation:**

```
User Input → Config → Scanner → Walker → Filter → Analyzer
    → Cache Check → Results → Presenter → User
```

**Clean Operation:**

```
User Selection → Validator → Cleaner → Filesystem → Logger
    → Reporter → User
```

---

## 3. Technology Stack

### 3.1 Core Dependencies

```go
// go.mod
module github.com/yourusername/nodecleaner

go 1.20

require (
    github.com/spf13/cobra v1.8.0          // CLI framework
    github.com/spf13/viper v1.18.0         // Configuration management
    github.com/charmbracelet/bubbletea v0.25.0  // TUI framework
    github.com/charmbracelet/bubbles v0.18.0    // TUI components
    github.com/charmbracelet/lipgloss v0.9.1    // Styling
    github.com/dustin/go-humanize v1.0.1   // Human-readable sizes
    github.com/schollz/progressbar/v3 v3.14.1   // Progress bars
    go.uber.org/zap v1.26.0                // Structured logging
    github.com/stretchr/testify v1.8.4     // Testing utilities
)
```

### 3.2 Standard Library Usage

- **os**: File operations, environment variables
- **path/filepath**: Cross-platform path handling
- **encoding/json**: Cache persistence
- **sync**: Concurrent operations (WaitGroup, Mutex)
- **time**: Timestamp handling
- **io/fs**: Filesystem abstraction (Go 1.16+)

### 3.3 Rationale for Key Choices

**Cobra**: Industry-standard CLI framework, excellent documentation, Git-style subcommands

**Bubbletea**: Modern TUI framework for interactive selection, clean architecture

**Zap**: High-performance structured logging with minimal allocations

**Standard Library First**: Minimize dependencies for security and binary size

---

## 4. Data Models

### 4.1 Core Types

```go
// pkg/models/types.go

package models

import "time"

// DependencyFolder represents a detected dependency directory
type DependencyFolder struct {
    Path         string    `json:"path"`
    AbsolutePath string    `json:"absolute_path"`
    Size         int64     `json:"size"`           // bytes
    ModTime      time.Time `json:"mod_time"`
    AccessTime   time.Time `json:"access_time"`
    Type         string    `json:"type"`           // "node_modules", etc.
    Selected     bool      `json:"-"`              // UI state, not persisted
}

// ScanResult encapsulates all findings from a scan
type ScanResult struct {
    Folders      []DependencyFolder `json:"folders"`
    TotalSize    int64              `json:"total_size"`
    TotalCount   int                `json:"total_count"`
    ScanPath     string             `json:"scan_path"`
    ScanTime     time.Time          `json:"scan_time"`
    Duration     time.Duration      `json:"duration"`
    CacheHits    int                `json:"cache_hits"`
    CacheMisses  int                `json:"cache_misses"`
}

// CacheEntry stores cached folder information
type CacheEntry struct {
    Path       string    `json:"path"`
    Size       int64     `json:"size"`
    ModTime    time.Time `json:"mod_time"`
    LastScan   time.Time `json:"last_scan"`
    Hash       string    `json:"hash,omitempty"`  // Optional content hash
}

// CacheIndex is the root cache structure
type CacheIndex struct {
    Version   string                `json:"version"`
    Entries   map[string]CacheEntry `json:"entries"`
    UpdatedAt time.Time             `json:"updated_at"`
}

// CleanResult records deletion outcomes
type CleanResult struct {
    Deleted       []string      `json:"deleted"`
    Failed        []FailedOp    `json:"failed"`
    SpaceReclaimed int64        `json:"space_reclaimed"`
    Duration      time.Duration `json:"duration"`
    DryRun        bool          `json:"dry_run"`
}

// FailedOp records a failed deletion
type FailedOp struct {
    Path   string `json:"path"`
    Reason string `json:"reason"`
}

// Config holds application configuration
type Config struct {
    ScanPath      string   `json:"scan_path"`
    IgnorePaths   []string `json:"ignore_paths"`
    CachePath     string   `json:"cache_path"`
    LogPath       string   `json:"log_path"`
    FollowSymlinks bool    `json:"follow_symlinks"`
    MaxDepth      int      `json:"max_depth"`
    Workers       int      `json:"workers"`
}
```

---

## 5. Core Components

### 5.1 Scanner Package

**Purpose**: Traverse filesystem and detect dependency directories

```go
// internal/scanner/scanner.go

package scanner

import (
    "context"
    "path/filepath"
    "sync"

    "github.com/yourusername/nodecleaner/pkg/models"
)

type Scanner struct {
    config  *models.Config
    cache   CacheProvider
    filter  *Filter
    results chan models.DependencyFolder
    errors  chan error
}

// CacheProvider abstracts cache operations
type CacheProvider interface {
    Get(path string) (*models.CacheEntry, bool)
    Set(path string, entry models.CacheEntry) error
}

// NewScanner creates a configured scanner
func NewScanner(cfg *models.Config, cache CacheProvider) *Scanner {
    return &Scanner{
        config:  cfg,
        cache:   cache,
        filter:  NewFilter(cfg.IgnorePaths),
        results: make(chan models.DependencyFolder, 100),
        errors:  make(chan error, 10),
    }
}

// Scan performs filesystem traversal
func (s *Scanner) Scan(ctx context.Context, rootPath string) (*models.ScanResult, error) {
    result := &models.ScanResult{
        ScanPath: rootPath,
        ScanTime: time.Now(),
    }

    var wg sync.WaitGroup

    // Launch worker pool
    for i := 0; i < s.config.Workers; i++ {
        wg.Add(1)
        go s.worker(ctx, &wg)
    }

    // Walk filesystem
    go func() {
        defer close(s.results)
        s.walk(ctx, rootPath, 0)
    }()

    // Collect results
    go func() {
        wg.Wait()
        close(s.errors)
    }()

    // Aggregate
    for folder := range s.results {
        result.Folders = append(result.Folders, folder)
        result.TotalSize += folder.Size
        result.TotalCount++
    }

    result.Duration = time.Since(result.ScanTime)
    return result, nil
}

// walk recursively traverses directories
func (s *Scanner) walk(ctx context.Context, path string, depth int) {
    // Implementation in next section
}

// worker processes discovered directories
func (s *Scanner) worker(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    // Implementation in next section
}
```

### 5.2 Analyzer Package

**Purpose**: Collect metadata and calculate statistics

```go
// internal/analyzer/analyzer.go

package analyzer

import (
    "os"
    "path/filepath"
    "syscall"

    "github.com/yourusername/nodecleaner/pkg/models"
)

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
    return &Analyzer{}
}

// Analyze collects metadata for a directory
func (a *Analyzer) Analyze(path string) (*models.DependencyFolder, error) {
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    folder := &models.DependencyFolder{
        Path:         path,
        AbsolutePath: path,
        ModTime:      info.ModTime(),
        Type:         a.detectType(path),
    }

    // Calculate size
    size, err := a.calculateSize(path)
    if err != nil {
        return nil, err
    }
    folder.Size = size

    // Get access time (platform-specific)
    folder.AccessTime = a.getAccessTime(info)

    return folder, nil
}

// calculateSize recursively sums directory size
func (a *Analyzer) calculateSize(path string) (int64, error) {
    var size int64

    err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
        if err != nil {
            return nil // Skip inaccessible files
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return nil
    })

    return size, err
}

// detectType determines folder type
func (a *Analyzer) detectType(path string) string {
    base := filepath.Base(path)
    switch base {
    case "node_modules":
        return "node_modules"
    default:
        return "unknown"
    }
}

// getAccessTime extracts access time (platform-specific)
func (a *Analyzer) getAccessTime(info os.FileInfo) time.Time {
    stat := info.Sys().(*syscall.Stat_t)
    // macOS/Linux specific
    return time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
}
```

### 5.3 Cache Package

**Purpose**: Persist scan results to avoid redundant work

```go
// internal/cache/cache.go

package cache

import (
    "encoding/json"
    "os"
    "path/filepath"
    "sync"

    "github.com/yourusername/nodecleaner/pkg/models"
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

    // Load existing cache
    if err := c.load(); err != nil && !os.IsNotExist(err) {
        return nil, err
    }

    return c, nil
}

// Get retrieves a cache entry
func (c *Cache) Get(path string) (*models.CacheEntry, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    entry, ok := c.index.Entries[path]
    return &entry, ok
}

// Set stores a cache entry
func (c *Cache) Set(path string, entry models.CacheEntry) error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.index.Entries[path] = entry
    c.modified = true
    return nil
}

// IsValid checks if cached entry is still valid
func (c *Cache) IsValid(path string, currentModTime time.Time) bool {
    entry, ok := c.Get(path)
    if !ok {
        return false
    }

    // Invalid if modification time changed
    return entry.ModTime.Equal(currentModTime)
}

// Save persists cache to disk
func (c *Cache) Save() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if !c.modified {
        return nil
    }

    c.index.UpdatedAt = time.Now()

    // Ensure directory exists
    dir := filepath.Dir(c.path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    // Write to temp file, then rename (atomic)
    tempPath := c.path + ".tmp"
    f, err := os.Create(tempPath)
    if err != nil {
        return err
    }

    encoder := json.NewEncoder(f)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(c.index)
    f.Close()

    if err != nil {
        os.Remove(tempPath)
        return err
    }

    return os.Rename(tempPath, c.path)
}

// load reads cache from disk
func (c *Cache) load() error {
    f, err := os.Open(c.path)
    if err != nil {
        return err
    }
    defer f.Close()

    return json.NewDecoder(f).Decode(&c.index)
}

// Clear removes all cache entries
func (c *Cache) Clear() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.index.Entries = make(map[string]models.CacheEntry)
    c.modified = true
    return c.Save()
}

// Prune removes cache entries for non-existent paths
func (c *Cache) Prune() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    for path := range c.index.Entries {
        if _, err := os.Stat(path); os.IsNotExist(err) {
            delete(c.index.Entries, path)
            c.modified = true
        }
    }

    return nil
}
```

### 5.4 UI Package

**Purpose**: Interactive terminal interface for selection

```go
// internal/ui/selector.go

package ui

import (
    "fmt"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/lipgloss"
    "github.com/dustin/go-humanize"

    "github.com/yourusername/nodecleaner/pkg/models"
)

type SelectionModel struct {
    table         table.Model
    folders       []models.DependencyFolder
    selected      map[int]bool
    totalSelected int64
}

func NewSelectionModel(folders []models.DependencyFolder) SelectionModel {
    // Create table columns
    columns := []table.Column{
        {Title: "✓", Width: 3},
        {Title: "Size", Width: 10},
        {Title: "Last Accessed", Width: 12},
        {Title: "Path", Width: 60},
    }

    // Convert folders to rows
    rows := make([]table.Row, len(folders))
    for i, f := range folders {
        rows[i] = table.Row{
            "",
            humanize.Bytes(uint64(f.Size)),
            humanize.Time(f.AccessTime),
            f.Path,
        }
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
        table.WithHeight(20),
    )

    return SelectionModel{
        table:    t,
        folders:  folders,
        selected: make(map[int]bool),
    }
}

func (m SelectionModel) Init() tea.Cmd {
    return nil
}

func (m SelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case " ":
            // Toggle selection
            idx := m.table.Cursor()
            m.selected[idx] = !m.selected[idx]

            // Update total
            if m.selected[idx] {
                m.totalSelected += m.folders[idx].Size
            } else {
                m.totalSelected -= m.folders[idx].Size
            }

            // Update row
            m.updateRow(idx)

        case "enter":
            // Confirm selection
            return m, tea.Quit
        }
    }

    var cmd tea.Cmd
    m.table, cmd = m.table.Update(msg)
    return m, cmd
}

func (m SelectionModel) View() string {
    header := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("12")).
        Render(fmt.Sprintf("Select folders to delete (Space to toggle, Enter to confirm)\nSelected: %s\n\n",
            humanize.Bytes(uint64(m.totalSelected))))

    return header + m.table.View() + "\n"
}

func (m *SelectionModel) updateRow(idx int) {
    checkmark := ""
    if m.selected[idx] {
        checkmark = "✓"
    }

    f := m.folders[idx]
    m.table.SetRows([]table.Row{
        {
            checkmark,
            humanize.Bytes(uint64(f.Size)),
            humanize.Time(f.AccessTime),
            f.Path,
        },
    })
}

// GetSelected returns selected folders
func (m SelectionModel) GetSelected() []models.DependencyFolder {
    var selected []models.DependencyFolder
    for i, isSelected := range m.selected {
        if isSelected {
            selected = append(selected, m.folders[i])
        }
    }
    return selected
}
```

### 5.5 Cleaner Package

**Purpose**: Safe deletion operations

```go
// internal/cleaner/cleaner.go

package cleaner

import (
    "context"
    "os"
    "sync"

    "github.com/yourusername/nodecleaner/pkg/models"
)

type Cleaner struct {
    dryRun bool
    logger Logger
}

type Logger interface {
    Info(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
}

func NewCleaner(dryRun bool, logger Logger) *Cleaner {
    return &Cleaner{
        dryRun: dryRun,
        logger: logger,
    }
}

// Clean removes selected folders
func (c *Cleaner) Clean(ctx context.Context, folders []models.DependencyFolder) (*models.CleanResult, error) {
    result := &models.CleanResult{
        DryRun: c.dryRun,
    }

    var mu sync.Mutex
    var wg sync.WaitGroup

    // Process deletions concurrently
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
                c.logger.Error("Failed to delete", "path", f.Path, "error", err)
            } else {
                mu.Lock()
                result.Deleted = append(result.Deleted, f.Path)
                result.SpaceReclaimed += f.Size
                mu.Unlock()
                c.logger.Info("Deleted", "path", f.Path, "size", f.Size)
            }
        }(folder)
    }

    wg.Wait()
    return result, nil
}

// deleteFolder removes a directory
func (c *Cleaner) deleteFolder(ctx context.Context, path string) error {
    // Validate path still exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return fmt.Errorf("path no longer exists")
    }

    if c.dryRun {
        c.logger.Info("DRY RUN: Would delete", "path", path)
        return nil
    }

    // Check context cancellation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    return os.RemoveAll(path)
}
```

---

## 6. Algorithms & Logic

### 6.1 Filesystem Traversal Algorithm

**Approach**: Breadth-First Search (BFS) with early termination

```go
// internal/scanner/walker.go

func (s *Scanner) walk(ctx context.Context, path string, depth int) {
    // Respect max depth
    if s.config.MaxDepth > 0 && depth > s.config.MaxDepth {
        return
    }

    // Check context cancellation
    select {
    case <-ctx.Done():
        return
    default:
    }

    // Check if path should be ignored
    if s.filter.ShouldIgnore(path) {
        return
    }

    entries, err := os.ReadDir(path)
    if err != nil {
        s.errors <- fmt.Errorf("reading %s: %w", path, err)
        return
    }

    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }

        fullPath := filepath.Join(path, entry.Name())

        // Check if this is a target directory
        if s.isTargetDir(entry.Name()) {
            // Check cache first
            info, _ := entry.Info()
            if s.cache != nil && s.cache.IsValid(fullPath, info.ModTime()) {
                // Use cached data
                cached, _ := s.cache.Get(fullPath)
                s.results <- models.DependencyFolder{
                    Path:       fullPath,
                    Size:       cached.Size,
                    ModTime:    cached.ModTime,
                    AccessTime: cached.ModTime, // Approximation
                    Type:       s.detectType(entry.Name()),
                }
                continue
            }

            // Send for analysis
            s.results <- s.queueAnalysis(fullPath)

            // Don't recurse into node_modules
            continue
        }

        // Recurse into subdirectories
        s.walk(ctx, fullPath, depth+1)
    }
}

func (s *Scanner) isTargetDir(name string) bool {
    return name == "node_modules"
}
```

### 6.2 Cache Invalidation Logic

**Strategy**: Timestamp-based with lazy pruning

```go
// internal/cache/invalidation.go

// ShouldRescan determines if a path needs rescanning
func (c *Cache) ShouldRescan(path string) (bool, error) {
    entry, exists := c.Get(path)
    if !exists {
        return true, nil // Not in cache
    }

    // Check if path still exists
    info, err := os.Stat(path)
    if os.IsNotExist(err) {
        // Path deleted, remove from cache
        c.mu.Lock()
        delete(c.index.Entries, path)
        c.modified = true
        c.mu.Unlock()
        return false, nil
    }
    if err != nil {
        return true, err
    }

    // Compare modification times
    if !info.ModTime().Equal(entry.ModTime) {
        return true, nil // Modified since last scan
    }

    // Check cache age (optional: invalidate after 7 days)
    cacheAge := time.Since(entry.LastScan)
    if cacheAge > 7*24*time.Hour {
        return true, nil
    }

    return false, nil
}
```

### 6.3 Size Calculation Optimization

**Strategy**: Parallel calculation with goroutines

```go
// internal/analyzer/stats.go

func (a *Analyzer) calculateSizeParallel(path string) (int64, error) {
    var (
        totalSize atomic.Int64
        wg        sync.WaitGroup
        errChan   = make(chan error, 1)
        semaphore = make(chan struct{}, 10) // Limit concurrent operations
    )

    err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil // Skip errors
        }

        if d.IsDir() {
            return nil
        }

        wg.Add(1)
        semaphore <- struct{}{} // Acquire

        go func(path string) {
            defer wg.Done()
            defer func() { <-semaphore }() // Release

            info, err := os.Stat(path)
            if err == nil {
                totalSize.Add(info.Size())
            }
        }(p)

        return nil
    })

    wg.Wait()
    close(errChan)

    if err != nil {
        return 0, err
    }

    return totalSize.Load(), nil
}
```

---

## 7. File System Operations

### 7.1 Path Filtering

```go
// internal/scanner/filter.go

type Filter struct {
    ignorePaths []string
    ignoreRegex []*regexp.Regexp
}

func NewFilter(ignorePaths []string) *Filter {
    return &Filter{
        ignorePaths: append(defaultIgnorePaths(), ignorePaths...),
    }
}

func defaultIgnorePaths() []string {
    return []string{
        "/System",
        "/Library",
        "/Applications",
        "/private/var",
        "/dev",
        "/proc",
        "/sys",
        "/.Trash",
        "/Network",
    }
}

func (f *Filter) ShouldIgnore(path string) bool {
    // Check exact matches
    for _, ignore := range f.ignorePaths {
        if strings.HasPrefix(path, ignore) {
            return true
        }
    }

    // Check hidden directories
    base := filepath.Base(path)
    if strings.HasPrefix(base, ".") && base != ".config" {
        return true
    }

    return false
}
```

### 7.2 Permission Handling

```go
// internal/scanner/permissions.go

func (s *Scanner) canAccess(path string) bool {
    // Try to stat the path
    _, err := os.Stat(path)
```
