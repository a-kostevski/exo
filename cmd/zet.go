package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/pkg/zettel"
)

// NewZetCmd returns a new cobra.Command for the "zet" command.
func NewZetCmd(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zet [title]",
		Short: "Create a new Zettel note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			title := args[0]
			zNote, err := zettel.NewZettelNote(title, *deps.Config, deps.TemplateManager, deps.Logger, deps.FS)
			if err != nil {
				return fmt.Errorf("failed to create zettel note: %w", err)
			}
			if err := zNote.Save(); err != nil {
				return fmt.Errorf("failed to save zettel note: %w", err)
			}
			if err := zNote.Open(); err != nil {
				return fmt.Errorf("failed to open zettel note: %w", err)
			}
			return nil
		},
	}
	return cmd
}
