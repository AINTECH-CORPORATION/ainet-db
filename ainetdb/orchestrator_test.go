package ainetdb

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOrchestratorSyncCycle(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "sync.db")
	rep, err := RunSyncCycle(dbPath, filepath.Join(tmp, "snaps"), "auto-sync")
	if err != nil {
		t.Fatalf("RunSyncCycle: %v", err)
	}
	if rep.ExtensionsCount != 4 || !rep.VectorReady {
		t.Fatalf("unexpected report: %+v", rep)
	}
	if _, err := os.Stat(rep.SnapshotPath); err != nil {
		t.Fatalf("snapshot missing: %v", err)
	}
	if _, err := os.Stat(rep.BranchPath); err != nil {
		t.Fatalf("branch missing: %v", err)
	}
}
