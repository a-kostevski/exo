package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/pkg/templates"
)

// NewTemplateCmd creates a new "templates" command.
// By default, it lists available templates in the custom template directory.
// When the --install flag is provided, it installs the built-in default templates.
func NewTemplateCmd(deps Dependencies) *cobra.Command {
	var installFlag bool

	cmd := &cobra.Command{
		Use:   "templates",
		Short: "List available templates or install defaults",
		Long: `Manage templates.

By default, this command lists the available custom templates.
Use the --install flag to install built-in default templates into your custom template directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if installFlag {
				// Create a default template store using embedded default templates.
				defaultStore := templates.NewEmbedTemplateStore(templates.DefaultTemplatesFS, templates.DefaultTemplateBaseDir)
				opts := templates.InstallOptions{
					TargetDir: deps.Config.Dir.TemplateDir,
					Force:     false,
					Reader:    &defaultInputReader{},
				}
				if err := templates.InstallDefaultTemplates(templates.TemplateConfig{
					TemplateDir:       deps.Config.Dir.TemplateDir,
					TemplateExtension: ".md",
					FilePermissions:   0644,
					Logger:            deps.Logger,
					FS:                deps.FS,
				}, opts, defaultStore); err != nil {
					return fmt.Errorf("failed to install default templates: %w", err)
				}
				deps.Logger.Info("Default templates installed successfully")
				return nil
			}

			// Otherwise, list available templates.
			names, err := deps.TemplateManager.ListTemplates()
			if err != nil {
				return fmt.Errorf("failed to list templates: %w", err)
			}
			if len(names) == 0 {
				fmt.Println("No templates found")
				return nil
			}
			fmt.Println("Available templates:")
			for _, name := range names {
				customPath := filepath.Join(deps.Config.Dir.TemplateDir, name+".md")
				var source string
				if _, err := os.Stat(customPath); err == nil {
					source = "[Custom]"
				} else {
					source = "[Built-in]"
				}
				fmt.Printf("  - %s %s\n", source, name)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&installFlag, "install", "i", false, "Install default templates into the custom template directory")
	return cmd
}
