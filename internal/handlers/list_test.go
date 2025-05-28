package handlers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	tmpDir := t.TempDir()

	repoPaths := []struct {
		host  string
		owner string
		repo  string
		dir   string
	}{
		{"github.com", "alice", "repo1", "base"},
		{"github.com", "bob", "repo2", "base"},
		{"gitlab.com", "carol", "repo3", "base"},
		{"github.com", "alice", "repoX", "base"},
		{"github.com", "alice", "repoX", "dir1"},
		{"github.com", "alice", "repoX", "dir2"},
	}
	for _, rp := range repoPaths {
		basePath := filepath.Join(tmpDir, rp.host, rp.owner, rp.repo, rp.dir)
		assert.NoError(t, os.MkdirAll(basePath, 0o755))
		if rp.dir == "base" {
			assert.NoError(t, os.MkdirAll(filepath.Join(basePath, ".git"), 0o755))
		}
	}

	// 2階層リポジトリ（無視されるべき）
	twoLevelPath := filepath.Join(tmpDir, "gitlab.com", "repo3")
	assert.NoError(t, os.MkdirAll(twoLevelPath, 0o755))

	wantBaseCount := 0
	for _, rp := range repoPaths {
		if rp.dir == "base" {
			wantBaseCount++
		}
	}

	tests := []struct {
		name    string
		baseDir string
		wantErr bool
		wantLen int
	}{
		{
			name:    "normal repositories",
			baseDir: tmpDir,
			wantErr: false,
			wantLen: wantBaseCount,
		},
		{
			name:    "not exist baseDir",
			baseDir: filepath.Join(tmpDir, "not-exist-dir"),
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "empty baseDir",
			baseDir: "",
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "only 2-level repo",
			baseDir: twoLevelPath,
			wantErr: false,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repos, err := List(context.Background(), tt.baseDir)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Len(t, repos, 0)
			} else {
				assert.NoError(t, err)
				if tt.wantLen == 0 {
					assert.Empty(t, repos)
				} else {
					assert.Len(t, repos, tt.wantLen)
					for _, rp := range repoPaths {
						if rp.dir != "base" {
							continue
						}
						found := false
						for _, r := range repos {
							if r.Host == rp.host && r.Owner == rp.owner && r.Repo == rp.repo {
								found = true
								break
							}
						}
						assert.True(t, found, "repo %s/%s/%s not found", rp.host, rp.owner, rp.repo)
					}
				}
			}
		})
	}
}
