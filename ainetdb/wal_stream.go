package ainetdb

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// WALSnapshot copies the SQLite database (and -wal if present) to dstDir (Litestream-style).
func WALSnapshot(dbPath, dstDir string) (string, error) {
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return "", err
	}
	stamp := time.Now().UTC().Format("20060102T150405Z")
	dstPath := filepath.Join(dstDir, fmt.Sprintf("%s.%s.bak", filepath.Base(dbPath), stamp))
	if err := copyFile(dbPath, dstPath); err != nil {
		return "", err
	}
	walSrc := dbPath + "-wal"
	if _, err := os.Stat(walSrc); err == nil {
		_ = copyFile(walSrc, dstPath+"-wal")
	}
	return dstPath, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
