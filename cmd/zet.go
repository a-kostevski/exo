package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/internal/note"
)

var zetCmd = &cobra.Command{
	Use:   "zet [title]",
	Short: "Create a new Zettel note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]

		zettel, err := note.NewZetNote(title)
		if err != nil {
			return fmt.Errorf("failed to create zettel: %w", err)
		}

		if err := zettel.Save(); err != nil {
			return fmt.Errorf("failed to save zettel: %w", err)
		}

		if err := zettel.Open(); err != nil {
			return fmt.Errorf("failed to open zettel: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(zetCmd)
}
