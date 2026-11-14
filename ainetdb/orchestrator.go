package ainetdb

import (
	"database/sql"
	"path/filepath"
)

// SyncReport summarizes one automated synchronization cycle across AINET DB modules.
type SyncReport struct {
	ExtensionsCount int
	SnapshotPath    string
	BranchPath      string
	VectorReady     bool
}

// RunSyncCycle executes an automated, synchronized cycle across vector-vec0, wal-stream and edge-branch.
func RunSyncCycle(dbPath, backupDir, branchName string) (SyncReport, error) {
	var rep SyncReport
	rep.ExtensionsCount = len(DefaultExtensions())

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return rep, err
	}
	if _, err := NewVectorStore(db); err != nil {
		_ = db.Close()
		return rep, err
	}
	rep.VectorReady = true
	_ = db.Close()

	if backupDir == "" {
		backupDir = filepath.Join(filepath.Dir(dbPath), ".snapshots")
	}
	snap, err := WALSnapshot(dbPath, backupDir)
	if err != nil {
		return rep, err
	}
	rep.SnapshotPath = snap

	if branchName != "" {
		br, err := CreateEdgeBranch(dbPath, branchName)
		if err != nil {
			return rep, err
		}
		rep.BranchPath = br
	}
	return rep, nil
}
