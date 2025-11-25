package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Simple authentication (in production, use proper password hashing)
	if req.Username == "admin" && req.Password == "admin123" {
		session, _ := store.Get(r, "session")
		session.Values["authenticated"] = true
		session.Values["username"] = req.Username
		if err := session.Save(r, w); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save session")
			return
		}
		respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Login successful"})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
	}
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}
	respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Logout successful"})
}

// ListFilesHandler lists files in a directory
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		path = "."
	}

	// Sanitize path to prevent directory traversal
	cleanPath := filepath.Clean(path)
	if strings.HasPrefix(cleanPath, "..") {
		respondWithError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	// Get absolute path relative to working directory
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to resolve path")
		return
	}

	// Ensure path is within working directory
	workDir, _ := os.Getwd()
	if !strings.HasPrefix(absPath, workDir) {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read directory")
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   entry.IsDir(),
		})
	}

	response := DirectoryResponse{
		Path:  cleanPath,
		Files: files,
	}
	respondWithJSON(w, http.StatusOK, response)
}

// UploadHandler handles file uploads
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (32MB max)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	path := r.FormValue("path")
	if path == "" {
		path = "."
	}

	// Sanitize path
	cleanPath := filepath.Clean(path)
	if strings.HasPrefix(cleanPath, "..") {
		respondWithError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	// Create destination file
	destPath := filepath.Join(cleanPath, handler.Filename)
	absPath, err := filepath.Abs(destPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to resolve path")
		return
	}

	// Ensure path is within working directory
	workDir, _ := os.Getwd()
	if !strings.HasPrefix(absPath, workDir) {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	dst, err := os.Create(absPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save file")
		return
	}

	respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "File uploaded successfully"})
}

// DownloadHandler handles file downloads
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		respondWithError(w, http.StatusBadRequest, "Path required")
		return
	}

	// Sanitize path
	cleanPath := filepath.Clean(path)
	if strings.HasPrefix(cleanPath, "..") {
		respondWithError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to resolve path")
		return
	}

	// Ensure path is within working directory
	workDir, _ := os.Getwd()
	if !strings.HasPrefix(absPath, workDir) {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	// Check if file exists and is not a directory
	info, err := os.Stat(absPath)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "File not found")
		return
	}
	if info.IsDir() {
		respondWithError(w, http.StatusBadRequest, "Cannot download directory")
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(absPath)))
	http.ServeFile(w, r, absPath)
}

// MkdirHandler creates a new directory
func MkdirHandler(w http.ResponseWriter, r *http.Request) {
	var req MkdirRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Directory name required")
		return
	}

	basePath := req.Path
	if basePath == "" {
		basePath = "."
	}

	// Sanitize path
	cleanPath := filepath.Clean(filepath.Join(basePath, req.Name))
	if strings.HasPrefix(cleanPath, "..") {
		respondWithError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to resolve path")
		return
	}

	// Ensure path is within working directory
	workDir, _ := os.Getwd()
	if !strings.HasPrefix(absPath, workDir) {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create directory")
		return
	}

	respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Directory created successfully"})
}

// DeleteHandler deletes a file or directory
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	if req.Path == "" {
		respondWithError(w, http.StatusBadRequest, "Path required")
		return
	}

	// Sanitize path
	cleanPath := filepath.Clean(req.Path)
	if strings.HasPrefix(cleanPath, "..") || cleanPath == "." {
		respondWithError(w, http.StatusBadRequest, "Invalid path")
		return
	}

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to resolve path")
		return
	}

	// Ensure path is within working directory
	workDir, _ := os.Getwd()
	if !strings.HasPrefix(absPath, workDir) || absPath == workDir {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := os.RemoveAll(absPath); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete")
		return
	}

	respondWithJSON(w, http.StatusOK, SuccessResponse{Message: "Deleted successfully"})
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, ErrorResponse{Error: message})
}
