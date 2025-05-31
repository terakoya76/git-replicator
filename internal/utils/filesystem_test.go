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

func TestGetGitReplicatorRoot(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir() failed: %v", err)
	}

	want := filepath.Join(home, "git-replicator")
	got, err := utils.GetGitReplicatorRoot()
	if err != nil {
		t.Errorf("GetGitReplicatorRoot() error = %v, want nil", err)
	}
	if got != want {
		t.Errorf("GetGitReplicatorRoot() = %q, want %q", got, want)
	}
}

func TestRemoveDir(t *testing.T) {
	tmp := t.TempDir()

	dirs := []struct {
		name       string
		prepare    func() string // returns the path to remove
		wantErr    bool
		checkAfter func(path string) error
	}{
		{
			name: "remove directory with files and subdirs",
			prepare: func() string {
				dir := filepath.Join(tmp, "dir-with-files")
				sub := filepath.Join(dir, "subdir")
				if err := os.MkdirAll(sub, 0o755); err != nil {
					t.Fatalf("failed to create subdir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("data"), 0o644); err != nil {
					t.Fatalf("failed to create file1: %v", err)
				}
				return dir
			},
			wantErr: false,
			checkAfter: func(path string) error {
				_, err := os.Stat(path)
				if !os.IsNotExist(err) {
					return err
				}
				return nil
			},
		},
		{
			name: "remove non-existent directory",
			prepare: func() string {
				return filepath.Join(tmp, "not-exist-dir")
			},
			wantErr: false,
			checkAfter: func(path string) error {
				_, err := os.Stat(path)
				if !os.IsNotExist(err) {
					return err
				}
				return nil
			},
		},
		{
			name: "remove file (not directory)",
			prepare: func() string {
				file := filepath.Join(tmp, "file.txt")
				if err := os.WriteFile(file, []byte("data"), 0o644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
				return file
			},
			wantErr: false, // os.RemoveAll does not error on files
			checkAfter: func(path string) error {
				_, err := os.Stat(path)
				if !os.IsNotExist(err) {
					return err
				}
				return nil
			},
		},
	}

	for _, tt := range dirs {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.prepare()
			err := utils.RemoveDir(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkAfter(path); err != nil {
				t.Errorf("post-check failed: %v", err)
			}
		})
	}
}
