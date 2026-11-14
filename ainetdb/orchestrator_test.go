package ainetdb

import (
	"context"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
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

func TestStartAutoSync(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "autosync.db")
	ctx, cancel := context.WithCancel(context.Background())
	var count atomic.Int32
	done := make(chan struct{})
	go func() {
		StartAutoSync(ctx, dbPath, filepath.Join(tmp, "snaps"), "auto", 20*time.Millisecond, func(_ SyncReport, err error) {
			if err == nil && count.Add(1) >= 2 {
				cancel()
			}
		})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		cancel()
		t.Fatal("StartAutoSync did not stop in time")
	}
	if count.Load() < 2 {
		t.Fatalf("expected at least 2 cycles, got %d", count.Load())
	}
}
