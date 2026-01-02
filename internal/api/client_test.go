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
		_ = json.NewEncoder(w).Encode(response)
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
		_ = json.NewEncoder(w).Encode(response)
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
		_ = json.NewEncoder(w).Encode(response)
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

func TestSearchIssues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"issues": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":         "issue-1",
							"identifier": "ENG-123",
							"title":      "Fix login bug",
							"priority":   2,
							"url":        "https://linear.app/test/issue/ENG-123",
							"state": map[string]interface{}{
								"id":   "state-1",
								"name": "In Progress",
							},
							"assignee": map[string]interface{}{
								"id":   "user-1",
								"name": "John Doe",
							},
							"team": map[string]interface{}{
								"id":  "team-1",
								"key": "ENG",
							},
						},
						{
							"id":         "issue-2",
							"identifier": "ENG-124",
							"title":      "Login page redesign",
							"priority":   3,
							"url":        "https://linear.app/test/issue/ENG-124",
							"state": map[string]interface{}{
								"id":   "state-2",
								"name": "Backlog",
							},
							"assignee": nil,
							"team": map[string]interface{}{
								"id":  "team-1",
								"key": "ENG",
							},
						},
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &LinearClient{
		gql: createTestClient(server.URL, "test-api-key"),
	}

	ctx := context.Background()
	issues, err := client.SearchIssues(ctx, "login", IssueListOptions{First: 10})

	require.NoError(t, err)
	assert.Len(t, issues, 2)

	// Check first issue
	assert.Equal(t, "ENG-123", issues[0].Identifier)
	assert.Equal(t, "Fix login bug", issues[0].Title)
	assert.Equal(t, 2, issues[0].Priority)
	assert.Equal(t, "In Progress", issues[0].State.Name)
	assert.Equal(t, "John Doe", issues[0].Assignee.Name)
	assert.Equal(t, "ENG", issues[0].Team.Key)

	// Check second issue (no assignee)
	assert.Equal(t, "ENG-124", issues[1].Identifier)
	assert.Nil(t, issues[1].Assignee)
}

func TestSearchIssues_EmptyResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"issues": map[string]interface{}{
					"nodes": []map[string]interface{}{},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &LinearClient{
		gql: createTestClient(server.URL, "test-api-key"),
	}

	ctx := context.Background()
	issues, err := client.SearchIssues(ctx, "nonexistent", IssueListOptions{})

	require.NoError(t, err)
	assert.Empty(t, issues)
}

func TestEscapeGraphQLString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "simple query",
			expected: "simple query",
		},
		{
			name:     "double quotes",
			input:    `search "term"`,
			expected: `search \"term\"`,
		},
		{
			name:     "backslash",
			input:    `path\to\file`,
			expected: `path\\to\\file`,
		},
		{
			name:     "newline",
			input:    "line1\nline2",
			expected: `line1\nline2`,
		},
		{
			name:     "tab",
			input:    "col1\tcol2",
			expected: `col1\tcol2`,
		},
		{
			name:     "carriage return",
			input:    "line1\rline2",
			expected: `line1\rline2`,
		},
		{
			name:     "mixed special characters",
			input:    "test \"value\"\nwith\\special",
			expected: `test \"value\"\nwith\\special`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeGraphQLString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
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
