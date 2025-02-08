package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command using the injected dependencies.
// It provides GNU-friendly usage and help text.
func NewRootCmd(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exo [OPTIONS] <command> [arguments]",
		Short: "Exo is a note-taking system",
		Long: `Exo is a note-taking system that helps you organize your thoughts, ideas, and knowledge.

Usage:
  exo [OPTIONS] <command> [arguments]

Examples:
  exo init               Initialize exo configuration and directories.
  exo day                Open today's daily note.
  exo zet "My Note"      Create a new Zettel note with the title "My Note".

Global Options:
  -c, --config FILE      Specify configuration file (default: $HOME/.config/exo/config.yaml)
  -d, --debug            Enable debug logging (sets log level to "debug")
  -v, --verbose          Enable verbose output (sets log level to "info")
  -q, --quiet            Suppress all output except errors (sets log level to "error")
      --version          Print version information
  -h, --help             Show this help message and exit.
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Handle version flag.
			ver, err := cmd.Flags().GetBool("version")
			if err == nil && ver {
				fmt.Println("exo version 0.1.0")
				os.Exit(0)
			}
			// At this point, configuration and logger are already constructed.
			deps.Logger.Infof("Configuration loaded successfully: %+v", deps.Config)
			return nil
		},
	}

	// Define GNU-friendly persistent flags.
	flags := cmd.PersistentFlags()
	flags.StringP("config", "c", "", "Configuration file (default: $HOME/.config/exo/config.yaml)")
	flags.BoolP("debug", "d", false, "Enable debug logging (sets log level to 'debug')")
	flags.BoolP("verbose", "v", false, "Enable verbose output (sets log level to 'info')")
	flags.BoolP("quiet", "q", false, "Suppress all output except errors (sets log level to 'error')")
	flags.Bool("version", false, "Print version information")
	flags.BoolP("help", "h", false, "Show help message and exit")

	// Set a GNU-friendly help template.
	cmd.SetHelpTemplate(`Usage: {{.CommandPath}} [OPTIONS] <command> [arguments]

{{.Long}}

Global Options:
{{.PersistentFlags.FlagUsages | trimTrailingWhitespaces}}

Use "{{.CommandPath}} <command> --help" for more information about a command.
`)
	return cmd
}

// Execute runs the root command.
func Execute(deps Dependencies) {
	rootCmd := NewRootCmd(deps)
	// Subcommands will be added in main.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
