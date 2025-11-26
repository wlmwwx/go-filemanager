package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

var store *sessions.CookieStore

func main() {
	// Load configuration
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Change to root directory if configured
	if config.Server.RootDir != "" {
		absRootDir, err := filepath.Abs(config.Server.RootDir)
		if err != nil {
			log.Fatalf("Failed to resolve root directory: %v", err)
		}
		if err := os.Chdir(absRootDir); err != nil {
			log.Fatalf("Failed to change to root directory %s: %v", absRootDir, err)
		}
		log.Printf("Root directory set to: %s", absRootDir)
	} else {
		cwd, _ := os.Getwd()
		log.Printf("Root directory: %s (current working directory)", cwd)
	}

	// Initialize session store with configured secret
	store = sessions.NewCookieStore([]byte(config.Secret))

	// Setup logging
	if err := setupLogging(config.Log); err != nil {
		log.Fatalf("Failed to setup logging: %v", err)
	}

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
	protected.HandleFunc("/change-password", ChangePasswordHandler).Methods("POST", "OPTIONS")

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

	addr := fmt.Sprintf(":%d", config.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Println("Default credentials: admin / admin123")
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

// setupLogging configures logging based on config
func setupLogging(logConfig LogConfig) error {
	if logConfig.File != "" {
		// Open log file
		f, err := os.OpenFile(logConfig.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("error opening log file: %v", err)
		}

		// Write to both file and stdout
		multiWriter := io.MultiWriter(os.Stdout, f)
		log.SetOutput(multiWriter)
		log.Printf("Logging to file: %s", logConfig.File)
	}

	// Set log flags based on level
	if logConfig.Level == "debug" {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	return nil
}
