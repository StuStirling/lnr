package api

import "context"

// MockClient is a mock implementation of the Client interface for testing
type MockClient struct {
	GetViewerFunc         func(ctx context.Context) (*User, error)
	GetOrganisationFunc   func(ctx context.Context) (*Organisation, error)
	GetUsersFunc          func(ctx context.Context) ([]User, error)
	GetTeamsFunc          func(ctx context.Context) ([]Team, error)
	GetTeamFunc           func(ctx context.Context, id string) (*Team, error)
	GetLabelsFunc         func(ctx context.Context, teamID *string) ([]Label, error)
	GetWorkflowStatesFunc func(ctx context.Context, teamID *string) ([]WorkflowState, error)
	GetIssuesFunc         func(ctx context.Context, opts IssueListOptions) ([]Issue, error)
	GetIssueFunc          func(ctx context.Context, id string) (*Issue, error)
	SearchIssuesFunc      func(ctx context.Context, query string, opts IssueListOptions) ([]Issue, error)
	GetProjectsFunc       func(ctx context.Context, opts ProjectListOptions) ([]Project, error)
	GetProjectFunc        func(ctx context.Context, id string) (*Project, error)
	GetInitiativesFunc    func(ctx context.Context) ([]Initiative, error)
	GetInitiativeFunc     func(ctx context.Context, id string) (*Initiative, error)
	GetCyclesFunc         func(ctx context.Context, teamID *string) ([]Cycle, error)
	GetActiveCycleFunc    func(ctx context.Context, teamID string) (*Cycle, error)
	GetCycleFunc          func(ctx context.Context, id string) (*Cycle, error)
}

func (m *MockClient) GetViewer(ctx context.Context) (*User, error) {
	if m.GetViewerFunc != nil {
		return m.GetViewerFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetOrganisation(ctx context.Context) (*Organisation, error) {
	if m.GetOrganisationFunc != nil {
		return m.GetOrganisationFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetUsers(ctx context.Context) ([]User, error) {
	if m.GetUsersFunc != nil {
		return m.GetUsersFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetTeams(ctx context.Context) ([]Team, error) {
	if m.GetTeamsFunc != nil {
		return m.GetTeamsFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetTeam(ctx context.Context, id string) (*Team, error) {
	if m.GetTeamFunc != nil {
		return m.GetTeamFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClient) GetLabels(ctx context.Context, teamID *string) ([]Label, error) {
	if m.GetLabelsFunc != nil {
		return m.GetLabelsFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockClient) GetWorkflowStates(ctx context.Context, teamID *string) ([]WorkflowState, error) {
	if m.GetWorkflowStatesFunc != nil {
		return m.GetWorkflowStatesFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockClient) GetIssues(ctx context.Context, opts IssueListOptions) ([]Issue, error) {
	if m.GetIssuesFunc != nil {
		return m.GetIssuesFunc(ctx, opts)
	}
	return nil, nil
}

func (m *MockClient) GetIssue(ctx context.Context, id string) (*Issue, error) {
	if m.GetIssueFunc != nil {
		return m.GetIssueFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClient) SearchIssues(ctx context.Context, query string, opts IssueListOptions) ([]Issue, error) {
	if m.SearchIssuesFunc != nil {
		return m.SearchIssuesFunc(ctx, query, opts)
	}
	return nil, nil
}

func (m *MockClient) GetProjects(ctx context.Context, opts ProjectListOptions) ([]Project, error) {
	if m.GetProjectsFunc != nil {
		return m.GetProjectsFunc(ctx, opts)
	}
	return nil, nil
}

func (m *MockClient) GetProject(ctx context.Context, id string) (*Project, error) {
	if m.GetProjectFunc != nil {
		return m.GetProjectFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClient) GetInitiatives(ctx context.Context) ([]Initiative, error) {
	if m.GetInitiativesFunc != nil {
		return m.GetInitiativesFunc(ctx)
	}
	return nil, nil
}

func (m *MockClient) GetInitiative(ctx context.Context, id string) (*Initiative, error) {
	if m.GetInitiativeFunc != nil {
		return m.GetInitiativeFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockClient) GetCycles(ctx context.Context, teamID *string) ([]Cycle, error) {
	if m.GetCyclesFunc != nil {
		return m.GetCyclesFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockClient) GetActiveCycle(ctx context.Context, teamID string) (*Cycle, error) {
	if m.GetActiveCycleFunc != nil {
		return m.GetActiveCycleFunc(ctx, teamID)
	}
	return nil, nil
}

func (m *MockClient) GetCycle(ctx context.Context, id string) (*Cycle, error) {
	if m.GetCycleFunc != nil {
		return m.GetCycleFunc(ctx, id)
	}
	return nil, nil
}
