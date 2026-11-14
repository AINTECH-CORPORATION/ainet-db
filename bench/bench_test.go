package bench

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func setupBenchDB(b *testing.B) *sql.DB {
	b.Helper()
	dbPath := filepath.Join(b.TempDir(), "ainet-bench.db")
	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)")
	if err != nil {
		b.Fatalf("open sqlite: %v", err)
	}
	b.Cleanup(func() { _ = db.Close() })
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS items (id INTEGER PRIMARY KEY AUTOINCREMENT, k TEXT, v TEXT)`); err != nil {
		b.Fatalf("create table: %v", err)
	}
	return db
}

func BenchmarkAINETDBWrite(b *testing.B) {
	db := setupBenchDB(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := db.Exec(`INSERT INTO items (k, v) VALUES (?, ?)`, "key", "ainet-db-val"); err != nil {
			b.Fatalf("insert: %v", err)
		}
	}
}

func BenchmarkAINETDBRead(b *testing.B) {
	db := setupBenchDB(b)
	for i := 0; i < 1000; i++ {
		_, _ = db.Exec(`INSERT INTO items (k, v) VALUES (?, ?)`, "key", "ainet-db-val")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM items`).Scan(&count); err != nil {
			b.Fatalf("select: %v", err)
		}
	}
}

func BenchmarkAINETDBVectorSearch(b *testing.B) {
	db := setupBenchDB(b)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS vec0_items (id TEXT PRIMARY KEY, vec BLOB NOT NULL)`)
	vec := make([]byte, 384*4)
	for i := 0; i < 100; i++ {
		_, _ = db.Exec(`INSERT OR REPLACE INTO vec0_items (id, vec) VALUES (?, ?)`, string(rune('a'+i)), vec)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`SELECT id, vec FROM vec0_items`)
		if err != nil {
			b.Fatalf("query: %v", err)
		}
		for rows.Next() {
		}
		_ = rows.Close()
	}
}
