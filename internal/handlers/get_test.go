package handlers_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/assert"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

func getBaseDir(url string, rootDir string) string {
	u, _ := utils.ParseGitURL(url)
	return filepath.Join(rootDir, u.Host, u.Owner, u.Repo, "base")
}

func cleanupTestRepo(t *testing.T, url string, rootDir string) {
	dir := getBaseDir(url, rootDir)
	_ = os.RemoveAll(dir)
}

func TestGet(t *testing.T) {
	tmpDir := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatalf("failed to chdir in cleanup: %v", err)
		}
	})

	validRepoURL := "https://github.com/terakoya76/git-replicator-test"
	validRepoURLWithGit := "https://github.com/terakoya76/git-replicator-test.git"
	otherRepoURL := "https://github.com/terakoya76/git-replicator-dummy"

	t.Run("table-driven: valid/invalid url", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		tests := []struct {
			name    string
			url     string
			wantErr bool
		}{
			{"invalid url", "", true},
			{"invalid url 2", "http://", true},
			{"valid url", validRepoURL, false},
			{"valid url with .git", validRepoURLWithGit, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				cleanupTestRepo(t, validRepoURL, tmpDir)
				err := handlers.Get(context.Background(), tt.url, tmpDir)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					if !tt.wantErr && tt.url != "" {
						_, statErr := os.Stat(getBaseDir(tt.url, tmpDir))
						assert.NoError(t, statErr)
					}
				}
			})
		}
	})

	t.Run("non-git base subdirectory exists", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getBaseDir(validRepoURL, tmpDir)
		if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
			t.Fatalf("failed to mkdir parent dir: %v", err)
		}
		assert.NoError(t, os.Mkdir(dir, 0o755))
		err := handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.Error(t, err)
	})

	t.Run("non-git directory exists in subdirectory", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getBaseDir(validRepoURL, tmpDir)
		if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
			t.Fatalf("failed to mkdir parent dir: %v", err)
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to mkdir dir: %v", err)
		}
		f, err := os.Create(filepath.Join(dir, "dummy.txt"))
		assert.NoError(t, err)
		if err := f.Close(); err != nil {
			t.Errorf("failed to close file: %v", err)
		}
		err = handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.Error(t, err)
	})

	t.Run("git repo exists with same remote", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getBaseDir(validRepoURL, tmpDir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to mkdir dir: %v", err)
		}
		r, err := git.PlainInit(dir, false)
		assert.NoError(t, err)
		_, err = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{validRepoURLWithGit}})
		assert.NoError(t, err)
		indexPath := filepath.Join(dir, ".git", "index")
		f, err := os.Create(indexPath)
		assert.NoError(t, err)
		if err := f.Close(); err != nil {
			t.Errorf("failed to close file: %v", err)
		}
		err = handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.NoError(t, err)
	})

	t.Run("git repo exists with different remote", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getBaseDir(validRepoURL, tmpDir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("failed to mkdir dir: %v", err)
		}
		r, err := git.PlainInit(dir, false)
		assert.NoError(t, err)
		_, err = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{otherRepoURL}})
		assert.NoError(t, err)
		indexPath := filepath.Join(dir, ".git", "index")
		f, err := os.Create(indexPath)
		assert.NoError(t, err)
		if err := f.Close(); err != nil {
			t.Errorf("failed to close file: %v", err)
		}
		err = handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.Error(t, err)
	})
}
