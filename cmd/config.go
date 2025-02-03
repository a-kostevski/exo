package cmd

import (
	"fmt"
	"strings"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/logger"
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

	if err := cfg.Save(); err != nil {
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
		return string(cfg.Log.Level)
	case "log.format", "logformat":
		return string(cfg.Log.Format)
	case "log.output", "logoutput":
		return string(cfg.Log.Output)
	default:
		return ""
	}
}

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
		cfg.Log.Level = logger.Level(value)
	case "log.format", "logformat":
		cfg.Log.Format = logger.Format(value)
	case "log.output", "logoutput":
		cfg.Log.Output = logger.OutputType(value)
	default:
		return false
	}
	return true
}
