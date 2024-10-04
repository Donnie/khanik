package main

import (
	"fmt"
	"khanik/vidur"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "khanik",
	Short: "Manage SSH surangs",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the surang manager daemon",
	Run: func(cmd *cobra.Command, args []string) {
		vidur.StartDaemon()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the surang manager daemon",
	Run: func(cmd *cobra.Command, args []string) {
		vidur.StopDaemon()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running surangs",
	Run: func(cmd *cobra.Command, args []string) {
		vidur.ListSurangs()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
