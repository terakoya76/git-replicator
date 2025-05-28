package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/terakoya76/git-replicator/internal/handlers"
)

var (
	sourceRepo string
	targetRepo string
)

var replicateCmd = &cobra.Command{
	Use:   "replicate",
	Short: "Replicate a git repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		opts := handlers.ReplicateOptions{
			SourceRepo: sourceRepo,
			TargetRepo: targetRepo,
		}
		if err := handlers.Replicate(ctx, opts); err != nil {
			return fmt.Errorf("failed to replicate: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(replicateCmd)
	replicateCmd.Flags().StringVar(&sourceRepo, "source", "", "Source git repository URL")
	replicateCmd.Flags().StringVar(&targetRepo, "target", "", "Target git repository URL")
	replicateCmd.MarkFlagRequired("source")
	replicateCmd.MarkFlagRequired("target")
}
