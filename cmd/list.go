package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories under $HOME/git-replicator",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir, err := utils.GetGitReplicatorRoot()
		if err != nil {
			return fmt.Errorf("failed to get git-replicator root: %w", err)
		}
		repos, err := handlers.List(context.Background(), rootDir)
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
