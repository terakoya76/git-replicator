package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RepoInfo struct {
	Host  string
	Owner string
	Repo  string
	Path  string
}

// List traverses the baseDir and returns a list of repositories found under the structure baseDir/host/owner/repo/base/.git
func List(ctx context.Context, baseDir string) ([]RepoInfo, error) {
	var repos []RepoInfo

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(baseDir, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." || rel == "" {
			return nil
		}
		walkParts := strings.Split(rel, string(filepath.Separator))
		if len(walkParts) == 4 && walkParts[3] == "base" {
			gitDir := filepath.Join(path, ".git")
			stat, err := os.Stat(gitDir)
			if err == nil && stat.IsDir() {
				repos = append(repos, RepoInfo{
					Host:  walkParts[0],
					Owner: walkParts[1],
					Repo:  walkParts[2],
					Path:  path,
				})
			}
			return filepath.SkipDir
		}
		// それ以外は無視
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}
	return repos, nil
}
