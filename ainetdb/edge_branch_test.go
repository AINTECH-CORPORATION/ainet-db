package ainetdb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEdgeBranch(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "main.db")
	if err := os.WriteFile(dbPath, []byte("sqlite-main"), 0o644); err != nil {
		t.Fatalf("write mock db: %v", err)
	}
	branchPath, err := CreateEdgeBranch(dbPath, "feature-x")
	if err != nil {
		t.Fatalf("CreateEdgeBranch: %v", err)
	}
	if _, err := os.Stat(branchPath); err != nil {
		t.Fatalf("branch missing: %v", err)
	}
}
