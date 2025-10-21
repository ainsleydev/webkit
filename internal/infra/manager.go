package infra

import (
	"context"

	"github.com/ainsleydev/webkit/pkg/env"
)

//go:generate go tool go.uber.org/mock/mockgen -source=manager.go -destination ../mocks/infra.go -package=mocks

type Manager interface {
	Init(ctx context.Context) error
	Plan(ctx context.Context, env env.Environment) (PlanOutput, error)
	Apply(ctx context.Context, env env.Environment) (ApplyOutput, error)
	Destroy(ctx context.Context, env env.Environment) (DestroyOutput, error)
	Output(ctx context.Context, env env.Environment) (OutputResult, error)
	Cleanup()
}
