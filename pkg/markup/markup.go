package markup

import (
	"context"
	"io"
)

type Component interface {
	Render(context.Context, io.Writer) error
}
