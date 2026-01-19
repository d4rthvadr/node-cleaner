package cmd

import (
	"fmt"

	"github.com/d4rthvadr/node-cleaner/internal/cache"
	"github.com/d4rthvadr/node-cleaner/internal/config"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cache",
}

var cacheClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cache",
	RunE:  runCacheClear,
}

func init() {
	cacheCmd.AddCommand(cacheClearCmd)
	rootCmd.AddCommand(cacheCmd)
}

func runCacheClear(cmd *cobra.Command, args []string) error {
	cfg := config.Load()

	c, err := cache.NewCache(cfg.CachePath)
	if err != nil {
		fmt.Printf("failed to load cache: %v\n", err)
		return err
	}

	if err := c.Clear(); err != nil {
		fmt.Printf("failed to clear cache: %v\n", err)
		return err
	}

	fmt.Println("Cache cleared successfully.")

	return nil
}
