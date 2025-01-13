package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/internal/config"
	"github.com/a-kostevski/exo/internal/note"
	"github.com/a-kostevski/exo/internal/note/builtin"
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

		// Create Zettel note
		note, err := builtin.NewZettelNote(title, tm)
		if err != nil {
			return fmt.Errorf("failed to create zettel note: %w", err)
		}

		if err := note.Save(); err != nil {
			return fmt.Errorf("failed to save note: %w", err)
		}

		dailyNote, err := builtin.GetOrCreateTodayNote(tm)
		if err != nil {
			return fmt.Errorf("failed to get/create daily note: %w", err)
		}

		// Append Zettel reference to daily note
		if err := appendZettelReference(dailyNote, note); err != nil {
			return fmt.Errorf("failed to update daily note: %w", err)
		}

		if err := note.Open(); err != nil {
			return fmt.Errorf("failed to open note: %w", err)
		}

		return nil
	},
}

func appendZettelReference(daily *builtin.DailyNote, zettel note.Note) error {
	zettelLink := fmt.Sprintf("\n- [[%s]]", zettel.Title())

	content := daily.Content()
	content += zettelLink

	if err := daily.SetContent(content); err != nil {
		return fmt.Errorf("failed to update daily note content: %w", err)
	}

	return daily.Save()
}

func init() {
	rootCmd.AddCommand(zetCmd)
}
