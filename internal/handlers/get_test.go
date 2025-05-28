package handlers_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/assert"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

func getDirName(url string) string {
	return filepath.Base(strings.TrimSuffix(url, ".git"))
}

func getCloneDir(url string, baseDir string) string {
	u, _ := utils.ParseGitURL(url)
	return filepath.Join(baseDir, u.Host, u.Owner, u.Repo, "base")
}

func cleanupTestRepo(t *testing.T, url string, baseDir string) {
	dir := getCloneDir(url, baseDir)
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
	t.Cleanup(func() { os.Chdir(cwd) })

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
						_, statErr := os.Stat(getCloneDir(tt.url, tmpDir))
						assert.NoError(t, statErr)
					}
				}
			})
		}
	})

	t.Run("non-git base subdirectory exists", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getCloneDir(validRepoURL, tmpDir)
		os.MkdirAll(filepath.Dir(dir), 0o755)
		assert.NoError(t, os.Mkdir(dir, 0o755))
		err := handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.Error(t, err)
	})

	t.Run("non-git directory exists in subdirectory", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getCloneDir(validRepoURL, tmpDir)
		os.MkdirAll(filepath.Dir(dir), 0o755)
		assert.NoError(t, os.Mkdir(dir, 0o755))
		f, err := os.Create(filepath.Join(dir, "dummy.txt"))
		assert.NoError(t, err)
		f.Close()
		err = handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.Error(t, err)
	})

	t.Run("git repo exists with same remote", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getCloneDir(validRepoURL, tmpDir)
		os.MkdirAll(dir, 0o755)
		r, err := git.PlainInit(dir, false)
		assert.NoError(t, err)
		_, err = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{validRepoURLWithGit}})
		assert.NoError(t, err)
		indexPath := filepath.Join(dir, ".git", "index")
		f, err := os.Create(indexPath)
		assert.NoError(t, err)
		f.Close()
		err = handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.NoError(t, err)
	})

	t.Run("git repo exists with different remote", func(t *testing.T) {
		cleanupTestRepo(t, validRepoURL, tmpDir)
		dir := getCloneDir(validRepoURL, tmpDir)
		os.MkdirAll(dir, 0o755)
		r, err := git.PlainInit(dir, false)
		assert.NoError(t, err)
		_, err = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{otherRepoURL}})
		assert.NoError(t, err)
		indexPath := filepath.Join(dir, ".git", "index")
		f, err := os.Create(indexPath)
		assert.NoError(t, err)
		f.Close()
		err = handlers.Get(context.Background(), validRepoURL, tmpDir)
		assert.Error(t, err)
	})
}
