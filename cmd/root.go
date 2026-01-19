/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/d4rthvadr/node-cleaner/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	workers int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "depo-cleaner",
	Short: "Clean up large dependency folders (node_modules, vendor, venv, target)",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig() {
	// Initialize configuration
	err := config.Init(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Configuration initialized")
	cfg := config.Load()
	cfg.Workers = workers

}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.depocleaner/config.yaml)")
	rootCmd.PersistentFlags().IntVar(&workers, "workers", 4, "Number of concurrent workers")

}
