package main

import (
	"os"

	"github.com/a-kostevski/exo/cmd"
	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
)

func main() {
	// Initialize configuration.
	cfg, err := config.NewConfig("")
	if err != nil {
		os.Exit(1)
	}

	// Build remaining dependencies.
	log := logger.NewLogger()
	fsys := fs.NewOSFileSystem()
	tm, err := templates.NewTemplateManager(templates.TemplateConfig{
		TemplateDir:       cfg.Dir.TemplateDir,
		TemplateExtension: ".md",
		FilePermissions:   0644,
		Logger:            log,
		FS:                fsys,
	})
	if err != nil {
		os.Exit(1)
	}

	// Build the dependencies container.
	deps := cmd.Dependencies{
		Config:          cfg,
		Logger:          log,
		FS:              fsys,
		TemplateManager: tm,
	}

	// Create the root command and add subcommands.
	rootCmd := cmd.NewRootCmd(deps)
	rootCmd.AddCommand(cmd.NewConfigCmd(deps))
	rootCmd.AddCommand(cmd.NewZetCmd(deps))
	rootCmd.AddCommand(cmd.NewDayCmd(deps))
	rootCmd.AddCommand(cmd.NewTemplateCmd(deps))
	// (Add additional commands like day, zet, init, etc.)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
