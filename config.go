package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server ServerConfig `yaml:"server"`
	Secret string       `yaml:"secret"`
	Log    LogConfig    `yaml:"log"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port int `yaml:"port"`
}

// LogConfig represents logging configuration
type LogConfig struct {
	File  string `yaml:"file"`
	Level string `yaml:"level"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
		Secret: "your-secret-key-change-this-in-production",
		Log: LogConfig{
			File:  "",
			Level: "info",
		},
	}
}

// LoadConfig loads configuration from file or returns default
func LoadConfig(filename string) (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("Config file not found, using defaults")
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveDefaultConfig saves the default configuration to file
func SaveDefaultConfig(filename string) error {
	config := DefaultConfig()
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
