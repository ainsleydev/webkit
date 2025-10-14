package log

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	webkitctx "github.com/ainsleydev/webkit/pkg/context"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestJsonHandler(t *testing.T) {
	ctx := context.WithValue(context.Background(), webkitctx.ContextKeyRequestID, "12345")

	err := os.Setenv(env.AppEnvironmentKey, env.Production.String())
	require.NoError(t, err)

	Bootstrap("test")

	slog.InfoContext(ctx, "Test", "id", 1)
}
