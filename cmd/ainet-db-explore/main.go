package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase/ainetdb"
)

func main() {
	addr := flag.String("addr", ":8099", "HTTP listen address")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("usage: %s [--addr :8099] <archivo.db>", os.Args[0])
	}
	dbPath := flag.Arg(0)

	h, closeDB, err := ainetdb.NewExploreHandler(dbPath)
	if err != nil {
		log.Fatalf("open db %s: %v", dbPath, err)
	}
	defer closeDB()

	log.Printf("AINET DB Explore listening on http://localhost%s (db: %s)", *addr, dbPath)
	if err := http.ListenAndServe(*addr, h); err != nil {
		log.Fatal(err)
	}
}
