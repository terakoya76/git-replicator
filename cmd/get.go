package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
)

var getCmd = &cobra.Command{
	Use:   "get <url>",
	Short: "Clone a git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		url := args[0]
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		baseDir := filepath.Join(home, "git-replicator")
		if err := handlers.Get(ctx, url, baseDir); err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
