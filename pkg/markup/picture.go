package markup

import (
	"context"
	"html/template"
	"io"
	"path/filepath"

	"github.com/ainsleydev/webkit/pkg/tpl"
)

// PictureProvider is a common - TODO
type PictureProvider interface {
	ToMarkup() PictureProps
}

// Picture - TODO
//
// TODO: How are we going to apply classes to the picture element?
// At the moment, we are only applying classes to the source elements.
func Picture(provider PictureProvider, opts ...ImageOptions) PictureProps {
	props := provider.ToMarkup()
	props.FileExtension = filepath.Ext(props.URL)

	for i, img := range props.Sources {
		// Assign the file extension to the source images.
		props.Sources[i].FileExtension = filepath.Ext(img.URL)

		// Apply all options to the source images
		for _, opt := range opts {
			opt(&img)
		}

		// Apply upwards
		if img.Alt != "" {
			props.Alt = img.Alt
		}
	}

	return props
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

	// List of class names to apply to the <picture> element.
	Classes []string

	// A unique identifier for the <picture> element.
	ID string

	// The file extension of the image, for example (jpg).
	FileExtension string

	// Determines if loading=lazy should be added to the image.
	Loading LoadingAttribute

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

// Image transforms the PictureProps into an ImageProps type.
//
// This is useful when you want to render a single image element, instead
// of the entire picture.
func (p PictureProps) Image() ImageProps {
	return ImageProps{
		URL:           p.URL,
		Alt:           p.Alt,
		IsSource:      false,
		Media:         "", // Default image should not output a media query.
		MimeType:      "",
		FileExtension: p.FileExtension,
		Loading:       p.Loading,
		Width:         p.Width,
		Height:        p.Height,
		Attributes:    p.Attributes,
	}
}

// headTemplate is the template for the head of the HTML document.
// It requires a HeadProps struct to be passed in when executing the template.
var mediaTemplates = template.Must(template.New("").Funcs(tpl.Funcs).ParseFS(templatesFS,
	"picture.html",
	"image.html",
))

// Render renders a picture element to the provided writer.
func (p PictureProps) Render(_ context.Context, w io.Writer) error {
	return mediaTemplates.ExecuteTemplate(w, "picture.html", p)
}

// PictureOptions allows for optional settings to be applied to a <picture>.
type PictureOptions func(p *PictureProps)

// PictureWithAlt attaches alternative text to the picture.
func PictureWithAlt(alt string) PictureOptions {
	return func(p *PictureProps) {
		p.Alt = alt
	}
}

// PictureWithLazyLoading sets loading=lazy to the picture.
func PictureWithLazyLoading() PictureOptions {
	return func(p *PictureProps) {
		p.Loading = LoadingLazy
	}
}

// PictureWithEagerLoading sets loading=eager to the picture.
func PictureWithEagerLoading() PictureOptions {
	return func(p *PictureProps) {
		p.Loading = LoadingEager
	}
}

func PictureWithClasses(classes ...string) PictureOptions {
	return func(p *PictureProps) {
		for _, v := range classes {
			p.Classes = append(p.Classes, v)
		}
	}
}
