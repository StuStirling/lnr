package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WithAPIKey(t *testing.T) {
	// Set up
	os.Setenv(EnvAPIKey, "test-api-key")
	defer os.Unsetenv(EnvAPIKey)

	// Execute
	cfg, err := Load()

	// Verify
	require.NoError(t, err)
	assert.Equal(t, "test-api-key", cfg.APIKey)
}

func TestLoad_WithoutAPIKey(t *testing.T) {
	// Set up
	os.Unsetenv(EnvAPIKey)

	// Execute
	cfg, err := Load()

	// Verify
	assert.Nil(t, cfg)
	assert.Equal(t, ErrNoAPIKey, err)
}

func TestMustLoad_Panics(t *testing.T) {
	// Set up
	os.Unsetenv(EnvAPIKey)

	// Verify it panics
	assert.Panics(t, func() {
		MustLoad()
	})
}
