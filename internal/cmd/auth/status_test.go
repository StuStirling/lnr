package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stustirling/lnr/internal/api"
	"github.com/stustirling/lnr/pkg/cmdutil"
)

func TestRunStatusWithFactory(t *testing.T) {
	tests := []struct {
		name       string
		mockUser   *api.User
		mockOrg    *api.Organisation
		mockError  error
		jsonOutput bool
		wantErr    bool
	}{
		{
			name: "authenticates successfully",
			mockUser: &api.User{
				Name:  "Test User",
				Email: "test@example.com",
				Admin: false,
			},
			mockOrg: &api.Organisation{
				Name:      "Test Org",
				UserCount: 10,
			},
			jsonOutput: false,
			wantErr:    false,
		},
		{
			name: "authenticates admin successfully",
			mockUser: &api.User{
				Name:  "Admin User",
				Email: "admin@example.com",
				Admin: true,
			},
			mockOrg: &api.Organisation{
				Name:      "Test Org",
				UserCount: 5,
			},
			jsonOutput: false,
			wantErr:    false,
		},
		{
			name:       "handles API error",
			mockError:  assert.AnError,
			jsonOutput: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &api.MockClient{
				GetViewerFunc: func(ctx context.Context) (*api.User, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockUser, nil
				},
				GetOrganisationFunc: func(ctx context.Context) (*api.Organisation, error) {
					if tt.mockError != nil {
						return nil, tt.mockError
					}
					return tt.mockOrg, nil
				},
			}

			factory := cmdutil.NewFactoryWithClient(mockClient, tt.jsonOutput)
			var buf bytes.Buffer
			factory.Formatter.SetWriter(&buf)

			err := runStatusWithFactory(factory, tt.jsonOutput)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestRunStatusWithFactory_JSONOutput(t *testing.T) {
	mockClient := &api.MockClient{
		GetViewerFunc: func(ctx context.Context) (*api.User, error) {
			return &api.User{
				ID:    "user-1",
				Name:  "Test User",
				Email: "test@example.com",
				Admin: false,
			}, nil
		},
		GetOrganisationFunc: func(ctx context.Context) (*api.Organisation, error) {
			return &api.Organisation{
				ID:        "org-1",
				Name:      "Test Org",
				UserCount: 10,
			}, nil
		},
	}

	factory := cmdutil.NewFactoryWithClient(mockClient, true)
	var buf bytes.Buffer
	factory.Formatter.SetWriter(&buf)

	err := runStatusWithFactory(factory, true)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `"authenticated"`)
	assert.Contains(t, output, `"user"`)
	assert.Contains(t, output, `"organisation"`)
}
