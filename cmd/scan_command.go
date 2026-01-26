package cmd

import (
	"fmt"
	"os"

	"github.com/d4rthvadr/node-cleaner/internal/cache"
	"github.com/d4rthvadr/node-cleaner/internal/config"
	"github.com/d4rthvadr/node-cleaner/internal/scanner"
	"github.com/d4rthvadr/node-cleaner/internal/ui"
	"github.com/spf13/cobra"
)

var (
	scanPath string
	noCache  bool
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan for dependency folders",
	Args:  cobra.MaximumNArgs(1),
	Run:   runScan,
}

func init() {

	// Scan command flags
	scanCmd.Flags().StringVarP(&scanPath, "path", "p", "", "Path to scan for dependency folders(default: $HOME)")
	scanCmd.Flags().BoolVar(&noCache, "no-cache", false, "Disable cache")

	rootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, args []string) {

	ctx := cmd.Context()

	path := scanPath

	if len(args) > 0 {
		path = args[0]
	}

	// if path is still empty, load from config
	if path == "" {
		path = os.Getenv("HOME")
	}

	cfg := config.Load()
	cfg.ScanPath = path
	fmt.Printf("config loaded %v", cfg)
	fmt.Printf("properties loaded workers: %v, scanPaths: %v, cachePath: %v, logPath: %v\n", cfg.Workers, cfg.ScanPath, cfg.CachePath, cfg.LogPath)

	// Initialize cache
	var c *cache.Cache
	var err error

	if !noCache {
		c, err = cache.NewCache(cfg.CachePath)
		if err != nil {
			fmt.Printf("failed to initialize cache: %v", err)
			os.Exit(1)
		}
		err = c.Save() // Save cache in case unsaved changes or exit occurs
		if err != nil {
			fmt.Printf("failed to save cache: %v", err)
		}
		fmt.Println("Cache initialized at", cfg.CachePath)
	}

	// Create scanner
	s := scanner.NewScanner(cfg, c)

	// Start scan
	fmt.Printf("Starting scan on path: %s\n", path)
	result, err := s.Scan(ctx, path)
	if err != nil {
		fmt.Printf("Scan failed: %v\n", err)
		os.Exit(1)
	}

	// Display results
	ui.DisplayScanResults(result)

}
