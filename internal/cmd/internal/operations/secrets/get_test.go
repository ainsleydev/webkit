package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestGet(t *testing.T) {
	ctx := t.Context()

	t.Run("Client Error", func(t *testing.T) {
		input := cmdtools.CommandInput{Command: GetCmd}

		err := Encrypt(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "age")
	})

	t.Run("Decode Error", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY: "1234"\ninvalid`)

		err := Get(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "decoding sops to map")
	})

	t.Run("No Value", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY: "1234"`)

		err := Encrypt(ctx, input)
		require.NoError(t, err)

		require.NoError(t, input.Command.Set("env", env.Production))
		require.NoError(t, input.Command.Set("key", "wrong"))

		err = Get(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "key wrong not found")
	})

	t.Run("Success", func(t *testing.T) {
		input := setupEncryptedProdFile(t, `KEY: "1234"`)

		err := Encrypt(ctx, input)
		require.NoError(t, err)

		require.NoError(t, input.Command.Set("env", env.Production))
		require.NoError(t, input.Command.Set("key", "KEY"))

		err = Get(ctx, input)
		assert.NoError(t, err)
		// TODO: Assert that output has value.
	})
}
