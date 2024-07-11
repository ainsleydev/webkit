package markup

import (
	"context"
	"html/template"
	"io"

	"github.com/ainsleydev/webkit/pkg/tpl"
)

type PictureProps struct {
	Images     []Image
	Lazy       bool
	Alt        string
	Classes    []string
	Attributes []string
}

type Image struct {
	URL        string
	Media      string // Media query for responsive images
	Type       string // Mimetype such as image/jpeg or image/webp
	Width      int
	Height     int
	Attributes []string
}

type PictureProvider interface {
	Data(ctx context.Context) PictureProps
}

func Picture(ctx context.Context, provider PictureProvider) PictureProps {
	return provider.Data(ctx)
}

// headTemplate is the template for the head of the HTML document.
// It requires a HeadProps struct to be passed in when executing the template.
var pictureTemplate = template.Must(template.New("").Funcs(tpl.Funcs).ParseFS(templatesFS,
	"picture.html",
))

// Render renders a picture element to the provided writer.
func (p PictureProps) Render(ctx context.Context, w io.Writer) error {
	return pictureTemplate.ExecuteTemplate(w, "picture.html", p)
}

type Static string

func (s Static) Data(ctx context.Context) PictureProps {
	return PictureProps{
		Images: []Image{
			{
				URL: string(s),
			},
		},
	}
}
