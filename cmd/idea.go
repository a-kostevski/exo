package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/internal/note"
)

var ideaCmd = &cobra.Command{
	Use:   "idea [title]",
	Short: "Create and store an idea",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]

		idea, err := note.NewIdeaNote(title)
		if err != nil {
			return fmt.Errorf("failed to create idea: %w", err)
		}

		if err := idea.Save(); err != nil {
			return fmt.Errorf("failed to save idea: %w", err)
		}

		if err := idea.Open(); err != nil {
			return fmt.Errorf("failed to open idea: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ideaCmd)
}
