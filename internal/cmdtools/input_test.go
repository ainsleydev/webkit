package cmdtools

import (
	"context"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/pkg/env"
)

// Note: Testing the error path (missing app.json) is not possible because AppDef()
// calls os.Exit(1) directly. This would require a subprocess testing pattern which
// is beyond the scope of these unit tests.
func TestCommandInput_AppDef(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "app.json", []byte(`{
			"project": {
				"name": "test-app",
				"repo": {
					"owner": "test",
					"name": "repo"
				}
			}
		}`), 0o644)
		require.NoError(t, err)

		input := CommandInput{
			FS:       fs,
			Manifest: manifest.NewTracker(),
		}

		def := input.AppDef()
		assert.NotNil(t, def)
		assert.Equal(t, "test-app", def.Project.Name)
		assert.Equal(t, "test", def.Project.Repo.Owner)
	})

	t.Run("Caching", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, "app.json", []byte(`{
			"project": {
				"name": "cached-app",
				"repo": {
					"owner": "cached",
					"name": "repo"
				}
			}
		}`), 0o644)
		require.NoError(t, err)

		input := CommandInput{
			FS:       fs,
			Manifest: manifest.NewTracker(),
		}

		def1 := input.AppDef()
		require.NotNil(t, def1)

		err = afero.WriteFile(fs, "app.json", []byte(`{
			"project": {
				"name": "modified-app",
				"repo": {
					"owner": "modified",
					"name": "repo"
				}
			}
		}`), 0o644)
		require.NoError(t, err)

		def2 := input.AppDef()
		assert.Same(t, def1, def2)
		assert.Equal(t, "cached-app", def2.Project.Name)
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("Production mode", func(t *testing.T) {
		t.Parallel()

		// Ensure we're not in development mode.
		oldEnv := os.Getenv(env.AppEnvironmentKey)
		defer func() {
			if oldEnv != "" {
				_ = os.Setenv(env.AppEnvironmentKey, oldEnv)
			} else {
				_ = os.Unsetenv(env.AppEnvironmentKey)
			}
		}()
		_ = os.Unsetenv(env.AppEnvironmentKey)

		called := false
		command := func(ctx context.Context, input CommandInput) error {
			called = true
			assert.NotNil(t, input.FS)
			assert.NotNil(t, input.Command)
			assert.Equal(t, "./", input.BaseDir)
			assert.NotNil(t, input.Manifest)
			return nil
		}

		wrappedFunc := Wrap(command)
		err := wrappedFunc(context.Background(), &cli.Command{})

		require.NoError(t, err)
		assert.True(t, called)
	})

	t.Run("Development mode", func(t *testing.T) {
		t.Parallel()

		oldEnv := os.Getenv(env.AppEnvironmentKey)
		defer func() {
			if oldEnv != "" {
				_ = os.Setenv(env.AppEnvironmentKey, oldEnv)
			} else {
				_ = os.Unsetenv(env.AppEnvironmentKey)
			}
		}()
		_ = os.Setenv(env.AppEnvironmentKey, env.Development.String())

		called := false
		command := func(ctx context.Context, input CommandInput) error {
			called = true
			assert.NotNil(t, input.FS)
			assert.Equal(t, "./internal/playground", input.BaseDir)
			return nil
		}

		wrappedFunc := Wrap(command)
		err := wrappedFunc(context.Background(), &cli.Command{})

		require.NoError(t, err)
		assert.True(t, called)
	})
}

func TestCommandInput_Generator(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()
	input := CommandInput{
		FS:       fs,
		Manifest: manifest.NewTracker(),
	}

	gen := input.Generator()
	assert.NotNil(t, gen)
}

func TestCommandInput_Printer(t *testing.T) {
	t.Parallel()

	t.Run("Creates printer", func(t *testing.T) {
		t.Parallel()

		input := CommandInput{}
		printer := input.Printer()
		assert.NotNil(t, printer)
	})

	t.Run("Caches printer", func(t *testing.T) {
		t.Parallel()

		input := CommandInput{}
		printer1 := input.Printer()
		printer2 := input.Printer()

		assert.Same(t, printer1, printer2)
	})
}

func TestCommandInput_SOPSClient(t *testing.T) {
	t.Parallel()

	t.Run("Returns cached client", func(t *testing.T) {
		t.Parallel()

		input := CommandInput{}

		client := input.SOPSClient()
		assert.NotNil(t, client)

		client2 := input.SOPSClient()
		assert.Same(t, client, client2)
	})
}

func TestCommandInput_Spinner(t *testing.T) {
	t.Parallel()

	input := CommandInput{}
	spinner := input.Spinner()

	assert.NotNil(t, spinner)
}
