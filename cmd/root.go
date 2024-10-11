package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "khanik",
	Short: "Manage SSH surangs",
}

// Execute runs the root command.
func Execute(version string) error {
	addCommands(version)
	return rootCmd.Execute()
}

func addCommands(version string) {
	rootCmd.AddCommand(
		newStartCmd(),
		newStopCmd(),
		newRestartCmd(),
		newListCmd(),
		newVersionCmd(version),
	)
}

func init() {
	cobra.OnInitialize()
}
