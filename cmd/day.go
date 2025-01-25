package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/internal/note"
)

var dayCmd = &cobra.Command{
	Use:   "day",
	Short: "Create or open today's daily note",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get or create today's note
		daily, err := note.GetOrCreateTodayNote()
		if err != nil {
			return fmt.Errorf("failed to get/create daily note: %w", err)
		}

		if err := daily.Open(); err != nil {
			return fmt.Errorf("failed to open daily note: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dayCmd)
}
