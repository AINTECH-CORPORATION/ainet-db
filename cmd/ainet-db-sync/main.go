package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pocketbase/pocketbase/ainetdb"
)

func main() {
	dbPath := flag.String("db", "ainet.db", "Path to SQLite .db file")
	backupDir := flag.String("backups", "", "Snapshot directory (default: .snapshots)")
	branch := flag.String("branch", "auto-sync", "Edge branch name")
	interval := flag.Duration("interval", time.Minute, "Auto-sync interval")
	once := flag.Bool("once", false, "Run a single sync cycle and exit")
	flag.Parse()

	if *once {
		rep, err := ainetdb.RunSyncCycle(*dbPath, *backupDir, *branch)
		if err != nil {
			log.Fatalf("sync failed: %v", err)
		}
		log.Printf("sync ok: snap=%s branch=%s", rep.SnapshotPath, rep.BranchPath)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Printf("AINET DB AutoSync started (db=%s, interval=%s)", *dbPath, *interval)
	ainetdb.StartAutoSync(ctx, *dbPath, *backupDir, *branch, *interval, func(rep ainetdb.SyncReport, err error) {
		if err != nil {
			log.Printf("sync error: %v", err)
			return
		}
		log.Printf("sync ok: snap=%s branch=%s", rep.SnapshotPath, rep.BranchPath)
	})
}
