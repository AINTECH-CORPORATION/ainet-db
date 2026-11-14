package ainetdb

import (
	"database/sql"
	"encoding/binary"
	"math"
	"sort"

	_ "modernc.org/sqlite"
)

// VectorMatch represents a nearest-neighbor match from vector-vec0.
type VectorMatch struct {
	ID    string
	Score float64
}

// VectorStore provides native SQLite-backed vector search (vector-vec0).
type VectorStore struct {
	db *sql.DB
}

// NewVectorStore opens or initializes a vector-vec0 table in the given SQLite DB.
func NewVectorStore(db *sql.DB) (*VectorStore, error) {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS vec0_items (id TEXT PRIMARY KEY, vec BLOB NOT NULL)`); err != nil {
		return nil, err
	}
	return &VectorStore{db: db}, nil
}

// Upsert stores a float32 embedding for id.
func (s *VectorStore) Upsert(id string, vec []float32) error {
	buf := make([]byte, len(vec)*4)
	for i, v := range vec {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(v))
	}
	_, err := s.db.Exec(`INSERT OR REPLACE INTO vec0_items (id, vec) VALUES (?, ?)`, id, buf)
	return err
}

// Search returns topK items ordered by cosine similarity descending.
func (s *VectorStore) Search(query []float32, topK int) ([]VectorMatch, error) {
	rows, err := s.db.Query(`SELECT id, vec FROM vec0_items`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []VectorMatch
	for rows.Next() {
		var id string
		var raw []byte
		if err := rows.Scan(&id, &raw); err != nil {
			return nil, err
		}
		matches = append(matches, VectorMatch{ID: id, Score: cosineBytes(query, raw)})
	}
	sort.Slice(matches, func(i, j int) bool { return matches[i].Score > matches[j].Score })
	if topK > 0 && len(matches) > topK {
		matches = matches[:topK]
	}
	return matches, rows.Err()
}

func cosineBytes(q []float32, raw []byte) float64 {
	n := len(raw) / 4
	if n != len(q) || n == 0 {
		return 0
	}
	var dot, nq, nv float64
	for i := 0; i < n; i++ {
		v := float64(math.Float32frombits(binary.LittleEndian.Uint32(raw[i*4:])))
		qf := float64(q[i])
		dot += qf * v
		nq += qf * qf
		nv += v * v
	}
	if nq == 0 || nv == 0 {
		return 0
	}
	return dot / (math.Sqrt(nq) * math.Sqrt(nv))
}
