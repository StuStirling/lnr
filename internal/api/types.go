package api

import "time"

// User represents a Linear user
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	Admin       bool   `json:"admin"`
	AvatarURL   string `json:"avatarUrl"`
}

// Team represents a Linear team
type Team struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
}

// WorkflowState represents a workflow state in Linear
type WorkflowState struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Type     string `json:"type"`
	Position int    `json:"position"`
	Team     *Team  `json:"team,omitempty"`
}

// Label represents an issue label
type Label struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Team        *Team  `json:"team,omitempty"`
}

// Issue represents a Linear issue
type Issue struct {
	ID          string         `json:"id"`
	Identifier  string         `json:"identifier"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    int            `json:"priority"`
	Estimate    *float64       `json:"estimate"`
	State       *WorkflowState `json:"state"`
	Assignee    *User          `json:"assignee"`
	Creator     *User          `json:"creator"`
	Team        *Team          `json:"team"`
	Project     *Project       `json:"project"`
	Cycle       *Cycle         `json:"cycle"`
	Labels      []Label        `json:"labels"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DueDate     *string        `json:"dueDate"`
	URL         string         `json:"url"`
}

// Project represents a Linear project
type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	State       string    `json:"state"`
	Progress    float64   `json:"progress"`
	TargetDate  *string   `json:"targetDate"`
	StartDate   *string   `json:"startDate"`
	Lead        *User     `json:"lead"`
	Teams       []Team    `json:"teams"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	URL         string    `json:"url"`
}

// Initiative represents a Linear initiative
type Initiative struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TargetDate  *string   `json:"targetDate"`
	Owner       *User     `json:"owner"`
	Projects    []Project `json:"projects"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Cycle represents a Linear cycle (sprint)
type Cycle struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Number      int       `json:"number"`
	StartsAt    time.Time `json:"startsAt"`
	EndsAt      time.Time `json:"endsAt"`
	Progress    float64   `json:"progress"`
	Team        *Team     `json:"team"`
	Description string    `json:"description"`
}

// Organisation represents the Linear organisation
type Organisation struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URLKey    string `json:"urlKey"`
	LogoURL   string `json:"logoUrl"`
	UserCount int    `json:"userCount"`
}

// PageInfo contains pagination information
type PageInfo struct {
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
}
