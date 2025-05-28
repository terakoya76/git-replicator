package handlers_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/terakoya76/git-replicator/internal/handlers"
	"github.com/terakoya76/git-replicator/internal/utils"
)

func TestSwitch(t *testing.T) {
	tmpDir := t.TempDir()
	gitReplicatorRoot := filepath.Join(tmpDir, "git-replicator")
	repoDir := filepath.Join(gitReplicatorRoot, "github.com", "terakoya76", "git-replicator-test")
	baseDir := filepath.Join(repoDir, "base")
	branchName := "feature-x"
	branchDir := filepath.Join(repoDir, branchName)

	tests := []struct {
		name             string
		prepare          func(tmpDir, baseDir, branchDir string)
		getRemoteURL     handlers.GetRemoteURLFunc
		cloneFunc        handlers.CloneFunc
		switchBranchFunc handlers.SwitchBranchFunc
		err              error
		check            func(t *testing.T, tmpDir, baseDir, branchDir string, oldwd string)
	}{
		{
			name: "success (actual clone)",
			prepare: func(tmpDir, baseDir, branchDir string) {
				_ = utils.DefaultCloneFunc(context.Background(), "https://github.com/terakoya76/git-replicator-test", baseDir)
			},
			getRemoteURL:     utils.DefaultGetRemoteURL,
			cloneFunc:        utils.DefaultCloneFunc,
			switchBranchFunc: utils.DefaultSwitchBranchFunc,
			err:              nil,
			check: func(t *testing.T, tmpDir, baseDir, branchDir string, oldwd string) {
				// Assert that branchDir exists
				_, statErr := os.Stat(branchDir)
				assert.NoError(t, statErr)
				// Assert that HEAD branch is branchName
				repo, err := git.PlainOpen(branchDir)
				assert.NoError(t, err)
				head, err := repo.Head()
				assert.NoError(t, err)
				assert.Equal(t, "refs/heads/"+branchName, head.Name().String())
			},
		},
		{
			name: "already existing branch dir",
			prepare: func(tmpDir, baseDir, branchDir string) {
				_ = utils.DefaultCloneFunc(context.Background(), "https://github.com/terakoya76/git-replicator-test", baseDir)
				_ = utils.DefaultCloneFunc(context.Background(), "https://github.com/terakoya76/git-replicator-test", branchDir)
			},
			getRemoteURL:     utils.DefaultGetRemoteURL,
			cloneFunc:        utils.DefaultCloneFunc,
			switchBranchFunc: utils.DefaultSwitchBranchFunc,
			err:              errors.New("branch directory already exists: " + branchDir),
			check: func(t *testing.T, tmpDir, baseDir, branchDir string, oldwd string) {
				// Assert that branchDir exists
				_, statErr := os.Stat(branchDir)
				assert.NoError(t, statErr)
			},
		},
		{
			name: "remote fetch failure",
			prepare: func(tmpDir, baseDir, branchDir string) {
				_ = utils.DefaultCloneFunc(context.Background(), "https://github.com/terakoya76/git-replicator-test", baseDir)
			},
			getRemoteURL: func(base string, _ string) (string, error) {
				return "", errors.New("remote error")
			},
			cloneFunc:        utils.DefaultCloneFunc,
			switchBranchFunc: utils.DefaultSwitchBranchFunc,
			err:              errors.New("failed to get remote url: remote error"),
			check: func(t *testing.T, tmpDir, baseDir, branchDir string, oldwd string) {
				// Assert that branchDir does not exist
				_, statErr := os.Stat(branchDir)
				assert.Error(t, statErr)
			},
		},
		{
			name: "clone failure",
			prepare: func(tmpDir, baseDir, branchDir string) {
				_ = utils.DefaultCloneFunc(context.Background(), "https://github.com/terakoya76/git-replicator-test", baseDir)
			},
			getRemoteURL: utils.DefaultGetRemoteURL,
			cloneFunc: func(ctx context.Context, url, dir string) error {
				return errors.New("clone error")
			},
			switchBranchFunc: utils.DefaultSwitchBranchFunc,
			err:              errors.New("failed to clone to branch dir: clone error"),
			check: func(t *testing.T, tmpDir, baseDir, branchDir string, oldwd string) {
				// Assert that branchDir does not exist
				_, statErr := os.Stat(branchDir)
				assert.Error(t, statErr)
			},
		},
		{
			name: "git switch failure",
			prepare: func(tmpDir, baseDir, branchDir string) {
				_ = utils.DefaultCloneFunc(context.Background(), "https://github.com/terakoya76/git-replicator-test", baseDir)
			},
			getRemoteURL: utils.DefaultGetRemoteURL,
			cloneFunc:    utils.DefaultCloneFunc,
			switchBranchFunc: func(ctx context.Context, repoDir, branchName string) error {
				return errors.New("mock switch branch error")
			},
			err: errors.New("mock switch branch error"),
			check: func(t *testing.T, tmpDir, baseDir, branchDir string, oldwd string) {
				// Assert that branchDir exists
				_, statErr := os.Stat(branchDir)
				assert.NoError(t, statErr)
				// Assert that HEAD branch is not branchName
				repo, err := git.PlainOpen(branchDir)
				assert.NoError(t, err)
				head, err := repo.Head()
				assert.NoError(t, err)
				assert.NotEqual(t, "refs/heads/"+branchName, head.Name().String())
			},
		},
	}

	os.Chdir(tmpDir)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldwd, _ := os.Getwd()
			defer os.Chdir(oldwd)
			if tt.prepare != nil {
				tt.prepare(tmpDir, baseDir, branchDir)
			}
			opts := handlers.SwitchOptions{
				RepoDir:           repoDir,
				BranchName:        branchName,
				GitReplicatorRoot: gitReplicatorRoot,
			}
			err := handlers.Switch(context.Background(), opts, tt.getRemoteURL, tt.cloneFunc, tt.switchBranchFunc)
			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			if tt.check != nil {
				tt.check(t, tmpDir, baseDir, branchDir, oldwd)
			}

			// cleanup
			os.RemoveAll(baseDir)
			os.RemoveAll(branchDir)
		})
	}
}
