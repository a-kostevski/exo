package cmd

import (
	"fmt"
	"os"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	debug   bool
	quiet   bool
	verbose bool
	version bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "exo",
	Short: "Exo is a note-taking system",
	Long: `Exo is a note-taking system that helps you organize your thoughts, ideas, and knowledge.
	
It provides commands for:
- Creating and managing daily notes
- Creating and organizing Zettel notes
- Managing ideas and projects
- Customizing templates and configuration`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization for version command
		if cmd.Name() == "version" || version {
			fmt.Println("exo version 0.1.0")
			os.Exit(0)
		}

		// Initialize configuration
		if err := initConfig(); err != nil {
			return err
		}

		// Initialize logger
		if err := initLogger(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	// Add GNU-compliant flags
	flags := rootCmd.PersistentFlags()

	// Configuration file flag
	flags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/exo/config.yaml)")

	// Logging control flags
	flags.BoolVarP(&debug, "debug", "d", false, "enable debug logging")
	flags.BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	flags.BoolVarP(&quiet, "quiet", "q", false, "suppress all output except errors")

	//  Version flag
	flags.BoolVar(&version, "version", false, "print version information")

	// Add help flag (override default)
	flags.BoolP("help", "h", false, "display this help message")

	// 	// Update help template to be more GNU-like
	// 	rootCmd.SetHelpTemplate(`Usage: {{.CommandPath}} [options] <command> [<args>]
	//
	// {{.Short}}
	//
	// {{if .Long}}{{.Long}}
	//
	// {{end}}Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
	//   {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}
	//
	// Options:
	// {{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
	//
	// Use "{{.CommandPath}} <command> --help" for more information about a command.
	// `)

	// rootCmd.AddCommand(&cobra.Command{
	// 	Use:   "version",
	// 	Short: "Print version information",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("exo version 0.1.0")
	// 	},
	// })
}

func initConfig() error {
	if err := config.Initialize(cfgFile); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get configuration: %w", err)
	}

	// Set log level based on flags (with priority)
	if debug {
		cfg.Log.Level = logger.DebugLevel
	} else if verbose {
		cfg.Log.Level = logger.InfoLevel
	} else if quiet {
		cfg.Log.Level = logger.ErrorLevel
	}

	return nil
}

func initLogger() error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get configuration: %w", err)
	}

	// Reinitialize logger with new configuration
	if err := logger.Reinitialize(cfg.Log); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Log initial state
	logger.Debug("Debug logging enabled")
	logger.Info("Configuration loaded successfully")
	logger.Debugf("Using configuration: %+v", cfg)

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
