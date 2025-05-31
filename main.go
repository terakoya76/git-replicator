package main

import (
	"os"

	"github.com/terakoya76/git-replicator/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
