package cmd

import (
	"fmt"
	"khanik/vidur"

	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the surang manager daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := vidur.StartDaemon(); err != nil {
				return fmt.Errorf("failed to start daemon: %w", err)
			}
			fmt.Println("Daemon started successfully.")
			return nil
		},
	}
}
