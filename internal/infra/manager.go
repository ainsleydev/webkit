package infra

import (
	"context"

	"github.com/ainsleydev/webkit/pkg/env"
)

//go:generate go tool go.uber.org/mock/mockgen -source=manager.go -destination ./mocks/infra.go -package=mockinfra

type (
	// Manager defines the interface for managing infrastructure operations.
	Manager interface {
		Init(ctx context.Context) error
		Plan(ctx context.Context, env env.Environment) (PlanOutput, error)
		Apply(ctx context.Context, env env.Environment) (ApplyOutput, error)
		Destroy(ctx context.Context, env env.Environment) (DestroyOutput, error)
		Output(ctx context.Context, env env.Environment) (OutputResult, error)
		Import(ctx context.Context, input ImportInput) (ImportOutput, error)
		Cleanup()
	}

	// ImportInput contains the configuration for importing existing resources.
	ImportInput struct {
		// ResourceName is the name of the resource in app.json.
		ResourceName string
		// ResourceID is the provider-specific ID (e.g., DigitalOcean cluster ID).
		ResourceID string
		// Environment specifies which environment to import into.
		Environment env.Environment
	}

	// ImportOutput contains the results of an import operation.
	ImportOutput struct {
		// ImportedResources lists the Terraform addresses that were imported.
		ImportedResources []string
		// Output contains the human-readable output from the import operations.
		Output string
	}
)
