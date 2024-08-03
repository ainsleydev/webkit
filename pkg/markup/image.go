package markup

import (
	"context"
	"github.com/ainsleydev/webkit/pkg/markup/internal/templates"
	"io"
)

// ImageProvider is a common - TODO
type ImageProvider interface {
	ImageMarkup() ImageProps
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

	// Defines text that can replace the image in the page.
	// Note: Will not output if IsSource is true, as alt  is not valid for source elements.
	Alt string

	// Determines if the image is a source element.
	IsSource bool

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
	MimeType ImageMimeType

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
	// For example: markup.Attributes{"id": "main-picture", "class": "responsive-picture"}
	Attributes Attributes
}

// Image - TODO
func Image(props ImageProps, opts ...ImageOptions) ImageProps {
	for _, opt := range opts {
		opt(&props)
	}
	return props
}

// Render renders a picture element to the provided writer.
func (i ImageProps) Render(ctx context.Context, w io.Writer) error {
	return templates.Render(ctx, w, "image.html", i)
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

// ImageSize specifies the size of the image, as defined in the static and
// dynamic renderers. With any luck, all sizes listed below will be
// rendered.
type ImageSize string

// ImageSize constants that are defined by sharp when resizing images.
const (
	ImageSizeThumbnail ImageSize = "thumbnail"
	ImageSizeMobile    ImageSize = "mobile"
	ImageSizeTablet    ImageSize = "tablet"
	ImageSizeDesktop   ImageSize = "desktop"
)

// ImageMimeType specifies a mimetype four a <source> element that is outputted
// on the type attribute.
type ImageMimeType string

// ImageMimeType constants that are defined at:
// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/Image_types
const (
	ImageMimeTypeAPNG ImageMimeType = "image/apng"
	ImageMimeTypeAVIF ImageMimeType = "image/avif"
	ImageMimeTypeGif  ImageMimeType = "image/gif"
	ImageMimeTypeJPG  ImageMimeType = "image/jpeg"
	ImageMimeTypePNG  ImageMimeType = "image/png"
	ImageMimeTypeSVG  ImageMimeType = "image/svg+xml"
	ImageMimeTypeWebP ImageMimeType = "image/webp"
)

// ImageOptions allows for optional settings to be applied to an <img> or <source>.
type ImageOptions func(p *ImageProps)

// ImageWithAlt attaches alternative text to the image.
func ImageWithAlt(alt string) ImageOptions {
	return func(p *ImageProps) {
		p.Alt = alt
	}
}

// ImageWithLazyLoading sets loading=lazy to the image.
func ImageWithLazyLoading() ImageOptions {
	return func(p *ImageProps) {
		p.Loading = LoadingLazy
	}
}

// ImageWithEagerLoading sets loading=eager to the image.
func ImageWithEagerLoading() ImageOptions {
	return func(p *ImageProps) {
		p.Loading = LoadingEager
	}
}

// ImageWithWidth sets the width of the image.
func ImageWithWidth(width int) ImageOptions {
	return func(p *ImageProps) {
		p.Width = &width
	}
}

// ImageWithHeight sets the height of the image.
func ImageWithHeight(height int) ImageOptions {
	return func(p *ImageProps) {
		p.Height = &height
	}
}

// ImageWithAttribute attaches a custom attribute to the image that
// will be rendered in the HTML, for example data-id="main-image".
func ImageWithAttribute(key, value string) ImageOptions {
	return func(p *ImageProps) {
		if p.Attributes == nil {
			p.Attributes = make(Attributes)
		}
		p.Attributes[key] = value
	}
}
