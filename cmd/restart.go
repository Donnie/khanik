package cmd

import (
	"fmt"
	"khanik/vidur"

	"github.com/spf13/cobra"
)

func newRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Restart the surang manager daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := vidur.RestartDaemon(); err != nil {
				return fmt.Errorf("failed to restart daemon: %w", err)
			}
			fmt.Println("Daemon restarted successfully.")
			return nil
		},
	}
}
