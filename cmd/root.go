package cmd

import (
	"fmt"
	"os"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/logger"
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
	Short: "A personal knowledge management system",
	Long: `Exo is a command-line tool for managing personal knowledge through
various note types and organizational structures.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check for version flag
		if version {
			fmt.Println("exo version 0.1.0")
			os.Exit(0)
		}

		// Set up initial logger with default settings
		initialLogConfig := logger.Config{
			Level:  "warn",
			Format: "text",
			Output: "stderr",
		}
		if err := logger.Initialize(initialLogConfig); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}

		if err := config.Initialize(cfgFile); err != nil {
			return fmt.Errorf("failed to initialize configuration: %w", err)
		}

		cfg, err := config.Get()
		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		// Update log level based on flags
		if debug {
			cfg.Log.Level = "debug"
		} else if verbose {
			cfg.Log.Level = "info"
		} else if quiet {
			cfg.Log.Level = "error"
		}

		// Re-initialize logger with final configuration
		if err := logger.Initialize(cfg.Log); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}

		logger.Debug("Configuration loaded successfully")
		logger.Debugf("Using configuration: %+v", cfg)

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add GNU-compliant flags
	flags := rootCmd.PersistentFlags()

	// Configuration file flag
	flags.StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/exo/config.yaml)")

	// Debug mode flag
	flags.BoolVarP(&debug, "debug", "d", false, "enable debug logging")

	// Add version flag
	flags.BoolVarP(&version, "version", "v", false, "print version information")

	// Add help flag (override default)
	flags.BoolP("help", "h", false, "display this help message")

	// Add quiet mode flag
	flags.BoolVarP(&quiet, "quiet", "q", false, "suppress all output except errors")

	// Add verbose mode flag
	flags.BoolVarP(&verbose, "verbose", "V", false, "verbose output")

	// Update help template to be more GNU-like
	rootCmd.SetHelpTemplate(`Usage: {{.CommandPath}} [options] <command> [<args>]

{{.Short}}

{{if .Long}}{{.Long}}

{{end}}Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Use "{{.CommandPath}} <command> --help" for more information about a command.
`)

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("exo version 0.1.0")
		},
	})
}
