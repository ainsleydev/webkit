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
	//redirect()

	Head(context.Context) markup.HeadProps
	Robots() webkit.Handler
	Sitemap() webkit.Handler
}

func Scratch() {
	app := webkit.New()
	p := &payload.Adapter{}

	app.Get("/robots.txt", p.Robots())
	app.Get("/sitemap.xml", p.Sitemap())
}

func PayloadScratch() {
	app := webkit.New()
	p, _ := payload.NewAdapter(
		payload.WithBaseURL("https://api.payloadcms.com"),
		payload.WithAPIKey(""),
	)

	app.Get("/robots.txt", p.Robots())
	app.Get("/sitemap.xml", p.Sitemap())
}
