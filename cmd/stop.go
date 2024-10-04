package cmd

import (
	"fmt"
	"khanik/vidur"

	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop the surang manager daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := vidur.StopDaemon(); err != nil {
				return fmt.Errorf("failed to stop daemon: %w", err)
			}
			fmt.Println("Daemon stopped successfully.")
			return nil
		},
	}
}
