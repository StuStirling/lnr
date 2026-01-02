package issue

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stustirling/lnr/internal/api"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

func TestRunListWithFactory(t *testing.T) {
	tests := []struct {
		name        string
		mockIssues  []api.Issue
		mockError   error
		opts        api.IssueListOptions
		wantErr     bool
		wantOutput  string
	}{
		{
			name: "lists issues successfully",
			mockIssues: []api.Issue{
				{
					Identifier: "ENG-123",
					Title:      "Fix login bug",
					Priority:   2,
					State:      &api.WorkflowState{Name: "In Progress"},
					Assignee:   &api.User{Name: "John Doe"},
				},
				{
					Identifier: "ENG-124",
					Title:      "Add dark mode",
					Priority:   3,
					State:      &api.WorkflowState{Name: "Backlog"},
					Assignee:   nil,
				},
			},
			opts:       api.IssueListOptions{First: 50},
			wantErr:    false,
			wantOutput: "ENG-123",
		},
		{
			name:       "handles empty results",
			mockIssues: []api.Issue{},
			opts:       api.IssueListOptions{First: 50},
			wantErr:    false,
			wantOutput: "No results found",
		},
		{
			name:      "handles API error",
			mockError: assert.AnError,
			opts:      api.IssueListOptions{First: 50},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock client
			mockClient := &api.MockClient{
				GetIssuesFunc: func(ctx context.Context, opts api.IssueListOptions) ([]api.Issue, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockIssues, nil
				},
			}

			// Create factory with mock client
			factory := cmdutil.NewFactoryWithClient(mockClient, false)

			// Capture output
			var buf bytes.Buffer
			factory.Formatter.SetWriter(&buf)

			// Run the command
			err := runListWithFactory(factory, tt.opts)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Contains(t, buf.String(), tt.wantOutput)
		})
	}
}

func TestRunListWithFactory_WithFilters(t *testing.T) {
	teamID := "team-123"
	assigneeID := "user-456"

	mockClient := &api.MockClient{
		GetIssuesFunc: func(ctx context.Context, opts api.IssueListOptions) ([]api.Issue, error) {
			// Verify filters are passed correctly
			assert.Equal(t, &teamID, opts.TeamID)
			assert.Equal(t, &assigneeID, opts.AssigneeID)
			return []api.Issue{
				{
					Identifier: "ENG-123",
					Title:      "Test Issue",
					Priority:   1,
					State:      &api.WorkflowState{Name: "Todo"},
				},
			}, nil
		},
	}

	factory := cmdutil.NewFactoryWithClient(mockClient, false)
	var buf bytes.Buffer
	factory.Formatter.SetWriter(&buf)

	opts := api.IssueListOptions{
		TeamID:     &teamID,
		AssigneeID: &assigneeID,
		First:      50,
	}

	err := runListWithFactory(factory, opts)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "ENG-123")
}

func TestRunListWithFactory_JSONOutput(t *testing.T) {
	mockClient := &api.MockClient{
		GetIssuesFunc: func(ctx context.Context, opts api.IssueListOptions) ([]api.Issue, error) {
			return []api.Issue{
				{
					ID:         "issue-1",
					Identifier: "ENG-123",
					Title:      "Test Issue",
					Priority:   1,
				},
			}, nil
		},
	}

	factory := cmdutil.NewFactoryWithClient(mockClient, true) // JSON output
	var buf bytes.Buffer
	factory.Formatter.SetWriter(&buf)

	err := runListWithFactory(factory, api.IssueListOptions{First: 50})
	require.NoError(t, err)

	// Verify JSON output contains expected fields
	output := buf.String()
	assert.Contains(t, output, `"identifier"`)
	assert.Contains(t, output, `"ENG-123"`)
}
