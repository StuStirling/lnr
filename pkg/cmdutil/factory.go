package cmdutil

import (
	"github.com/stustirling/lnr/internal/api"
	"github.com/stustirling/lnr/internal/config"
	"github.com/stustirling/lnr/internal/output"
)

// Factory provides dependencies for commands
type Factory struct {
	Config    *config.Config
	Client    api.Client
	Formatter *output.Formatter
}

// NewFactory creates a new factory with dependencies
func NewFactory(jsonOutput bool) (*Factory, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	client := api.NewClient(cfg.APIKey)
	formatter := output.NewFormatter(jsonOutput)

	return &Factory{
		Config:    cfg,
		Client:    client,
		Formatter: formatter,
	}, nil
}

// NewFactoryWithClient creates a factory with a custom client (for testing)
func NewFactoryWithClient(client api.Client, jsonOutput bool) *Factory {
	return &Factory{
		Config:    &config.Config{},
		Client:    client,
		Formatter: output.NewFormatter(jsonOutput),
	}
}
