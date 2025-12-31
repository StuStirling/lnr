package config

import (
	"errors"
	"os"
)

const (
	// EnvAPIKey is the environment variable name for the Linear API key
	EnvAPIKey = "LINEAR_API_KEY"
)

var (
	// ErrNoAPIKey is returned when the API key is not set
	ErrNoAPIKey = errors.New("LINEAR_API_KEY environment variable is not set")
)

// Config holds the application configuration
type Config struct {
	APIKey string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	apiKey := os.Getenv(EnvAPIKey)
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}

	return &Config{
		APIKey: apiKey,
	}, nil
}

// MustLoad loads configuration and panics if it fails
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}
