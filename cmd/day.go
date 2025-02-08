package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/a-kostevski/exo/pkg/periodic"
)

// NewDayCmd returns a new cobra.Command for the "day" command.
func NewDayCmd(deps Dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "day",
		Short: "Create or open today's daily note",
		RunE: func(cmd *cobra.Command, args []string) error {
			today := time.Now().Truncate(24 * time.Hour)
			// Create (or load) today's daily note using injected dependencies.
			daily, err := periodic.NewDailyNote(today, *deps.Config, deps.TemplateManager, deps.Logger, deps.FS)
			if err != nil {
				return fmt.Errorf("failed to create daily note: %w", err)
			}
			if err := daily.Open(); err != nil {
				return fmt.Errorf("failed to open daily note: %w", err)
			}
			return nil
		},
	}
	return cmd
}
