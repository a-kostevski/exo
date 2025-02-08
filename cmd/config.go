package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/pkg/config"
)

// NewConfigCmd creates a new "config" command with subcommands "get" and "set".
func NewConfigCmd(deps Dependencies) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long: `Manage exo configuration settings.

Without arguments, lists all configuration settings.
Use "get" to retrieve a specific setting.
Use "set" to modify a specific setting.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Simply print the configuration.
			fmt.Println(deps.Config)
		},
	}
	configCmd.AddCommand(NewConfigGetCmd(deps))
	configCmd.AddCommand(NewConfigSetCmd(deps))
	return configCmd
}

func NewConfigGetCmd(deps Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := getConfigValue(deps.Config, key)
			if value == "" {
				deps.Logger.Errorf("Invalid configuration key: %s", key)
				return
			}
			fmt.Printf("%s: %s\n", key, value)
		},
	}
}

func NewConfigSetCmd(deps Dependencies) *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			value := args[1]
			if !setConfigValue(deps.Config, key, value) {
				deps.Logger.Errorf("Invalid configuration key: %s", key)
				return
			}
			if err := deps.Config.Save(); err != nil {
				deps.Logger.Errorf("Failed to save configuration: %v", err)
				return
			}
			deps.Logger.Info("Configuration updated successfully")
			fmt.Printf("Set %s to %s\n", key, value)
		},
	}
}

// getConfigValue returns the configuration value for a given key.
func getConfigValue(cfg *config.Config, key string) string {
	key = strings.ToLower(key)
	switch key {
	case "editor":
		return cfg.General.Editor
	case "data_home", "datahome":
		return cfg.Dir.DataHome
	case "template_dir", "templatedir":
		return cfg.Dir.TemplateDir
	case "periodic_dir", "periodicdir":
		return cfg.Dir.PeriodicDir
	case "zettel_dir", "zetteldir":
		return cfg.Dir.ZettelDir
	case "log.level", "loglevel":
		return cfg.Log.Level
	case "log.format", "logformat":
		return cfg.Log.Format
	case "log.output", "logoutput":
		return cfg.Log.Output
	default:
		return ""
	}
}

// setConfigValue updates the configuration for a given key.
func setConfigValue(cfg *config.Config, key, value string) bool {
	key = strings.ToLower(key)
	switch key {
	case "editor":
		cfg.General.Editor = value
	case "data_home", "datahome":
		cfg.Dir.DataHome = value
	case "template_dir", "templatedir":
		cfg.Dir.TemplateDir = value
	case "periodic_dir", "periodicdir":
		cfg.Dir.PeriodicDir = value
	case "zettel_dir", "zetteldir":
		cfg.Dir.ZettelDir = value
	case "log.level", "loglevel":
		cfg.Log.Level = value
	case "log.format", "logformat":
		cfg.Log.Format = value
	case "log.output", "logoutput":
		cfg.Log.Output = value
	default:
		return false
	}
	return true
}
