package utils_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/terakoya76/git-replicator/internal/utils"
)

func TestParseGitURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    utils.GitURLParts
		wantErr bool
	}{
		{
			name:  "https url",
			input: "https://github.com/owner/repo",
			want:  utils.GitURLParts{Host: "github.com", Owner: "owner", Repo: "repo"},
		},
		{
			name:  "https url with .git",
			input: "https://github.com/owner/repo.git",
			want:  utils.GitURLParts{Host: "github.com", Owner: "owner", Repo: "repo"},
		},
		{
			name:  "ssh url",
			input: "git@github.com:owner/repo",
			want:  utils.GitURLParts{Host: "github.com", Owner: "owner", Repo: "repo"},
		},
		{
			name:  "ssh url with .git",
			input: "git@github.com:owner/repo.git",
			want:  utils.GitURLParts{Host: "github.com", Owner: "owner", Repo: "repo"},
		},
		{
			name:    "invalid ssh url",
			input:   "git@github.com:owner",
			wantErr: true,
		},
		{
			name:    "invalid https url",
			input:   "https://github.com/owner",
			wantErr: true,
		},
		{
			name:    "unsupported format",
			input:   "ftp://github.com/owner/repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.ParseGitURL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestBuildRemoteURLFromRepoDir(t *testing.T) {
	tmp := t.TempDir()
	// Simulate $HOME/git-replicator/owner/repo/base
	gitReplicatorRoot := filepath.Join(tmp, "git-replicator")
	host := "github.com"
	owner := "owner"
	repo := "repo"
	repoDir := filepath.Join(gitReplicatorRoot, host, owner, repo)
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("failed to mkdir repoDir: %v", err)
	}

	tests := []struct {
		name    string
		repoDir string
		want    string
		wantErr bool
	}{
		{
			name:    "normal case",
			repoDir: repoDir,
			want:    "https://github.com/owner/repo.git",
		},
		{
			name:    "missing repo",
			repoDir: filepath.Join(gitReplicatorRoot, host, owner, ""),
			wantErr: true,
		},
		{
			name:    "missing owner",
			repoDir: filepath.Join(gitReplicatorRoot, host, "", repo),
			wantErr: true,
		},
		{
			name:    "missing host",
			repoDir: filepath.Join(gitReplicatorRoot, "", owner, repo),
			wantErr: true,
		},
		{
			name:    "root dir",
			repoDir: "/",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.BuildRemoteURLFromRepoDir(tt.repoDir, gitReplicatorRoot)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSwitchBranch(t *testing.T) {
	tests := []struct {
		name             string
		beforeBranchName string
		afterBranchName  string
		branchName       string
		wantErr          bool
	}{
		{
			name:             "switch from main to feature/test-branch",
			beforeBranchName: "main",
			afterBranchName:  "feature/test-branch",
			branchName:       "feature/test-branch",
			wantErr:          false,
		},
		{
			name:             "switch from feature/test-branch to feature/test-branch (already on branch)",
			beforeBranchName: "feature/test-branch",
			afterBranchName:  "feature/test-branch",
			branchName:       "feature/test-branch",
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare
			tmp := t.TempDir()
			cloneDir := filepath.Join(tmp, "clone")
			assert.NoError(t, utils.DefaultCloneFunc(context.Background(), "https://github.com/terakoya76/git-replicator-test", cloneDir))

			if tt.beforeBranchName != "main" {
				repo, err := git.PlainOpen(cloneDir)
				assert.NoError(t, err)
				wt, err := repo.Worktree()
				assert.NoError(t, err)
				refName := plumbing.NewBranchReferenceName(tt.beforeBranchName)
				err = wt.Checkout(&git.CheckoutOptions{
					Branch: refName,
					Create: true,
					Force:  true,
				})
				assert.NoError(t, err)
			}

			err := utils.SwitchBranch(context.TODO(), cloneDir, tt.branchName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				repo, err := git.PlainOpen(cloneDir)
				assert.NoError(t, err)
				head, err := repo.Head()
				assert.NoError(t, err)
				assert.Equal(t, "refs/heads/"+tt.afterBranchName, head.Name().String())
			}
		})
	}
}
