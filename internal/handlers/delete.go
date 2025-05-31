package handlers

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/terakoya76/git-replicator/internal/utils"
)

// DeleteBranchDir deletes the branch directory under the given repo for a branch name.
func DeleteBranchDir(ctx context.Context, repoDir, branchName string) error {
	branchDir := filepath.Join(repoDir, branchName)
	if err := utils.RemoveDir(branchDir); err != nil {
		return fmt.Errorf("failed to delete branch directory %s: %w", branchDir, err)
	}
	return nil
}
