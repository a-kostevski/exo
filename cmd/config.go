package cmd

import (
	"fmt"
	"strings"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/logger"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Manage exo configuration settings.
	
Without arguments, lists all configuration settings.
Use 'get' to retrieve a specific setting.
Use 'set' to modify a specific setting.`,
	Run: func(cmd *cobra.Command, args []string) {
		listConfig()
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a configuration value",
	Long:  "Get the value of a specific configuration key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getConfig(args[0])
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  "Set the value of a specific configuration key",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		setConfig(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
}

func listConfig() {
	cfg, err := config.Get()
	if err != nil {
		logger.Errorf("Failed to get configuration: %v", err)
		return
	}

	fmt.Println(cfg)
}

func getConfig(key string) {
	cfg, err := config.Get()
	if err != nil {
		logger.Errorf("Failed to get configuration: %v", err)
		return
	}

	value := getConfigValue(cfg, key)
	if value == "" {
		logger.Errorf("Invalid configuration key: %s", key)
		return
	}

	fmt.Printf("%s: %s\n", key, value)
}

func setConfig(key, value string) {
	cfg, err := config.Get()
	if err != nil {
		logger.Errorf("Failed to get configuration: %v", err)
		return
	}

	if !setConfigValue(cfg, key, value) {
		logger.Errorf("Invalid configuration key: %s", key)
		return
	}

	if err := config.Save(); err != nil {
		logger.Errorf("Failed to save configuration: %v", err)
		return
	}
	logger.Info("Configuration updated successfully")
	fmt.Printf("Set %s to %s\n", key, value)
}

func getConfigValue(cfg *config.Config, key string) string {
	key = strings.ToLower(key)
	switch key {
	case "editor":
		return cfg.Editor
	case "data_home", "datahome":
		return cfg.DataHome
	case "template_dir", "templatedir":
		return cfg.TemplateDir
	case "periodic_dir", "periodicdir":
		return cfg.PeriodicDir
	case "zettel_dir", "zetteldir":
		return cfg.ZettelDir
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

func setConfigValue(cfg *config.Config, key, value string) bool {
	key = strings.ToLower(key)
	switch key {
	case "editor":
		cfg.Editor = value
	case "data_home", "datahome":
		cfg.DataHome = value
	case "template_dir", "templatedir":
		cfg.TemplateDir = value
	case "periodic_dir", "periodicdir":
		cfg.PeriodicDir = value
	case "zettel_dir", "zetteldir":
		cfg.ZettelDir = value
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
