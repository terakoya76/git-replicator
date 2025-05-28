package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

var getCmd = &cobra.Command{
	Use:   "get <url>",
	Short: "Clone a git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		url := args[0]
		rootDir, err := utils.GetGitReplicatorRoot()
		if err != nil {
			return fmt.Errorf("failed to get git-replicator root: %w", err)
		}
		if err := handlers.Get(ctx, url, rootDir); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
