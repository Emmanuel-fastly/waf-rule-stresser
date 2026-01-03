// main.go is the entry point for the WAF Rate Limit Tester backend server.
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Load attack payloads from text files
	err := LoadPayloads()
	if err != nil {
		log.Printf("Warning: %v - Attack mode will have limited functionality\n", err)
	}

	// Setup API routes
	setupRoutes()

	// Setup frontend serving
	setupFrontend()

	// Professional startup message
	fmt.Println("WAF Tester started")
	fmt.Printf("Open in browser: http://localhost:%s\n", port)

	// Start the HTTP server
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

// setupRoutes configures all the API endpoints
func setupRoutes() {
	http.HandleFunc("/api/test/start", handleStartTest)
	http.HandleFunc("/api/test/stream", handleStartTestStream)
	http.HandleFunc("/api/test/stop", handleStopTest)
	http.HandleFunc("/api/test/export", handleExport)
	http.HandleFunc("/api/exports/list", handleListExports)
}

// setupFrontend configures serving the embedded frontend
func setupFrontend() {
	// Get the embedded dist directory
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to load embedded frontend:", err)
	}

	// Create file server for static assets
	fileServer := http.FileServer(http.FS(distFS))

	// Handle all non-API routes with the frontend
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If it's an API request, skip (already handled by setupRoutes)
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			http.NotFound(w, r)
			return
		}

		// Try to serve the requested file
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists in embedded FS
		if _, err := distFS.Open(path[1:]); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// For SPA routing, serve index.html for all other routes
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
