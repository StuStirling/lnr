package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_WithAPIKey(t *testing.T) {
	// Set up - t.Setenv automatically cleans up after test
	t.Setenv(EnvAPIKey, "test-api-key")

	// Execute
	cfg, err := Load()

	// Verify
	require.NoError(t, err)
	assert.Equal(t, "test-api-key", cfg.APIKey)
}

func TestLoad_WithoutAPIKey(t *testing.T) {
	// Set up - ensure env var is not set
	t.Setenv(EnvAPIKey, "")

	// Execute
	cfg, err := Load()

	// Verify
	assert.Nil(t, cfg)
	assert.Equal(t, ErrNoAPIKey, err)
}

func TestMustLoad_Panics(t *testing.T) {
	// Set up - ensure env var is not set
	t.Setenv(EnvAPIKey, "")

	// Verify it panics
	assert.Panics(t, func() {
		MustLoad()
	})
}
