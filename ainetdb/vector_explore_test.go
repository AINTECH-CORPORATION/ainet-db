package ainetdb

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestVectorIndexAndSearch(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "vec.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	vs, err := NewVectorStore(db)
	if err != nil {
		t.Fatalf("new vector store: %v", err)
	}
	_ = vs.Upsert("doc-a", []float32{1, 0, 0})
	_ = vs.Upsert("doc-b", []float32{0, 1, 0})

	matches, err := vs.Search([]float32{0.9, 0.1, 0}, 1)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(matches) != 1 || matches[0].ID != "doc-a" {
		t.Fatalf("unexpected matches: %+v", matches)
	}
}

func TestExploreDB(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "explore.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	_, _ = db.Exec(`CREATE TABLE sessions (id TEXT, title TEXT); INSERT INTO sessions VALUES ('s1', 'AINET');`)
	_ = db.Close()

	h, closeDB, err := NewExploreHandler(dbPath)
	if err != nil {
		t.Fatalf("new explore handler: %v", err)
	}
	defer closeDB()
	req := httptest.NewRequest(http.MethodGet, "/api/tables", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var payload map[string][]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(payload["tables"]) != 1 || payload["tables"][0] != "sessions" {
		t.Fatalf("unexpected tables: %+v", payload)
	}
}
