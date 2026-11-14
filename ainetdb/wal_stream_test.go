package ainetdb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWALSnapshot(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "test.db")
	if err := os.WriteFile(dbPath, []byte("sqlite-mock"), 0o644); err != nil {
		t.Fatalf("write mock db: %v", err)
	}
	snapDir := filepath.Join(tmp, "snapshots")
	dst, err := WALSnapshot(dbPath, snapDir)
	if err != nil {
		t.Fatalf("WALSnapshot: %v", err)
	}
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("snapshot missing: %v", err)
	}
}
