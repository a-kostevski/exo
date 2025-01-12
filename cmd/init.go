package cmd

import (
	"fmt"
	"os"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/logger"
	"github.com/a-kostevski/exo/internal/templates"
	"github.com/spf13/cobra"
)

// newInitCmd creates a new init command
func newInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize exo configuration and directories",
		Long: `Initialize the exo configuration and create necessary directories.
		If configuration already exists, it will not be overwritten unless --force is used.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Info("Starting initialization...")
			logger.Debug("Initializing with force flag:", force)

			// Initialize configuration
			logger.Debug("Initializing configuration...")
			if err := config.Initialize(""); err != nil {
				return fmt.Errorf("failed to initialize configuration: %w", err)
			}
			logger.Info("Configuration initialized successfully")

			// Get configuration
			cfg := config.MustGet()
			logger.Debug("Configuration loaded")

			// Create required directories
			dirs := []string{
				cfg.DataHome,
				cfg.TemplateDir,
				cfg.PeriodicDir,
				cfg.ZettelDir,
			}

			logger.Info("Creating required directories...")
			for _, dir := range dirs {
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					logger.Debugf("Creating directory: %s", dir)
					if err := os.MkdirAll(dir, 0755); err != nil {
						return fmt.Errorf("failed to create directory %s: %w", dir, err)
					}
					logger.Infof("Created directory: %s", dir)
				} else {
					logger.Infof("Directory already exists: %s", dir)
				}
			}

			logger.Info("Installing default templates...")
			if err := templates.InstallDefault(cfg.TemplateDir, force); err != nil {
				return fmt.Errorf("failed to copy default templates: %w", err)
			}
			logger.Info("Default templates installed successfully")

			logger.Info("Initialization completed successfully")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing configuration and templates")
	return cmd
}

func init() {
	rootCmd.AddCommand(newInitCmd())
}
