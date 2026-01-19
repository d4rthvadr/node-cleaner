package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/d4rthvadr/node-cleaner/pkg/models"
	"github.com/spf13/viper"
)

var (
	globalConfig *models.Config
	configPath   string
)

func Init(cfgFile string) error {

	if cfgFile != "" {
		configPath = cfgFile
	} else {

		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting home directory: %w", err)
		}

		configDir := filepath.Join(home, ".nodecleaner")

		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}

		viper.AddConfigPath(configDir)
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		configPath = filepath.Join(configDir, "config.yaml")
	}

	setDefaults()

	// enable environment variables support
	viper.SetEnvPrefix("NODECLEANER")
	viper.AutomaticEnv()

	// read in config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			if err := viper.WriteConfigAs(configPath); err != nil {
				return fmt.Errorf("writing default config file: %w", err)
			}
		} else {
			// Config file was found but another error was produced
			return fmt.Errorf("reading config file: %w", err)
		}
	}

	// we are here means config file is found and read successfully
	globalConfig = &models.Config{}
	if err := viper.Unmarshal(globalConfig); err != nil {
		return fmt.Errorf("unmarshaling config file: %w", err)
	}

	return nil
}

func setDefaults() {

	home, _ := os.UserHomeDir()

	// TODO: extract to constants package
	configDir := filepath.Join(home, ".nodecleaner")

	viper.SetDefault("scan_path", home)
	viper.SetDefault("cache_path", filepath.Join(configDir, "cache.json"))
	viper.SetDefault("log_path", filepath.Join(configDir, "nodecleaner.log"))
	viper.SetDefault("follow_symlinks", false)
	viper.SetDefault("max_depth", 10)
	viper.SetDefault("workers", 4)
	// TODO: allow user to customize or add additional ignore paths
	viper.SetDefault("ignore_paths", []string{
		"/System",
		"/Library",
		"/Applications",
		"/private/var",
		"/dev",
		"/proc",
		"/sys",
		"/.Trash",
		"/Network",
	})

}

func Load() *models.Config {

	if globalConfig == nil {
		globalConfig = &models.Config{
			Workers:  viper.GetInt("workers"),
			MaxDepth: viper.GetInt("max_depth"),
		}

		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, ".nodecleaner")
		globalConfig.CachePath = filepath.Join(configDir, "cache.json")
		globalConfig.LogPath = filepath.Join(configDir, "nodecleaner.log")
	}
	return globalConfig

}

// Display prints the current configuration settings to the console
func Display() {
	fmt.Println("Current Configuration:")
	fmt.Println("======================")

	allSettings := viper.AllSettings()
	for key, value := range allSettings {
		fmt.Printf("%-20s: %v\n", key, value)
	}
}

func SaveDefaultsToFile() error {

	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".nodecleaner")
	configPath := filepath.Join(configDir, "config.yaml")

	// should not error if dir exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	return viper.SafeWriteConfigAs(configPath)

}

// Restores default configuration
func RestoreDefaults() error {

	for key := range viper.AllSettings() {
		viper.Set(key, nil)
	}

	setDefaults()

	return SaveDefaultsToFile()

}

// Helper functions to get specific config values

// persistConfigChanges writes the current in-memory config to the config file
func persistConfigChanges() error {
	if globalConfig != nil {
		if err := viper.Unmarshal(globalConfig); err != nil {
			return fmt.Errorf("unmarshaling config after set: %w", err)
		}
	}
	return nil
}

// Update a configuration value and save changes to file
func Set(key string, value interface{}) error {
	// this saves the value in viper's in-memory config
	viper.Set(key, value)

	// but we need to persist the change to the config file shared
	// across application runs

	err := persistConfigChanges()
	if err != nil {
		return fmt.Errorf("persisting config changes: %w", err)
	}

	return viper.WriteConfigAs(configPath)
}

func GetStringValue(key string) string {
	return viper.GetString(key)
}

func GetIntValue(key string) int {
	return viper.GetInt(key)
}

func GetBoolValue(key string) bool {
	return viper.GetBool(key)
}

func GetStringSliceValue(key string) []string {
	return viper.GetStringSlice(key)
}
