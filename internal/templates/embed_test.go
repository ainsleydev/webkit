package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Success README", func(t *testing.T) {
		t.Parallel()

		tmpl, err := LoadTemplate("README.md")
		require.NoError(t, err)
		assert.NotNil(t, tmpl)
		assert.Equal(t, "README.md", tmpl.Name())
	})

	t.Run("Success eslint config", func(t *testing.T) {
		t.Parallel()

		tmpl, err := LoadTemplate("eslint.config.js.tmpl")
		require.NoError(t, err)
		assert.NotNil(t, tmpl)
		assert.Equal(t, "eslint.config.js.tmpl", tmpl.Name())
	})

	t.Run("Template not found", func(t *testing.T) {
		t.Parallel()

		tmpl, err := LoadTemplate("nonexistent.tmpl")
		assert.Error(t, err)
		assert.Nil(t, tmpl)
	})

	t.Run("Custom functions available", func(t *testing.T) {
		t.Parallel()

		tmpl, err := LoadTemplate("README.md")
		require.NoError(t, err)

		funcs := tmpl.Funcs(templateFuncs())
		assert.NotNil(t, funcs)
	})
}

func TestMustLoadTemplate(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tmpl := MustLoadTemplate("README.md")
		assert.NotNil(t, tmpl)
		assert.Equal(t, "README.md", tmpl.Name())
	})

	t.Run("Panic on error", func(t *testing.T) {
		t.Parallel()

		assert.Panics(t, func() {
			MustLoadTemplate("nonexistent.tmpl")
		})
	})
}

func TestTemplateFuncs(t *testing.T) {
	t.Parallel()

	funcs := templateFuncs()

	t.Run("Contains custom functions", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, funcs["ghVar"])
		assert.NotNil(t, funcs["ghSecret"])
		assert.NotNil(t, funcs["ghInput"])
		assert.NotNil(t, funcs["ghEnv"])
	})

	t.Run("Contains sprig functions", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, funcs["upper"])
		assert.NotNil(t, funcs["lower"])
		assert.NotNil(t, funcs["trim"])
	})

	t.Run("Custom functions work correctly", func(t *testing.T) {
		t.Parallel()

		ghVarFunc := funcs["ghVar"].(func(string) string)
		assert.Equal(t, "${{ test }}", ghVarFunc("test"))

		ghSecretFunc := funcs["ghSecret"].(func(string) string)
		assert.Equal(t, "${{ secrets.TOKEN }}", ghSecretFunc("TOKEN"))

		ghInputFunc := funcs["ghInput"].(func(string) string)
		assert.Equal(t, "${{ inputs.env }}", ghInputFunc("env"))

		ghEnvFunc := funcs["ghEnv"].(func(string) string)
		assert.Equal(t, "${{ env.VAR }}", ghEnvFunc("VAR"))
	})
}

func TestEmbed(t *testing.T) {
	t.Parallel()

	t.Run("Embed FS not nil", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, Embed)
	})

	t.Run("Can read from embed", func(t *testing.T) {
		t.Parallel()

		data, err := Embed.ReadFile("README.md")
		require.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}
