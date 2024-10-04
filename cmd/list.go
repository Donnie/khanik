package cmd

import (
	"khanik/vidur"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all running surangs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return vidur.ListSurangs()
		},
	}
}
