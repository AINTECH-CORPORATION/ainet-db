package ainetdb

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(exploreHTML))
	})
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
	mux.HandleFunc("/api/export", func(w http.ResponseWriter, r *http.Request) {
		table := r.URL.Query().Get("table")
		rows, err := db.Query(`SELECT * FROM "`+table+`" LIMIT 1000`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()
		cols, _ := rows.Columns()
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", table+".csv"))
		cw := csv.NewWriter(w)
		_ = cw.Write(cols)
		for rows.Next() {
			vals := make([]any, len(cols))
			ptrs := make([]any, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err := rows.Scan(ptrs...); err == nil {
				rec := make([]string, len(cols))
				for i, v := range vals {
					rec[i] = fmt.Sprint(v)
				}
				_ = cw.Write(rec)
			}
		}
		cw.Flush()
	})
	return mux, db.Close, nil
}

const exploreHTML = `<!doctype html>
<html><head><meta charset="utf-8"><title>AINET DB Explore</title>
<style>body{background:#0C1B3A;color:#EEF2F7;font-family:Inter,system-ui,sans-serif;padding:24px}h1{color:#1DBCB5}button{background:#12254D;color:#86CDFF;border:1px solid #263B6A;padding:8px 12px;margin:4px;cursor:pointer;border-radius:6px}table{width:100%;border-collapse:collapse;margin-top:16px;font-size:13px}th,td{border:1px solid #263B6A;padding:8px;text-align:left}th{background:#12254D;color:#86CDFF}</style>
</head><body><h1>AINET DB Explore</h1><div id="tables"></div><div id="grid"></div>
<script>
fetch('/api/tables').then(r=>r.json()).then(d=>{
  const el=document.getElementById('tables');
  (d.tables||[]).forEach(t=>{const b=document.createElement('button');b.textContent=t;b.onclick=()=>load(t);el.appendChild(b);});
});
function load(t){
  fetch('/api/rows?table='+encodeURIComponent(t)+'&limit=25').then(r=>r.json()).then(d=>{
    const rows=d.rows||[]; if(!rows.length){document.getElementById('grid').innerHTML='<p>Sin filas</p>';return;}
    const cols=Object.keys(rows[0]);
    document.getElementById('grid').innerHTML='<table><thead><tr>'+cols.map(c=>'<th>'+c+'</th>').join('')+'</tr></thead><tbody>'+rows.map(r=>'<tr>'+cols.map(c=>'<td>'+String(r[c]??'')+'</td>').join('')+'</tr>').join('')+'</tbody></table>';
  });
}
</script></body></html>`
