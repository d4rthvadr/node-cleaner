package cmd

import (
	"fmt"
	"strconv"

	"github.com/d4rthvadr/node-cleaner/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration settings",
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config.Display()
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			cmd.Help()
			return
		}
		key, value := args[0], args[1]

		if err := config.Set(key, parseValue(value)); err != nil {
			cmd.PrintErrln("Error setting config:", err)
			return
		}

		fmt.Printf("Set %s = %s\n", key, value)

	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to default values",
	Run: func(cmd *cobra.Command, args []string) {

		if err := config.RestoreDefaults(); err != nil {
			cmd.PrintErrln("Error restoring defaults:", err)
			return
		}

		fmt.Println("Configuration reset to default values.")

	},
}

// attempt to parse value into appropriate type or return as string
func parseValue(value string) interface{} {

	if value == "true" || value == "false" {
		// Parse boolean
		return value == "true"
	}

	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	// Default to value
	return value
}

func init() {

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	rootCmd.AddCommand(configCmd)
}
