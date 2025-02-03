package cmd

import (
	"fmt"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize exo configuration and directories",
		Long: `Initialize the exo configuration and create necessary directories.
		If configuration already exists, it will not be overwritten unless --force is used.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize configuration
			logger.Debug("Initializing configuration...")
			if err := config.Initialize(""); err != nil {
				return fmt.Errorf("failed to initialize configuration: %w", err)
			}
			logger.Info("Configuration initialized successfully")

			// Get configuration
			cfg := config.MustGet()

			// Create required directories
			if err := ensureDirectories(cfg); err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}

			logger.Info("Installing default templates...")
			if err := installTemplates(cfg, force); err != nil {
				return fmt.Errorf("failed to install templates: %w", err)
			}
			logger.Info("Default templates installed successfully")

			logger.Info("Initialization completed successfully")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing configuration and templates")
	return cmd
}

// ensureDirectories creates all required directories
func ensureDirectories(cfg *config.Config) error {
	dirs := []string{
		cfg.Dir.DataHome,
		cfg.Dir.IdeaDir,
		cfg.Dir.TemplateDir,
		cfg.Dir.PeriodicDir,
		cfg.Dir.ZettelDir,
	}

	for _, dir := range dirs {
		if err := fs.EnsureDirectoryExists(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logger.Info("Created directory", logger.Field{Key: "path", Value: dir})
	}

	return nil
}

// installTemplates installs default templates
func installTemplates(cfg *config.Config, force bool) error {
	// Install default templates
	if err := templates.InstallDefault(cfg.Dir.TemplateDir, force); err != nil {
		return fmt.Errorf("failed to install templates: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(newInitCmd())
}
