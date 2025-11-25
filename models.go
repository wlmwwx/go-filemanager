package main

import "time"

// LoginRequest represents the login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// FileInfo represents file or directory metadata
type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
}

// DirectoryResponse represents the response for directory listing
type DirectoryResponse struct {
	Path  string     `json:"path"`
	Files []FileInfo `json:"files"`
}

// ErrorResponse represents an error message
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success message
type SuccessResponse struct {
	Message string `json:"message"`
}

// MkdirRequest represents a request to create a directory
type MkdirRequest struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// DeleteRequest represents a request to delete a file or directory
type DeleteRequest struct {
	Path string `json:"path"`
}
