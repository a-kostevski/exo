package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/note/factory"
	"github.com/a-kostevski/exo/internal/templates"
)

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Create or open today's daily note",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.MustGet()
		tm, err := templates.NewTemplateManager(cfg.TemplateDir)
		if err != nil {
			return fmt.Errorf("failed to create template manager: %w", err)
		}

		f := factory.New(tm)
		note, err := f.CreateDaily(time.Now())
		if err != nil {
			return fmt.Errorf("failed to create daily note: %w", err)
		}

		if err := note.Save(); err != nil {
			return fmt.Errorf("failed to save note: %w", err)
		}

		if err := note.Open(); err != nil {
			return fmt.Errorf("failed to open note: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dayCmd)
}
