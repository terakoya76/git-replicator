package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

var switchCmd = &cobra.Command{
	Use:   "switch <branch>",
	Short: "Clone current repo into a new branch directory (like git switch)",
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
		opts := handlers.SwitchOptions{
			RepoDir:           repoDir,
			BranchName:        branch,
			GitReplicatorRoot: rootDir,
		}
		if err := handlers.Switch(context.Background(), opts, utils.DefaultGetRemoteURL, utils.DefaultCloneFunc, utils.DefaultSwitchBranchFunc); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
