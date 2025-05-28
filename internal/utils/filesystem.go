package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// MoveDir moves (renames) a directory from src to dst. If dst exists, it returns an error.
func MoveDir(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("destination already exists: %s", dst)
	}
	parent := filepath.Dir(dst)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}
	return os.Rename(src, dst)
}
