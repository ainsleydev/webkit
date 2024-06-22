package adapters

import (
	"context"

	"github.com/ainsleydev/webkit/pkg/markup"
)

// Adapter for on different platforms such as Payload & Static
type Adapter interface {
	Head(context.Context) markup.HeadProps
}
