package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories under $HOME/git-replicator",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		baseDir := filepath.Join(home, "git-replicator")
		repos, err := handlers.List(context.Background(), baseDir)
		if err != nil {
			return fmt.Errorf("failed to list repositories: %w", err)
		}
		for _, repo := range repos {
			fmt.Printf("%s/%s/%s\n", repo.Host, repo.Owner, repo.Repo)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
