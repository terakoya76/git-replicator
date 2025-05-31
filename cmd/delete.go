package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <branch>",
	Short: "Delete a branch directory under the current repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := args[0]
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
		if err := handlers.DeleteBranchDir(context.Background(), repoDir, branch); err != nil {
			return err
		}
		fmt.Printf("Deleted branch directory: %s\n", branch)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
