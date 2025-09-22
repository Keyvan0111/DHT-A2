package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func shortHost() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "unknown"
	}

	// Trim domain (e.g., c1-0.ifi.uit.no -> c1-0) to match assignment samples
	if i := strings.IndexByte(h, '.'); i > 0 {
		return h[:i]
	}
	return h
}

func main() {
	// Pick a free port
	ln, err := net.Listen("tcp4", "0.0.0.0:0")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	host := shortHost()
	host1, _ := os.Hostname()

	// If PORT_FILE is set, write the chosen port there so run.sh can read it
	if path := os.Getenv("PORT_FILE"); path != "" {
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("failed to create PORT_FILE %q: %v", path, err)
		}
		if _, err := fmt.Fprintf(f, "%d", port); err != nil {
			_ = f.Close()
			log.Fatalf("failed to write PORT_FILE %q: %v", path, err)
		}
		_ = f.Sync()
		_ = f.Close()
	}

	// Serve the required endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s:%d", host, port)
		log.Println("Hello server guys!")
	})

	fmt.Println("Hello guys im here!")

	log.Printf("listening on %s:%d\n", host1, port)

	// Serve using the already-open socket (keeps same :0-chosen port)
	server := &http.Server{Handler: mux}
	if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http serve error: %v", err)
	}
}
