package handlers_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/terakoya76/git-replicator/internal/handlers"
)

func TestListBranchDirs(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("failed to create repo dir: %v", err)
	}
	// Create base and branch directories
	if err := os.Mkdir(filepath.Join(repoDir, "base"), 0o755); err != nil {
		t.Fatalf("failed to create base dir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(repoDir, "feature-x"), 0o755); err != nil {
		t.Fatalf("failed to create feature-x dir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(repoDir, "bugfix-y"), 0o755); err != nil {
		t.Fatalf("failed to create bugfix-y dir: %v", err)
	}

	t.Run("list branch dirs (including base)", func(t *testing.T) {
		branches, err := handlers.ListBranchDirs(context.Background(), repoDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := map[string]bool{"feature-x": true, "bugfix-y": true, "base": true}
		if len(branches) != 3 {
			t.Errorf("expected 3 branches, got %d", len(branches))
		}
		for _, b := range branches {
			if !want[b] {
				t.Errorf("unexpected branch: %s", b)
			}
		}
	})

	t.Run("repo dir does not exist", func(t *testing.T) {
		_, err := handlers.ListBranchDirs(context.Background(), filepath.Join(tmpDir, "not-exist"))
		if err == nil {
			t.Errorf("expected error for non-existent dir, got nil")
		}
	})
}
