package markup

import (
	"context"
	"embed"
	"io"
)

//go:embed *.html
var templatesFS embed.FS

type Component interface {
	Render(context.Context, io.Writer) error
}
