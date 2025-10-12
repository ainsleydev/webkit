package secrets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/env"
)

func TestGet(t *testing.T) {
	ctx := t.Context()

	t.Run("Decode Error", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY: "1234"\ninvalid`)

		err := Get(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "decoding sops to map")
	})

	t.Run("No Key or All Flag", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY: "1234"`)
		require.NoError(t, input.Command.Set("env", env.Production))

		err := Get(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "either --key or --all must be provided")
	})

	t.Run("Key Not Found", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY: "1234"`)
		require.NoError(t, input.Command.Set("env", env.Production))
		require.NoError(t, input.Command.Set("key", "wrong"))

		err := Get(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "key wrong not found")
	})

	t.Run("Single Key Success", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY: "1234"`)
		buf := &bytes.Buffer{}
		input.Printer().SetWriter(buf)

		require.NoError(t, input.Command.Set("env", env.Production))
		require.NoError(t, input.Command.Set("key", "KEY"))

		err := Get(ctx, input)
		require.NoError(t, err)

		out := buf.String()
		assert.Contains(t, out, "KEY=1234")
	})

	t.Run("All Keys Success", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY1: "1234"
KEY2: "abcd"`)

		buf := &bytes.Buffer{}
		input.Printer().SetWriter(buf)

		require.NoError(t, input.Command.Set("env", env.Production))
		require.NoError(t, input.Command.Set("all", "true"))

		err := Get(ctx, input)
		require.NoError(t, err)

		out := buf.String()
		assert.Contains(t, out, "KEY1: 1234")
		assert.Contains(t, out, "KEY2: abcd")
	})
}
