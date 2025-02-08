package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
)

// NewInitCmd returns a new "init" command that initializes configuration directories
// and installs default templates. All dependencies are injected via the deps parameter.
func NewInitCmd(deps Dependencies) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize exo configuration and directories",
		Long: `Initialize the exo configuration and create all necessary directories.
If configuration already exists, it will not be overwritten unless --force is used.

This command creates the required directories and installs the built-in default templates.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use the injected configuration.
			cfg := deps.Config

			// Create required directories.
			if err := ensureDirectories(cfg, deps.Logger, deps.FS); err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}

			// Install default templates.
			if err := installTemplates(cfg, force, deps.Logger, deps.FS); err != nil {
				return fmt.Errorf("failed to install default templates: %w", err)
			}

			deps.Logger.Info("Initialization completed successfully")
			return nil
		},
	}

	// Define GNU-friendly flag for forcing overwrites.
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing configuration and templates")
	return cmd
}

// ensureDirectories creates all required directories as defined in the configuration.
func ensureDirectories(cfg *config.Config, log logger.Logger, fsys fs.FileSystem) error {
	// List all directories that should exist.
	dirs := []string{
		cfg.Dir.DataHome,
		cfg.Dir.IdeaDir,
		cfg.Dir.TemplateDir,
		cfg.Dir.PeriodicDir,
		cfg.Dir.ZettelDir,
	}

	for _, dir := range dirs {
		if err := fsys.EnsureDirectoryExists(dir); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Infof("Created directory %s", dir)
	}
	return nil
}

// installTemplates installs default (built-in) templates into the custom template directory.
// It uses the embedded default templates from the templates package.
func installTemplates(cfg *config.Config, force bool, log logger.Logger, fsys fs.FileSystem) error {
	// Create a default template store from the embedded defaults.
	defaultStore := templates.NewEmbedTemplateStore(templates.DefaultTemplatesFS, templates.DefaultTemplateBaseDir)

	opts := templates.InstallOptions{
		TargetDir: cfg.Dir.TemplateDir,
		Force:     force,
		Reader:    &defaultInputReader{}, // Our interactive input reader implementation.
	}

	// Build a TemplateConfig using the injected logger and file system.
	tmplCfg := templates.TemplateConfig{
		TemplateDir:       cfg.Dir.TemplateDir,
		TemplateExtension: ".md",
		FilePermissions:   0644,
		Logger:            log,
		FS:                fsys,
	}

	if err := templates.InstallDefaultTemplates(tmplCfg, opts, defaultStore); err != nil {
		return err
	}
	log.Info("Default templates installed successfully")
	return nil
}
