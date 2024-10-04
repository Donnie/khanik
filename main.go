package main

import (
	"log"

	"khanik/cmd"
)

var Version = "dev"

func main() {
	if err := cmd.Execute(Version); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
