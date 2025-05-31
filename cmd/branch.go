package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "List branch directories under the current repository (like git switch)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		rootDir, err := utils.GetGitReplicatorRoot()
		if err != nil {
			return fmt.Errorf("failed to get git-replicator root: %w", err)
		}
		repoDir, err := utils.FindRepoDir(cwd, rootDir)
		if err != nil {
			return err
		}
		branches, err := handlers.ListBranchDirs(context.Background(), repoDir)
		if err != nil {
			return err
		}
		for _, b := range branches {
			fmt.Println(b)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
