package cmd

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/app"
	cgtools "github.com/ainsleydev/webkit/internal/scaffold"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	err := update(context.Background(), commandInput{
		FS:        afero.NewMemMapFs(),
		AppDef:    app.Definition{},
		Command:   nil,
		Generator: cgtools.Generator{},
	})

	assert.NoError(t, err)
}
