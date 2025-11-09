package files

import (
	"context"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// Schema copies the JSON schema to the .webkit folder.
func Schema(_ context.Context, input cmdtools.CommandInput) error {
	return input.Generator().Template(filepath.Join(".webkit", "schema.json"),
		templates.MustLoadTemplate("schema.json"),
		nil,
		scaffold.WithTracking(manifest.SourceProject()),
		scaffold.WithoutNotice(),
	)
}
