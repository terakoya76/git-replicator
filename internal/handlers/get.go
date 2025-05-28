package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/terakoya76/git-replicator/internal/utils"
)

func Get(ctx context.Context, url string, baseDir string) error {
	cloneURL := url
	if !strings.HasSuffix(url, ".git") {
		cloneURL = url + ".git"
	}

	// Parse the URL to extract host, owner, and repo
	u, err := utils.ParseGitURL(url)
	if err != nil {
		return fmt.Errorf("failed to parse git url: %w", err)
	}
	dir := filepath.Join(baseDir, u.Host, u.Owner, u.Repo, "base")
	gitIndex := filepath.Join(dir, ".git", "index")
	if _, err := os.Stat(dir); err == nil {
		// Directory exists, check if it's a git repo
		if _, err := os.Stat(gitIndex); err == nil {
			r, err := git.PlainOpen(dir)
			if err != nil {
				return fmt.Errorf("failed to open existing git repo: %w", err)
			}
			remotes, err := r.Remotes()
			if err != nil {
				return fmt.Errorf("failed to get remotes: %w", err)
			}
			for _, remote := range remotes {
				for _, u := range remote.Config().URLs {
					if u == cloneURL {
						// Same repo, do nothing
						return nil
					}
				}
			}
			return fmt.Errorf("directory %s exists and is a git repo, but remote does not match", dir)
		}
		return fmt.Errorf("directory %s exists but is not a git repo", dir)
	}
	// Directory does not exist, clone directly into the target directory
	_, err = git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:      cloneURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	return nil
}
