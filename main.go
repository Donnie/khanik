package main

import (
	"fmt"
	"os"
)

var Version = "dev"

func init() {
	rootCmd.AddCommand(
		startCmd,
		stopCmd,
		listCmd,
		versionCmd,
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
