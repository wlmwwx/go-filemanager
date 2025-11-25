package main

import (
	"encoding/json"
	"os"
	"sync"
)

// UserCredentials stores user authentication information
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	credentialsFile = ".filemanager_credentials.json"
	credentialsMux  sync.RWMutex
)

// LoadCredentials loads user credentials from file
func LoadCredentials() (*UserCredentials, error) {
	credentialsMux.RLock()
	defer credentialsMux.RUnlock()

	// Check if file exists
	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		// Return default credentials if file doesn't exist
		return &UserCredentials{
			Username: "admin",
			Password: "admin123",
		}, nil
	}

	data, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}

	var creds UserCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

// SaveCredentials saves user credentials to file
func SaveCredentials(creds *UserCredentials) error {
	credentialsMux.Lock()
	defer credentialsMux.Unlock()

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(credentialsFile, data, 0600)
}
