package ainetdb

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "modernc.org/sqlite"
)

// NewExploreHandler returns an HTTP handler (Datasette-style) for any .db file and a closer.
func NewExploreHandler(dbPath string) (http.Handler, func() error, error) {
	db, err := sql.Open("sqlite", dbPath+"?mode=ro")
	if err != nil {
		return nil, nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tables", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var tables []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err == nil {
				tables = append(tables, name)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"tables": tables})
	})
	mux.HandleFunc("/api/rows", func(w http.ResponseWriter, r *http.Request) {
		table := r.URL.Query().Get("table")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit <= 0 || limit > 200 {
			limit = 50
		}
		rows, err := db.Query(`SELECT * FROM "`+table+`" LIMIT ?`, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		cols, _ := rows.Columns()
		var out []map[string]any
		for rows.Next() {
			vals := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err := rows.Scan(ptrs...); err == nil {
				m := make(map[string]any, len(cols))
				for i, c := range cols {
					m[c] = vals[i]
				}
				out = append(out, m)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"table": table, "rows": out})
	})
	return mux, db.Close, nil
}
