package handlers

import (
	"context"
	"fmt"
	"os"
)

// ListBranchDirs returns a list of branch directory names under the given repoDir (including 'base').
func ListBranchDirs(ctx context.Context, repoDir string) ([]string, error) {
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read repo directory: %w", err)
	}
	var branches []string
	for _, entry := range entries {
		if entry.IsDir() {
			branches = append(branches, entry.Name())
		}
	}
	return branches, nil
}
