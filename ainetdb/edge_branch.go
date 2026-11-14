package ainetdb

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateEdgeBranch creates an instant isolated branch copy of dbPath (libsql-style).
func CreateEdgeBranch(dbPath, branchName string) (string, error) {
	if branchName == "" {
		return "", fmt.Errorf("branchName is required")
	}
	branchDir := filepath.Join(filepath.Dir(dbPath), ".branches")
	if err := os.MkdirAll(branchDir, 0o755); err != nil {
		return "", err
	}
	dstPath := filepath.Join(branchDir, fmt.Sprintf("%s.%s.db", filepath.Base(dbPath), branchName))
	if err := copyFile(dbPath, dstPath); err != nil {
		return "", err
	}
	if _, err := os.Stat(dbPath + "-wal"); err == nil {
		_ = copyFile(dbPath+"-wal", dstPath+"-wal")
	}
	return dstPath, nil
}
