package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetGitReplicatorRoot returns the path to $HOME/git-replicator
func GetGitReplicatorRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, "git-replicator"), nil
}

// FindRepoDir walks up from cwd to gitReplicatorRoot and returns the repo directory (the second-level directory under gitReplicatorRoot, e.g., $HOME/git-replicator/owner/repo).
func FindRepoDir(cwd, gitReplicatorRoot string) (string, error) {
	dir := cwd
	for {
		parent := filepath.Dir(dir)
		grandparent := filepath.Dir(parent)
		greatgrandparent := filepath.Dir(grandparent)
		if greatgrandparent == gitReplicatorRoot {
			return dir, nil
		}
		if dir == gitReplicatorRoot || dir == "/" || dir == "." {
			return "", fmt.Errorf("could not find repo directory under git-replicator root")
		}
		dir = parent
	}
}
