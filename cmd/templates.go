package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/templates"
	"github.com/spf13/cobra"
)

var templateCmd = cobra.Command{
	Use:   "templates",
	Short: "Lists all templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get configuration
		cfg, err := config.Get()
		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		tm, err := templates.NewTemplateManager(cfg.Dir.TemplateDir)
		if err != nil {
			return fmt.Errorf("failed to create template manager: %w", err)
		}

		templates, err := tm.ListTemplates()
		if err != nil {
			return fmt.Errorf("failed to list templates: %w", err)
		}

		if len(templates) == 0 {
			fmt.Println("No templates found")
			return nil
		}

		fmt.Println("Available templates:")
		for _, name := range templates {
			customPath := filepath.Join(cfg.Dir.TemplateDir, name+".md")
			if _, err := os.Stat(customPath); err == nil {
				fmt.Printf("  - [Custom] ")
			} else {
				fmt.Printf("  - [Built-in] ")
			}
			fmt.Printf("%s\n%s\n", name, customPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(&templateCmd)
}
