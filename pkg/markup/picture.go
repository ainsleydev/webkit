package markup

import (
	"context"
	"html/template"
	"io"

	"github.com/ainsleydev/webkit/pkg/tpl"
)

// PictureProvider is a common - TODO
type PictureProvider interface {
	ToMarkup(ctx context.Context) PictureProps
}

func Picture(ctx context.Context, provider PictureProvider) PictureProps {
	return provider.ToMarkup(ctx)
}

// PictureProps defines the fields for to render a <picture> element onto the DOM.
//
// The <picture> HTML element contains zero or more <source> elements and one <img> element
// to offer alternative versions of an image for different display/device scenarios.
//
// See: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/picture
type PictureProps struct {
	// The URL of the image, which will map to the src attribute.
	// For example: "/images/dog.jpg"
	URL string

	// Defines text that can replace the image in the page.
	Alt string

	// Maps to <source> elements within the <picture> element.
	// The browser will consider each child <source> element and choose the best match among them.
	Sources []ImageProps

	// Determines if loading=lazy should be added to the image.
	Lazy LoadingAttribute

	// List of class names to apply to the <picture> element.
	Classes []string

	// A unique identifier for the <picture> element.
	ID string

	// The intrinsic width of the image in pixels , for example (300).
	// Must be an integer without a unit (optional).
	Width *int

	// The intrinsic height of the image, in pixels, for example (300).
	// Must be an integer without a unit (optional).
	Height *int

	// Attributes specifies additional attributes for the picture element as key-value pairs.
	// For example: markup.Attributes{"data-attribute-size": "large"}
	Attributes Attributes
}

// LoadingAttribute specifies the loading attribute for an image.
// Indicates how the browser should load the image:
type LoadingAttribute string

const (
	// LoadingEager loads the image immediately, regardless of whether or not the
	// image is currently within the visible viewport (this is the default value).
	LoadingEager LoadingAttribute = "eager"
	// LoadingLazy Defers loading the image until it reaches a calculated distance
	// from the viewport, as defined by the browser. The intent is to avoid the
	// network and storage bandwidth needed to handle the image until it's
	// reasonably certain that it will be needed. This generally improves
	// the performance of the content in most typical use cases.
	LoadingLazy LoadingAttribute = "lazy"
)

// headTemplate is the template for the head of the HTML document.
// It requires a HeadProps struct to be passed in when executing the template.
var mediaTemplates = template.Must(template.New("").Funcs(tpl.Funcs).ParseFS(templatesFS,
	"picture.html",
	//"image.html",
))

// Render renders a picture element to the provided writer.
func (p PictureProps) Render(ctx context.Context, w io.Writer) error {
	return mediaTemplates.ExecuteTemplate(w, "picture.html", p)
}

// ImageProps defines the fields for an individual image or source HTML element.
// The data type supports both <img> and <source> elements.
//
// See: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/img
// See: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/source
type ImageProps struct {
	// The URL of the image, which will map to the srcset or src attribute.
	// For example: "/images/dog.jpg"
	URL string

	// Media specifies the media condition (media query) for the source,
	// For example: "(min-width: 600px)"
	Media string

	// Mimetype such as
	// - image/jpeg
	// - image/png
	// - image/gif
	// - image/avif
	// - image/webp
	// - image/svg+xml
	MimeType *string

	// The intrinsic width of the image in pixels , for example (300).
	// Must be an integer without a unit (optional).
	Width *int

	// The intrinsic height of the image, in pixels, for example (300).
	// Must be an integer without a unit (optional).
	Height *int

	// Attributes specifies additional attributes for the picture element as key-value pairs.
	// For example: markup.Attributes{"id": "main-picture", "class": "responsive-picture"}
	Attributes Attributes
}

//////////////////////////////////////////////////////////////////

type PictureOptions func(p *PictureProps)

func PictureWithAlt(alt string) PictureOptions {
	return func(p *PictureProps) {
		p.Alt = alt
	}
}

func PictureWithLoadingLazy() PictureOptions {
	return func(p *PictureProps) {
		p.Lazy = LoadingLazy
	}
}

func PictureWithLoadingAuto() PictureOptions {
	return func(p *PictureProps) {
		p.Lazy = LoadingEager
	}
}
