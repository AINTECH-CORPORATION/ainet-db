package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/pocketbase/pocketbase/ainetdb"
)

func main() {
	addr := flag.String("addr", ":8099", "HTTP listen address")
	idle := flag.Duration("idle-timeout", 5*time.Minute, "Auto-shutdown after inactivity (0 disables)")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatalf("usage: %s [--addr :8099] [--idle-timeout 5m] <archivo.db>", os.Args[0])
	}
	dbPath := flag.Arg(0)

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("port %s busy (no collision): %v", *addr, err)
	}

	h, closeDB, err := ainetdb.NewExploreHandler(dbPath)
	if err != nil {
		_ = ln.Close()
		log.Fatalf("open db %s: %v", dbPath, err)
	}
	defer closeDB()

	var lastHit atomic.Int64
	lastHit.Store(time.Now().UnixNano())
	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lastHit.Store(time.Now().UnixNano())
		h.ServeHTTP(w, r)
	})

	srv := &http.Server{Handler: wrapped}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if *idle > 0 {
		go func() {
			t := time.NewTicker(15 * time.Second)
			defer t.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					if time.Since(time.Unix(0, lastHit.Load())) >= *idle {
						log.Printf("idle timeout (%s) reached, shutting down cleanly", *idle)
						_ = srv.Shutdown(context.Background())
						return
					}
				}
			}
		}()
	}

	go func() {
		<-ctx.Done()
		_ = srv.Shutdown(context.Background())
	}()

	log.Printf("AINET DB Explore listening on http://localhost%s (idle-timeout: %s, db: %s)", *addr, *idle, dbPath)
	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
