package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
)

// ErrNoActiveCycle is returned when a team has no active cycle
var ErrNoActiveCycle = errors.New("no active cycle for this team")

const (
	// LinearAPIEndpoint is the Linear GraphQL API endpoint
	LinearAPIEndpoint = "https://api.linear.app/graphql"
)

// Client is the interface for the Linear API client
type Client interface {
	// Viewer
	GetViewer(ctx context.Context) (*User, error)
	GetOrganisation(ctx context.Context) (*Organisation, error)

	// Users
	GetUsers(ctx context.Context) ([]User, error)

	// Teams
	GetTeams(ctx context.Context) ([]Team, error)
	GetTeam(ctx context.Context, id string) (*Team, error)

	// Labels
	GetLabels(ctx context.Context, teamID *string) ([]Label, error)

	// Workflow States
	GetWorkflowStates(ctx context.Context, teamID *string) ([]WorkflowState, error)

	// Issues
	GetIssues(ctx context.Context, opts IssueListOptions) ([]Issue, error)
	GetIssue(ctx context.Context, id string) (*Issue, error)
	SearchIssues(ctx context.Context, query string, opts IssueListOptions) ([]Issue, error)

	// Projects
	GetProjects(ctx context.Context, opts ProjectListOptions) ([]Project, error)
	GetProject(ctx context.Context, id string) (*Project, error)

	// Initiatives
	GetInitiatives(ctx context.Context) ([]Initiative, error)
	GetInitiative(ctx context.Context, id string) (*Initiative, error)

	// Cycles
	GetCycles(ctx context.Context, teamID *string) ([]Cycle, error)
	GetActiveCycle(ctx context.Context, teamID string) (*Cycle, error)
	GetCycle(ctx context.Context, id string) (*Cycle, error)
}

// IssueListOptions contains options for listing issues
type IssueListOptions struct {
	TeamID     *string
	AssigneeID *string
	StateID    *string
	ProjectID  *string
	First      int
}

// ProjectListOptions contains options for listing projects
type ProjectListOptions struct {
	TeamID *string
	State  *string
	First  int
}

// LinearClient implements the Client interface
type LinearClient struct {
	gql *graphql.Client
}

// authTransport adds authorization header to requests
type authTransport struct {
	apiKey    string
	transport http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return t.transport.RoundTrip(req)
}

// retryTransport wraps a transport and retries on rate limit errors
type retryTransport struct {
	transport  http.RoundTripper
	maxRetries int
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Read and store the body so we can retry
	var bodyBytes []byte
	if req.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}
		_ = req.Body.Close()
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		// Restore the body for each attempt
		if bodyBytes != nil {
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		resp, err = t.transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		// If not rate limited, return immediately
		if resp.StatusCode != http.StatusTooManyRequests {
			return resp, nil
		}

		// Don't retry on last attempt
		if attempt == t.maxRetries {
			break
		}

		// Close the response body before retrying
		_ = resp.Body.Close()

		// Calculate backoff: 1s, 2s, 4s
		backoff := time.Duration(1<<attempt) * time.Second

		// Check for Retry-After header
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			if seconds, parseErr := time.ParseDuration(retryAfter + "s"); parseErr == nil {
				backoff = seconds
			}
		}

		// Wait before retrying
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		case <-time.After(backoff):
		}
	}

	return resp, nil
}

// NewClient creates a new Linear API client
func NewClient(apiKey string) *LinearClient {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &retryTransport{
			maxRetries: 3,
			transport: &authTransport{
				apiKey:    apiKey,
				transport: http.DefaultTransport,
			},
		},
	}

	return &LinearClient{
		gql: graphql.NewClient(LinearAPIEndpoint, httpClient),
	}
}

// GetViewer returns the currently authenticated user
func (c *LinearClient) GetViewer(ctx context.Context) (*User, error) {
	var query struct {
		Viewer struct {
			ID          string `graphql:"id"`
			Name        string `graphql:"name"`
			Email       string `graphql:"email"`
			DisplayName string `graphql:"displayName"`
			Active      bool   `graphql:"active"`
			Admin       bool   `graphql:"admin"`
			AvatarURL   string `graphql:"avatarUrl"`
		} `graphql:"viewer"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get viewer: %w", err)
	}

	return &User{
		ID:          query.Viewer.ID,
		Name:        query.Viewer.Name,
		Email:       query.Viewer.Email,
		DisplayName: query.Viewer.DisplayName,
		Active:      query.Viewer.Active,
		Admin:       query.Viewer.Admin,
		AvatarURL:   query.Viewer.AvatarURL,
	}, nil
}

// GetOrganisation returns the current organisation
func (c *LinearClient) GetOrganisation(ctx context.Context) (*Organisation, error) {
	var query struct {
		Organisation struct {
			ID        string `graphql:"id"`
			Name      string `graphql:"name"`
			URLKey    string `graphql:"urlKey"`
			LogoURL   string `graphql:"logoUrl"`
			UserCount int    `graphql:"userCount"`
		} `graphql:"organization"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get organisation: %w", err)
	}

	return &Organisation{
		ID:        query.Organisation.ID,
		Name:      query.Organisation.Name,
		URLKey:    query.Organisation.URLKey,
		LogoURL:   query.Organisation.LogoURL,
		UserCount: query.Organisation.UserCount,
	}, nil
}

// GetUsers returns all users in the organisation
func (c *LinearClient) GetUsers(ctx context.Context) ([]User, error) {
	var query struct {
		Users struct {
			Nodes []struct {
				ID          string `graphql:"id"`
				Name        string `graphql:"name"`
				Email       string `graphql:"email"`
				DisplayName string `graphql:"displayName"`
				Active      bool   `graphql:"active"`
				Admin       bool   `graphql:"admin"`
				AvatarURL   string `graphql:"avatarUrl"`
			} `graphql:"nodes"`
		} `graphql:"users(first: 100)"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	users := make([]User, len(query.Users.Nodes))
	for i, u := range query.Users.Nodes {
		users[i] = User{
			ID:          u.ID,
			Name:        u.Name,
			Email:       u.Email,
			DisplayName: u.DisplayName,
			Active:      u.Active,
			Admin:       u.Admin,
			AvatarURL:   u.AvatarURL,
		}
	}
	return users, nil
}

// GetTeams returns all teams in the organisation
func (c *LinearClient) GetTeams(ctx context.Context) ([]Team, error) {
	var query struct {
		Teams struct {
			Nodes []struct {
				ID          string `graphql:"id"`
				Name        string `graphql:"name"`
				Key         string `graphql:"key"`
				Description string `graphql:"description"`
				Private     bool   `graphql:"private"`
			} `graphql:"nodes"`
		} `graphql:"teams(first: 100)"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get teams: %w", err)
	}

	teams := make([]Team, len(query.Teams.Nodes))
	for i, t := range query.Teams.Nodes {
		teams[i] = Team{
			ID:          t.ID,
			Name:        t.Name,
			Key:         t.Key,
			Description: t.Description,
			Private:     t.Private,
		}
	}
	return teams, nil
}

// GetTeam returns a single team by ID
func (c *LinearClient) GetTeam(ctx context.Context, id string) (*Team, error) {
	var query struct {
		Team struct {
			ID          string `graphql:"id"`
			Name        string `graphql:"name"`
			Key         string `graphql:"key"`
			Description string `graphql:"description"`
			Private     bool   `graphql:"private"`
		} `graphql:"team(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(id),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get team: %w", err)
	}

	return &Team{
		ID:          query.Team.ID,
		Name:        query.Team.Name,
		Key:         query.Team.Key,
		Description: query.Team.Description,
		Private:     query.Team.Private,
	}, nil
}

// GetLabels returns labels, optionally filtered by team
func (c *LinearClient) GetLabels(ctx context.Context, teamID *string) ([]Label, error) {
	var query struct {
		IssueLabels struct {
			Nodes []struct {
				ID          string `graphql:"id"`
				Name        string `graphql:"name"`
				Description string `graphql:"description"`
				Color       string `graphql:"color"`
				Team        *struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
					Key  string `graphql:"key"`
				} `graphql:"team"`
			} `graphql:"nodes"`
		} `graphql:"issueLabels(first: 100)"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get labels: %w", err)
	}

	labels := make([]Label, 0, len(query.IssueLabels.Nodes))
	for _, l := range query.IssueLabels.Nodes {
		// Filter by team if specified
		if teamID != nil && (l.Team == nil || l.Team.ID != *teamID) {
			continue
		}

		label := Label{
			ID:          l.ID,
			Name:        l.Name,
			Description: l.Description,
			Color:       l.Color,
		}
		if l.Team != nil {
			label.Team = &Team{
				ID:   l.Team.ID,
				Name: l.Team.Name,
				Key:  l.Team.Key,
			}
		}
		labels = append(labels, label)
	}
	return labels, nil
}

// GetWorkflowStates returns workflow states, optionally filtered by team
func (c *LinearClient) GetWorkflowStates(ctx context.Context, teamID *string) ([]WorkflowState, error) {
	var query struct {
		WorkflowStates struct {
			Nodes []struct {
				ID       string `graphql:"id"`
				Name     string `graphql:"name"`
				Color    string `graphql:"color"`
				Type     string `graphql:"type"`
				Position int    `graphql:"position"`
				Team     struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
					Key  string `graphql:"key"`
				} `graphql:"team"`
			} `graphql:"nodes"`
		} `graphql:"workflowStates(first: 100)"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get workflow states: %w", err)
	}

	states := make([]WorkflowState, 0, len(query.WorkflowStates.Nodes))
	for _, s := range query.WorkflowStates.Nodes {
		// Filter by team if specified
		if teamID != nil && s.Team.ID != *teamID {
			continue
		}

		states = append(states, WorkflowState{
			ID:       s.ID,
			Name:     s.Name,
			Color:    s.Color,
			Type:     s.Type,
			Position: s.Position,
			Team: &Team{
				ID:   s.Team.ID,
				Name: s.Team.Name,
				Key:  s.Team.Key,
			},
		})
	}
	return states, nil
}

// GetIssues returns issues with optional filters
func (c *LinearClient) GetIssues(ctx context.Context, opts IssueListOptions) ([]Issue, error) {
	first := opts.First
	if first == 0 {
		first = 50
	}

	var query struct {
		Issues struct {
			Nodes []struct {
				ID          string  `graphql:"id"`
				Identifier  string  `graphql:"identifier"`
				Title       string  `graphql:"title"`
				Description string  `graphql:"description"`
				Priority    int     `graphql:"priority"`
				Estimate    *float64 `graphql:"estimate"`
				URL         string  `graphql:"url"`
				CreatedAt   string  `graphql:"createdAt"`
				UpdatedAt   string  `graphql:"updatedAt"`
				DueDate     *string `graphql:"dueDate"`
				State       *struct {
					ID    string `graphql:"id"`
					Name  string `graphql:"name"`
					Color string `graphql:"color"`
					Type  string `graphql:"type"`
				} `graphql:"state"`
				Assignee *struct {
					ID    string `graphql:"id"`
					Name  string `graphql:"name"`
					Email string `graphql:"email"`
				} `graphql:"assignee"`
				Team struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
					Key  string `graphql:"key"`
				} `graphql:"team"`
				Project *struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
				} `graphql:"project"`
				Labels struct {
					Nodes []struct {
						ID    string `graphql:"id"`
						Name  string `graphql:"name"`
						Color string `graphql:"color"`
					} `graphql:"nodes"`
				} `graphql:"labels"`
			} `graphql:"nodes"`
		} `graphql:"issues(first: $first)"`
	}

	vars := map[string]interface{}{
		"first": graphql.Int(first),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get issues: %w", err)
	}

	issues := make([]Issue, 0, len(query.Issues.Nodes))
	for _, i := range query.Issues.Nodes {
		// Apply filters
		if opts.TeamID != nil && i.Team.ID != *opts.TeamID {
			continue
		}
		if opts.AssigneeID != nil && (i.Assignee == nil || i.Assignee.ID != *opts.AssigneeID) {
			continue
		}
		if opts.StateID != nil && (i.State == nil || i.State.ID != *opts.StateID) {
			continue
		}
		if opts.ProjectID != nil && (i.Project == nil || i.Project.ID != *opts.ProjectID) {
			continue
		}

		issue := Issue{
			ID:          i.ID,
			Identifier:  i.Identifier,
			Title:       i.Title,
			Description: i.Description,
			Priority:    i.Priority,
			Estimate:    i.Estimate,
			URL:         i.URL,
			DueDate:     i.DueDate,
			Team: &Team{
				ID:   i.Team.ID,
				Name: i.Team.Name,
				Key:  i.Team.Key,
			},
		}

		if i.State != nil {
			issue.State = &WorkflowState{
				ID:    i.State.ID,
				Name:  i.State.Name,
				Color: i.State.Color,
				Type:  i.State.Type,
			}
		}

		if i.Assignee != nil {
			issue.Assignee = &User{
				ID:    i.Assignee.ID,
				Name:  i.Assignee.Name,
				Email: i.Assignee.Email,
			}
		}

		if i.Project != nil {
			issue.Project = &Project{
				ID:   i.Project.ID,
				Name: i.Project.Name,
			}
		}

		for _, l := range i.Labels.Nodes {
			issue.Labels = append(issue.Labels, Label{
				ID:    l.ID,
				Name:  l.Name,
				Color: l.Color,
			})
		}

		issues = append(issues, issue)
	}
	return issues, nil
}

// GetIssue returns a single issue by ID or identifier
func (c *LinearClient) GetIssue(ctx context.Context, id string) (*Issue, error) {
	var query struct {
		Issue struct {
			ID          string   `graphql:"id"`
			Identifier  string   `graphql:"identifier"`
			Title       string   `graphql:"title"`
			Description string   `graphql:"description"`
			Priority    int      `graphql:"priority"`
			Estimate    *float64 `graphql:"estimate"`
			URL         string   `graphql:"url"`
			CreatedAt   string   `graphql:"createdAt"`
			UpdatedAt   string   `graphql:"updatedAt"`
			DueDate     *string  `graphql:"dueDate"`
			State       *struct {
				ID    string `graphql:"id"`
				Name  string `graphql:"name"`
				Color string `graphql:"color"`
				Type  string `graphql:"type"`
			} `graphql:"state"`
			Assignee *struct {
				ID    string `graphql:"id"`
				Name  string `graphql:"name"`
				Email string `graphql:"email"`
			} `graphql:"assignee"`
			Creator *struct {
				ID    string `graphql:"id"`
				Name  string `graphql:"name"`
				Email string `graphql:"email"`
			} `graphql:"creator"`
			Team struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				Key  string `graphql:"key"`
			} `graphql:"team"`
			Project *struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"project"`
			Cycle *struct {
				ID     string `graphql:"id"`
				Name   string `graphql:"name"`
				Number int    `graphql:"number"`
			} `graphql:"cycle"`
			Labels struct {
				Nodes []struct {
					ID    string `graphql:"id"`
					Name  string `graphql:"name"`
					Color string `graphql:"color"`
				} `graphql:"nodes"`
			} `graphql:"labels"`
		} `graphql:"issue(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(id),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get issue: %w", err)
	}

	i := query.Issue
	issue := &Issue{
		ID:          i.ID,
		Identifier:  i.Identifier,
		Title:       i.Title,
		Description: i.Description,
		Priority:    i.Priority,
		Estimate:    i.Estimate,
		URL:         i.URL,
		DueDate:     i.DueDate,
		Team: &Team{
			ID:   i.Team.ID,
			Name: i.Team.Name,
			Key:  i.Team.Key,
		},
	}

	if i.State != nil {
		issue.State = &WorkflowState{
			ID:    i.State.ID,
			Name:  i.State.Name,
			Color: i.State.Color,
			Type:  i.State.Type,
		}
	}

	if i.Assignee != nil {
		issue.Assignee = &User{
			ID:    i.Assignee.ID,
			Name:  i.Assignee.Name,
			Email: i.Assignee.Email,
		}
	}

	if i.Creator != nil {
		issue.Creator = &User{
			ID:    i.Creator.ID,
			Name:  i.Creator.Name,
			Email: i.Creator.Email,
		}
	}

	if i.Project != nil {
		issue.Project = &Project{
			ID:   i.Project.ID,
			Name: i.Project.Name,
		}
	}

	if i.Cycle != nil {
		issue.Cycle = &Cycle{
			ID:     i.Cycle.ID,
			Name:   i.Cycle.Name,
			Number: i.Cycle.Number,
		}
	}

	for _, l := range i.Labels.Nodes {
		issue.Labels = append(issue.Labels, Label{
			ID:    l.ID,
			Name:  l.Name,
			Color: l.Color,
		})
	}

	return issue, nil
}

// searchIssuesResponse is the response structure for issue search
type searchIssuesResponse struct {
	Issues struct {
		Nodes []struct {
			ID         string `json:"id"`
			Identifier string `json:"identifier"`
			Title      string `json:"title"`
			Priority   int    `json:"priority"`
			URL        string `json:"url"`
			State      *struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"state"`
			Assignee *struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"assignee"`
			Team struct {
				ID  string `json:"id"`
				Key string `json:"key"`
			} `json:"team"`
		} `json:"nodes"`
	} `json:"issues"`
}

// SearchIssues searches for issues matching the query in title
func (c *LinearClient) SearchIssues(ctx context.Context, query string, opts IssueListOptions) ([]Issue, error) {
	first := opts.First
	if first == 0 {
		first = 50
	}

	// Use proper GraphQL variables to prevent injection
	// The filter is passed as a variable, not concatenated into the query
	rawQuery := `
		query SearchIssues($first: Int!, $filter: IssueFilter!) {
			issues(
				filter: $filter
				first: $first
			) {
				nodes {
					id
					identifier
					title
					priority
					url
					state { id name }
					assignee { id name }
					team { id key }
				}
			}
		}
	`

	// Build filter as a properly typed variable
	filter := map[string]interface{}{
		"title": map[string]interface{}{
			"containsIgnoreCase": query,
		},
	}

	var result searchIssuesResponse
	err := c.gql.Exec(ctx, rawQuery, &result, map[string]interface{}{
		"first":  first,
		"filter": filter,
	})
	if err != nil {
		return nil, fmt.Errorf("search issues: %w", err)
	}

	issues := make([]Issue, 0, len(result.Issues.Nodes))
	for _, i := range result.Issues.Nodes {
		issue := Issue{
			ID:         i.ID,
			Identifier: i.Identifier,
			Title:      i.Title,
			Priority:   i.Priority,
			URL:        i.URL,
			Team: &Team{
				ID:  i.Team.ID,
				Key: i.Team.Key,
			},
		}

		if i.State != nil {
			issue.State = &WorkflowState{
				ID:   i.State.ID,
				Name: i.State.Name,
			}
		}

		if i.Assignee != nil {
			issue.Assignee = &User{
				ID:   i.Assignee.ID,
				Name: i.Assignee.Name,
			}
		}

		issues = append(issues, issue)
	}
	return issues, nil
}

// GetProjects returns projects with optional filters
func (c *LinearClient) GetProjects(ctx context.Context, opts ProjectListOptions) ([]Project, error) {
	first := opts.First
	if first == 0 {
		first = 50
	}

	var query struct {
		Projects struct {
			Nodes []struct {
				ID          string  `graphql:"id"`
				Name        string  `graphql:"name"`
				Description string  `graphql:"description"`
				State       string  `graphql:"state"`
				Progress    float64 `graphql:"progress"`
				TargetDate  *string `graphql:"targetDate"`
				StartDate   *string `graphql:"startDate"`
				URL         string  `graphql:"url"`
				CreatedAt   string  `graphql:"createdAt"`
				UpdatedAt   string  `graphql:"updatedAt"`
				Lead        *struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
				} `graphql:"lead"`
				Teams struct {
					Nodes []struct {
						ID   string `graphql:"id"`
						Name string `graphql:"name"`
						Key  string `graphql:"key"`
					} `graphql:"nodes"`
				} `graphql:"teams"`
			} `graphql:"nodes"`
		} `graphql:"projects(first: $first)"`
	}

	vars := map[string]interface{}{
		"first": graphql.Int(first),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get projects: %w", err)
	}

	projects := make([]Project, 0, len(query.Projects.Nodes))
	for _, p := range query.Projects.Nodes {
		// Apply filters
		if opts.State != nil && p.State != *opts.State {
			continue
		}

		project := Project{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			State:       p.State,
			Progress:    p.Progress,
			TargetDate:  p.TargetDate,
			StartDate:   p.StartDate,
			URL:         p.URL,
		}

		if p.Lead != nil {
			project.Lead = &User{
				ID:   p.Lead.ID,
				Name: p.Lead.Name,
			}
		}

		for _, t := range p.Teams.Nodes {
			project.Teams = append(project.Teams, Team{
				ID:   t.ID,
				Name: t.Name,
				Key:  t.Key,
			})
		}

		// Filter by team if specified
		if opts.TeamID != nil {
			found := false
			for _, t := range project.Teams {
				if t.ID == *opts.TeamID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		projects = append(projects, project)
	}
	return projects, nil
}

// GetProject returns a single project by ID
func (c *LinearClient) GetProject(ctx context.Context, id string) (*Project, error) {
	var query struct {
		Project struct {
			ID          string  `graphql:"id"`
			Name        string  `graphql:"name"`
			Description string  `graphql:"description"`
			State       string  `graphql:"state"`
			Progress    float64 `graphql:"progress"`
			TargetDate  *string `graphql:"targetDate"`
			StartDate   *string `graphql:"startDate"`
			URL         string  `graphql:"url"`
			CreatedAt   string  `graphql:"createdAt"`
			UpdatedAt   string  `graphql:"updatedAt"`
			Lead        *struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"lead"`
			Teams struct {
				Nodes []struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
					Key  string `graphql:"key"`
				} `graphql:"nodes"`
			} `graphql:"teams"`
		} `graphql:"project(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(id),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get project: %w", err)
	}

	p := query.Project
	project := &Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		State:       p.State,
		Progress:    p.Progress,
		TargetDate:  p.TargetDate,
		StartDate:   p.StartDate,
		URL:         p.URL,
	}

	if p.Lead != nil {
		project.Lead = &User{
			ID:   p.Lead.ID,
			Name: p.Lead.Name,
		}
	}

	for _, t := range p.Teams.Nodes {
		project.Teams = append(project.Teams, Team{
			ID:   t.ID,
			Name: t.Name,
			Key:  t.Key,
		})
	}

	return project, nil
}

// GetInitiatives returns all initiatives
func (c *LinearClient) GetInitiatives(ctx context.Context) ([]Initiative, error) {
	var query struct {
		Initiatives struct {
			Nodes []struct {
				ID          string  `graphql:"id"`
				Name        string  `graphql:"name"`
				Description string  `graphql:"description"`
				TargetDate  *string `graphql:"targetDate"`
				CreatedAt   string  `graphql:"createdAt"`
				UpdatedAt   string  `graphql:"updatedAt"`
				Owner       *struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
				} `graphql:"owner"`
			} `graphql:"nodes"`
		} `graphql:"initiatives(first: 50)"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get initiatives: %w", err)
	}

	initiatives := make([]Initiative, len(query.Initiatives.Nodes))
	for i, init := range query.Initiatives.Nodes {
		initiatives[i] = Initiative{
			ID:          init.ID,
			Name:        init.Name,
			Description: init.Description,
			TargetDate:  init.TargetDate,
		}

		if init.Owner != nil {
			initiatives[i].Owner = &User{
				ID:   init.Owner.ID,
				Name: init.Owner.Name,
			}
		}
	}
	return initiatives, nil
}

// GetInitiative returns a single initiative by ID
func (c *LinearClient) GetInitiative(ctx context.Context, id string) (*Initiative, error) {
	var query struct {
		Initiative struct {
			ID          string  `graphql:"id"`
			Name        string  `graphql:"name"`
			Description string  `graphql:"description"`
			TargetDate  *string `graphql:"targetDate"`
			CreatedAt   string  `graphql:"createdAt"`
			UpdatedAt   string  `graphql:"updatedAt"`
			Owner       *struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
			} `graphql:"owner"`
			Projects struct {
				Nodes []struct {
					ID    string `graphql:"id"`
					Name  string `graphql:"name"`
					State string `graphql:"state"`
				} `graphql:"nodes"`
			} `graphql:"projects"`
		} `graphql:"initiative(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(id),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get initiative: %w", err)
	}

	i := query.Initiative
	initiative := &Initiative{
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		TargetDate:  i.TargetDate,
	}

	if i.Owner != nil {
		initiative.Owner = &User{
			ID:   i.Owner.ID,
			Name: i.Owner.Name,
		}
	}

	for _, p := range i.Projects.Nodes {
		initiative.Projects = append(initiative.Projects, Project{
			ID:    p.ID,
			Name:  p.Name,
			State: p.State,
		})
	}

	return initiative, nil
}

// GetCycles returns cycles, optionally filtered by team
func (c *LinearClient) GetCycles(ctx context.Context, teamID *string) ([]Cycle, error) {
	var query struct {
		Cycles struct {
			Nodes []struct {
				ID          string  `graphql:"id"`
				Name        string  `graphql:"name"`
				Number      int     `graphql:"number"`
				StartsAt    string  `graphql:"startsAt"`
				EndsAt      string  `graphql:"endsAt"`
				Progress    float64 `graphql:"progress"`
				Description string  `graphql:"description"`
				Team        struct {
					ID   string `graphql:"id"`
					Name string `graphql:"name"`
					Key  string `graphql:"key"`
				} `graphql:"team"`
			} `graphql:"nodes"`
		} `graphql:"cycles(first: 50)"`
	}

	if err := c.gql.Query(ctx, &query, nil); err != nil {
		return nil, fmt.Errorf("get cycles: %w", err)
	}

	cycles := make([]Cycle, 0, len(query.Cycles.Nodes))
	for _, cy := range query.Cycles.Nodes {
		// Filter by team if specified
		if teamID != nil && cy.Team.ID != *teamID {
			continue
		}

		cycles = append(cycles, Cycle{
			ID:          cy.ID,
			Name:        cy.Name,
			Number:      cy.Number,
			Progress:    cy.Progress,
			Description: cy.Description,
			Team: &Team{
				ID:   cy.Team.ID,
				Name: cy.Team.Name,
				Key:  cy.Team.Key,
			},
		})
	}
	return cycles, nil
}

// GetActiveCycle returns the currently active cycle for a team
func (c *LinearClient) GetActiveCycle(ctx context.Context, teamID string) (*Cycle, error) {
	var query struct {
		Team struct {
			ActiveCycle *struct {
				ID          string  `graphql:"id"`
				Name        string  `graphql:"name"`
				Number      int     `graphql:"number"`
				StartsAt    string  `graphql:"startsAt"`
				EndsAt      string  `graphql:"endsAt"`
				Progress    float64 `graphql:"progress"`
				Description string  `graphql:"description"`
			} `graphql:"activeCycle"`
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
			Key  string `graphql:"key"`
		} `graphql:"team(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(teamID),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get active cycle: %w", err)
	}

	if query.Team.ActiveCycle == nil {
		return nil, ErrNoActiveCycle
	}

	ac := query.Team.ActiveCycle
	return &Cycle{
		ID:          ac.ID,
		Name:        ac.Name,
		Number:      ac.Number,
		Progress:    ac.Progress,
		Description: ac.Description,
		Team: &Team{
			ID:   query.Team.ID,
			Name: query.Team.Name,
			Key:  query.Team.Key,
		},
	}, nil
}

// GetCycle returns a single cycle by ID
func (c *LinearClient) GetCycle(ctx context.Context, id string) (*Cycle, error) {
	var query struct {
		Cycle struct {
			ID          string  `graphql:"id"`
			Name        string  `graphql:"name"`
			Number      int     `graphql:"number"`
			StartsAt    string  `graphql:"startsAt"`
			EndsAt      string  `graphql:"endsAt"`
			Progress    float64 `graphql:"progress"`
			Description string  `graphql:"description"`
			Team        struct {
				ID   string `graphql:"id"`
				Name string `graphql:"name"`
				Key  string `graphql:"key"`
			} `graphql:"team"`
		} `graphql:"cycle(id: $id)"`
	}

	vars := map[string]interface{}{
		"id": graphql.String(id),
	}

	if err := c.gql.Query(ctx, &query, vars); err != nil {
		return nil, fmt.Errorf("get cycle: %w", err)
	}

	cy := query.Cycle
	return &Cycle{
		ID:          cy.ID,
		Name:        cy.Name,
		Number:      cy.Number,
		Progress:    cy.Progress,
		Description: cy.Description,
		Team: &Team{
			ID:   cy.Team.ID,
			Name: cy.Team.Name,
			Key:  cy.Team.Key,
		},
	}, nil
}
