package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type SwitchOptions struct {
	RepoDir           string
	BranchName        string
	GitReplicatorRoot string
}

// GetRemoteURLFunc defines a function type for getting remote URL
// This allows for dependency injection in tests
type GetRemoteURLFunc func(repoDir, gitReplicatorRoot string) (string, error)

// CloneFunc defines a function type for cloning a repository
// This allows for dependency injection in tests
type CloneFunc func(ctx context.Context, url, dir string) error

// SwitchBranchFunc defines a function type for switching branches
// This allows for dependency injection in tests
type SwitchBranchFunc func(ctx context.Context, repoDir, branchName string) error

func Switch(
	ctx context.Context,
	opts SwitchOptions,
	getRemoteURL GetRemoteURLFunc,
	cloneFunc CloneFunc,
	switchBranchFunc SwitchBranchFunc,
) error {
	if opts.RepoDir == "" || opts.BranchName == "" {
		return fmt.Errorf("repo dir and branch name are required")
	}

	remoteURL, err := getRemoteURL(opts.RepoDir, opts.GitReplicatorRoot)
	if err != nil {
		return fmt.Errorf("failed to get remote url: %w", err)
	}

	branchDir := filepath.Join(opts.RepoDir, opts.BranchName)
	if _, err := os.Stat(branchDir); err == nil {
		return fmt.Errorf("branch directory already exists: %s", branchDir)
	}

	if err := cloneFunc(ctx, remoteURL, branchDir); err != nil {
		return fmt.Errorf("failed to clone to branch dir: %w", err)
	}

	if err := switchBranchFunc(ctx, branchDir, opts.BranchName); err != nil {
		return err
	}

	fmt.Printf("cloned branch: %s to dir: %s", opts.BranchName, branchDir)
	return nil
}
