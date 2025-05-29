package utils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terakoya76/git-replicator/internal/utils"
)

func TestFindRepoDir(t *testing.T) {
	tmp := t.TempDir()
	// Simulate $HOME/git-replicator/owner/repo/base
	gitReplicatorRoot := filepath.Join(tmp, "git-replicator")
	hostDir := filepath.Join(gitReplicatorRoot, "github.com")
	ownerDir := filepath.Join(hostDir, "owner")
	repoDir := filepath.Join(ownerDir, "repo")
	baseDir := filepath.Join(repoDir, "base")
	subDir := filepath.Join(baseDir, "subdir")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("failed to mkdir subDir: %v", err)
	}

	branchDir := filepath.Join(repoDir, "branch")
	if err := os.MkdirAll(branchDir, 0o755); err != nil {
		t.Fatalf("failed to mkdir branchDir: %v", err)
	}

	tests := []struct {
		name    string
		cwd     string
		root    string
		want    string
		wantErr bool
	}{
		{
			name: "in repo dir",
			cwd:  repoDir,
			root: gitReplicatorRoot,
			want: repoDir,
		},
		{
			name: "in base dir",
			cwd:  baseDir,
			root: gitReplicatorRoot,
			want: repoDir,
		},
		{
			name: "in subdir under base",
			cwd:  subDir,
			root: gitReplicatorRoot,
			want: repoDir,
		},
		{
			name: "in branch dir",
			cwd:  branchDir,
			root: gitReplicatorRoot,
			want: repoDir,
		},
		{
			name:    "in owner dir (not repo)",
			cwd:     ownerDir,
			root:    gitReplicatorRoot,
			wantErr: true,
		},
		{
			name:    "in host dir (not repo)",
			cwd:     hostDir,
			root:    gitReplicatorRoot,
			wantErr: true,
		},
		{
			name:    "at gitReplicatorRoot",
			cwd:     gitReplicatorRoot,
			root:    gitReplicatorRoot,
			wantErr: true,
		},
		{
			name:    "at root dir",
			cwd:     "/",
			root:    gitReplicatorRoot,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := utils.FindRepoDir(tt.cwd, tt.root)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
