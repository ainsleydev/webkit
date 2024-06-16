package adapters

import (
	"context"

	"github.com/ainsleydev/webkit/pkg/adapters/payload"
	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/ainsleydev/webkit/pkg/webkit"
)

// Adapter for on different platforms such as Payload & Static
type Adapter interface {
	//Head() string
	//
	//Redirect()

	Head(ctx context.Context) (markup.HeadProps, error)
	Robots() webkit.Handler
	Sitemap() webkit.Handler
}

func Scratch() {
	app := webkit.New()
	p := &payload.Adapter{}

	app.Get("/robots.txt", p.Robots())
	app.Get("/sitemap.xml", p.Sitemap())
}
