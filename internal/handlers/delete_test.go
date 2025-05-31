package handlers_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/terakoya76/git-replicator/internal/handlers"
)

func TestDeleteBranchDir(t *testing.T) {
	tmpDir := t.TempDir()
	repoDir := filepath.Join(tmpDir, "repo")
	branchName := "feature-x"
	branchDir := filepath.Join(repoDir, branchName)

	tests := []struct {
		name       string
		branch     string
		prepare    func()
		wantErr    bool
		checkAfter func(branchDir string) error
	}{
		{
			name:   "delete existing branch dir",
			branch: branchName,
			prepare: func() {
				if err := os.MkdirAll(branchDir, 0o755); err != nil {
					t.Fatalf("failed to create branch dir: %v", err)
				}
				filePath := filepath.Join(branchDir, "dummy.txt")
				if err := os.WriteFile(filePath, []byte("dummy"), 0o644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			},
			wantErr: false,
			checkAfter: func(branchDir string) error {
				if _, err := os.Stat(branchDir); !os.IsNotExist(err) {
					return err
				}
				return nil
			},
		},
		{
			name:   "delete already deleted branch dir",
			branch: branchName,
			prepare: func() {
				_ = os.RemoveAll(branchDir)
			},
			wantErr: false,
			checkAfter: func(branchDir string) error {
				if _, err := os.Stat(branchDir); !os.IsNotExist(err) {
					return err
				}
				return nil
			},
		},
		{
			name:    "delete non-existent branch dir",
			branch:  "nonexistent",
			prepare: func() {},
			wantErr: false,
			checkAfter: func(branchDir string) error {
				if _, err := os.Stat(branchDir); !os.IsNotExist(err) {
					return err
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			branchDirPath := filepath.Join(repoDir, tt.branch)
			err := handlers.DeleteBranchDir(context.Background(), repoDir, tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteBranchDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := tt.checkAfter(branchDirPath); err != nil {
				t.Errorf("post-check failed: %v", err)
			}
		})
	}
}
