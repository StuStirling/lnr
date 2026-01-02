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

func TestGetUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"users": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "user-1",
							"name":        "Alice",
							"email":       "alice@example.com",
							"displayName": "Alice A",
							"active":      true,
							"admin":       true,
							"avatarUrl":   "https://example.com/alice.png",
						},
						{
							"id":          "user-2",
							"name":        "Bob",
							"email":       "bob@example.com",
							"displayName": "Bob B",
							"active":      true,
							"admin":       false,
							"avatarUrl":   "https://example.com/bob.png",
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
	users, err := client.GetUsers(ctx)

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	assert.True(t, users[0].Admin)
	assert.Equal(t, "Bob", users[1].Name)
	assert.False(t, users[1].Admin)
}

func TestGetIssues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"issues": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "issue-1",
							"identifier":  "ENG-100",
							"title":       "Test Issue",
							"description": "Test description",
							"priority":    1,
							"url":         "https://linear.app/test/issue/ENG-100",
							"createdAt":   "2024-01-01T00:00:00Z",
							"updatedAt":   "2024-01-02T00:00:00Z",
							"state": map[string]interface{}{
								"id":    "state-1",
								"name":  "Todo",
								"color": "#000",
								"type":  "unstarted",
							},
							"team": map[string]interface{}{
								"id":   "team-1",
								"name": "Engineering",
								"key":  "ENG",
							},
							"labels": map[string]interface{}{
								"nodes": []map[string]interface{}{},
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
	issues, err := client.GetIssues(ctx, IssueListOptions{First: 10})

	require.NoError(t, err)
	assert.Len(t, issues, 1)
	assert.Equal(t, "ENG-100", issues[0].Identifier)
	assert.Equal(t, "Test Issue", issues[0].Title)
	assert.Equal(t, 1, issues[0].Priority)
	assert.Equal(t, "Todo", issues[0].State.Name)
	assert.Equal(t, "ENG", issues[0].Team.Key)
}

func TestGetProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"projects": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "proj-1",
							"name":        "Project Alpha",
							"description": "First project",
							"state":       "started",
							"progress":    0.5,
							"url":         "https://linear.app/test/project/proj-1",
							"createdAt":   "2024-01-01T00:00:00Z",
							"updatedAt":   "2024-01-02T00:00:00Z",
							"teams": map[string]interface{}{
								"nodes": []map[string]interface{}{
									{"id": "team-1", "name": "Engineering", "key": "ENG"},
								},
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
	projects, err := client.GetProjects(ctx, ProjectListOptions{First: 10})

	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "Project Alpha", projects[0].Name)
	assert.Equal(t, "started", projects[0].State)
	assert.Equal(t, 0.5, projects[0].Progress)
	assert.Len(t, projects[0].Teams, 1)
	assert.Equal(t, "ENG", projects[0].Teams[0].Key)
}

func TestGetCycles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"cycles": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "cycle-1",
							"name":        "Sprint 1",
							"number":      1,
							"startsAt":    "2024-01-01T00:00:00Z",
							"endsAt":      "2024-01-14T00:00:00Z",
							"progress":    0.75,
							"description": "First sprint",
							"team": map[string]interface{}{
								"id":   "team-1",
								"name": "Engineering",
								"key":  "ENG",
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
	cycles, err := client.GetCycles(ctx, nil)

	require.NoError(t, err)
	assert.Len(t, cycles, 1)
	assert.Equal(t, "Sprint 1", cycles[0].Name)
	assert.Equal(t, 1, cycles[0].Number)
	assert.Equal(t, 0.75, cycles[0].Progress)
	assert.Equal(t, "ENG", cycles[0].Team.Key)
}

func TestGetActiveCycle_NoCycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"team": map[string]interface{}{
					"id":          "team-1",
					"name":        "Engineering",
					"key":         "ENG",
					"activeCycle": nil,
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
	cycle, err := client.GetActiveCycle(ctx, "team-1")

	assert.Nil(t, cycle)
	assert.ErrorIs(t, err, ErrNoActiveCycle)
}

func TestRetryTransport_RateLimited(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"viewer": map[string]interface{}{
					"id":    "user-1",
					"name":  "Test",
					"email": "test@example.com",
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with retry transport
	httpClient := &http.Client{
		Transport: &retryTransport{
			maxRetries: 3,
			transport: &authTransport{
				apiKey:    "test-key",
				transport: http.DefaultTransport,
			},
		},
	}
	client := &LinearClient{
		gql: graphql.NewClient(server.URL, httpClient),
	}

	ctx := context.Background()
	user, err := client.GetViewer(ctx)

	require.NoError(t, err)
	assert.Equal(t, "Test", user.Name)
	assert.Equal(t, 3, attempts) // Should have retried twice
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
