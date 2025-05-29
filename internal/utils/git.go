package utils

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func DefaultCloneFunc(ctx context.Context, url, dir string) error {
	var err error
	if testing.Testing() {
		_, err = git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
			URL:      url,
			Progress: io.Discard,
		})
	} else {
		_, err = git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
	}
	return err
}

// Validate that repoDir is under gitReplicatorRoot and build the remote URL
func DefaultGetRemoteURL(repoDir, gitReplicatorRoot string) (string, error) {
	return BuildRemoteURLFromRepoDir(repoDir, gitReplicatorRoot)
}

// DefaultSwitchBranchFunc is the default implementation for switching branches, for external use
func DefaultSwitchBranchFunc(ctx context.Context, repoDir, branchName string) error {
	return SwitchBranch(ctx, repoDir, branchName)
}

type GitURLParts struct {
	Host  string
	Owner string
	Repo  string
}

func ParseGitURL(rawurl string) (GitURLParts, error) {
	var u GitURLParts
	urlStr := strings.TrimSuffix(rawurl, ".git")
	if strings.HasPrefix(urlStr, "git@") {
		parts := strings.SplitN(urlStr, ":", 2)
		if len(parts) != 2 {
			return u, fmt.Errorf("invalid ssh url: %s", rawurl)
		}
		host := strings.TrimPrefix(parts[0], "git@")
		pathParts := strings.Split(parts[1], "/")
		if len(pathParts) < 2 {
			return u, fmt.Errorf("invalid ssh url path: %s", rawurl)
		}
		u.Host = host
		u.Owner = pathParts[0]
		u.Repo = pathParts[1]
		return u, nil
	}
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return u, err
		}
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) < 2 {
			return u, fmt.Errorf("invalid url path: %s", rawurl)
		}
		u.Host = parsed.Host
		u.Owner = parts[0]
		u.Repo = parts[1]
		return u, nil
	}
	return u, fmt.Errorf("unsupported git url format: %s", rawurl)
}

func BuildRemoteURLFromRepoDir(repoDir, gitReplicatorRoot string) (string, error) {
	absRepo, err := filepath.Abs(repoDir)
	if err != nil {
		return "", fmt.Errorf("invalid repoDir: %s", repoDir)
	}
	absRoot, err := filepath.Abs(gitReplicatorRoot)
	if err != nil {
		return "", fmt.Errorf("invalid gitReplicatorRoot: %s", gitReplicatorRoot)
	}
	rel, err := filepath.Rel(absRoot, absRepo)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("repoDir is not under gitReplicatorRoot: %s", repoDir)
	}
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid repoDir: %s", repoDir)
	}
	host, owner, repo := parts[0], parts[1], parts[2]
	if host == "" || owner == "" || repo == "" || host == "." || owner == "." || repo == "." {
		return "", fmt.Errorf("invalid repoDir: %s", repoDir)
	}
	return fmt.Sprintf("https://%s/%s/%s.git", host, owner, repo), nil
}

// SwitchBranch performs the equivalent of 'git switch -C branchname' using go-git
func SwitchBranch(ctx context.Context, repoDir, branchName string) error {
	repo, err := git.PlainOpen(repoDir)
	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	}
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}
	refName := plumbing.NewBranchReferenceName(branchName)
	err = wt.Checkout(&git.CheckoutOptions{
		Branch: refName,
		Create: true,
		Force:  true, // Equivalent to -C
	})
	if err != nil {
		if err == git.ErrBranchExists || strings.Contains(err.Error(), "already exists") {
			err = wt.Checkout(&git.CheckoutOptions{
				Branch: refName,
				Create: false,
				Force:  true,
			})
		}
		if err != nil {
			return fmt.Errorf("failed to switch to branch %s: %w", branchName, err)
		}
	}
	return nil
}
