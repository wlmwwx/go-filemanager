package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

var store = sessions.NewCookieStore([]byte("your-secret-key-change-this-in-production"))

func main() {
	r := mux.NewRouter()

	// Apply logging middleware to all routes
	r.Use(LoggingMiddleware)
	r.Use(CORSMiddleware)

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Public routes
	api.HandleFunc("/login", LoginHandler).Methods("POST", "OPTIONS")

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(AuthMiddleware)
	protected.HandleFunc("/logout", LogoutHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/files", ListFilesHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/upload", UploadHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/download", DownloadHandler).Methods("GET", "OPTIONS")
	protected.HandleFunc("/mkdir", MkdirHandler).Methods("POST", "OPTIONS")
	protected.HandleFunc("/delete", DeleteHandler).Methods("DELETE", "OPTIONS")

	// Serve frontend static files
	frontendSubFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal("Failed to load frontend assets:", err)
	}

	// Serve static files and SPA fallback
	r.PathPrefix("/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists
		if _, err := frontendSubFS.Open(path[1:]); err == nil {
			http.FileServer(http.FS(frontendSubFS)).ServeHTTP(w, r)
		} else {
			// Fallback to index.html for SPA routing
			r.URL.Path = "/"
			http.FileServer(http.FS(frontendSubFS)).ServeHTTP(w, r)
		}
	}))

	log.Println("Server starting on :8080")
	log.Println("Default credentials: admin / admin123")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
