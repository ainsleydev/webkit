package markup

import (
	"context"
	"io"
)

// Component is the interface that all templates implement.
type Component interface {
	// Render the template.
	Render(ctx context.Context, w io.Writer) error
}
