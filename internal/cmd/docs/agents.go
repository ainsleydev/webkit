package docs

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	internaldocs "github.com/ainsleydev/webkit/internal/docs"
	"github.com/ainsleydev/webkit/internal/manifest"
)

const (
	// customContentPath is the path for custom content in service/app repos.
	customContentPath = "docs"
)

// Agents creates the AGENTS.md file at the project root by combining
// the base template with optional custom content from docs/.
func Agents(_ context.Context, input cmdtools.CommandInput) error {
	opts := internaldocs.GenerateAgentsOptions{
		FS:                input.FS,
		Generator:         input.Generator(),
		CustomContentPath: customContentPath,
		OutputPath:        "AGENTS.md",
		TemplateData:      input.AppDef(),
		TrackingSource:    manifest.SourceProject(),
	}

	if err := internaldocs.GenerateAgents(opts); err != nil {
		return errors.Wrap(err, "generating AGENTS.md")
	}

	return nil
}
