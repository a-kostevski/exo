package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/fs"
	"github.com/a-kostevski/exo/internal/note/factory"
	"github.com/a-kostevski/exo/internal/templates"
)

var zetCmd = &cobra.Command{
	Use:   "zet [title]",
	Short: "Create a new Zettelkasten note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]

		cfg := config.MustGet()
		tm, err := templates.NewTemplateManager(cfg.TemplateDir)
		if err != nil {
			return fmt.Errorf("failed to create template manager: %w", err)
		}

		f := factory.New(tm)

		// Create Zettel note
		note, err := f.CreateZettel(title)
		if err != nil {
			return fmt.Errorf("failed to create zettel note: %w", err)
		}

		if err := note.Save(); err != nil {
			return fmt.Errorf("failed to save note: %w", err)
		}

		// Create/get today's daily note
		dailyNote, err := f.CreateDaily(time.Now())
		if err != nil {
			return fmt.Errorf("failed to create/get daily note: %w", err)
		}

		if err := dailyNote.Save(); err != nil {
			return fmt.Errorf("failed to save daily note: %w", err)
		}

		// Append Zettel link to daily note
		zettelLink := fmt.Sprintf("- [[%s]]", note.Title())

		if err := fs.AppendToFile(dailyNote.Path(), zettelLink); err != nil {
			return fmt.Errorf("failed to append link to daily note: %w", err)
		}

		if err := note.Open(); err != nil {
			return fmt.Errorf("failed to open note: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(zetCmd)
}
