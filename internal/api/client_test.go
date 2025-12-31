package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hasura/go-graphql-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key")
	assert.NotNil(t, client)
	assert.NotNil(t, client.gql)
}

func TestGetViewer(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth header
		assert.Equal(t, "test-api-key", r.Header.Get("Authorization"))

		// Return mock response
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"viewer": map[string]interface{}{
					"id":          "user-123",
					"name":        "Test User",
					"email":       "test@example.com",
					"displayName": "Test",
					"active":      true,
					"admin":       false,
					"avatarUrl":   "https://example.com/avatar.png",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client pointing to mock server
	client := &LinearClient{
		gql: createTestClient(server.URL, "test-api-key"),
	}

	// Execute
	ctx := context.Background()
	user, err := client.GetViewer(ctx)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.True(t, user.Active)
	assert.False(t, user.Admin)
}

func TestGetOrganisation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"organization": map[string]interface{}{
					"id":        "org-123",
					"name":      "Test Org",
					"urlKey":    "test-org",
					"logoUrl":   "https://example.com/logo.png",
					"userCount": 10,
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &LinearClient{
		gql: createTestClient(server.URL, "test-api-key"),
	}

	ctx := context.Background()
	org, err := client.GetOrganisation(ctx)

	require.NoError(t, err)
	assert.Equal(t, "org-123", org.ID)
	assert.Equal(t, "Test Org", org.Name)
	assert.Equal(t, 10, org.UserCount)
}

func TestGetTeams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"teams": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "team-1",
							"name":        "Engineering",
							"key":         "ENG",
							"description": "Engineering team",
							"private":     false,
						},
						{
							"id":          "team-2",
							"name":        "Product",
							"key":         "PRD",
							"description": "Product team",
							"private":     false,
						},
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &LinearClient{
		gql: createTestClient(server.URL, "test-api-key"),
	}

	ctx := context.Background()
	teams, err := client.GetTeams(ctx)

	require.NoError(t, err)
	assert.Len(t, teams, 2)
	assert.Equal(t, "ENG", teams[0].Key)
	assert.Equal(t, "PRD", teams[1].Key)
}

// Helper to create a test GraphQL client pointing to a test server
func createTestClient(url, apiKey string) *graphqlClient {
	httpClient := &http.Client{
		Transport: &authTransport{
			apiKey:    apiKey,
			transport: http.DefaultTransport,
		},
	}
	return graphql.NewClient(url, httpClient)
}

// graphqlClient is an alias for the graphql client type
type graphqlClient = graphql.Client
