package markup

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ainsleydev/webkit/pkg/markup/internal/templates"
)

// PictureProvider is a common - TODO
type PictureProvider interface {
	PictureMarkup() PictureProps
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

	// Determines if loading=lazy should be added to the image.
	Loading LoadingAttribute

	// The intrinsic width of the image in pixels , for example (300).
	// Must be an integer without a unit (optional).
	Width *int

	// The intrinsic height of the image, in pixels, for example (300).
	// Must be an integer without a unit (optional).
	Height *int

	// HideMediaSourcesWithSizeAttrs indicates if only next-gen image formats (AVIF & WebP)
	// should be used. When true, this will hide any <source> elements with size attributes,
	// effectively excluding them from the rendering process.
	HideMediaSizes bool

	// Attributes specifies additional attributes for the picture element as key-value pairs.
	// For example: markup.Attributes{"data-attribute-size": "large"}
	Attributes Attributes
}

// Picture returns picture properties - TODO
func Picture(provider PictureProvider, opts ...PictureOptions) PictureProps {
	props := provider.PictureMarkup()
	for _, opt := range opts {
		opt(&props)
	}

	// Add media query for sources with width attributes.
	for idx := range props.Sources {
		src := &props.Sources[idx]
		if src.IsSource && src.Width != nil {
			src.Media = fmt.Sprintf("(max-width: %vpx)", *src.Width+50)
		}
	}

	// If HideMediaSizes is true, remove sources with size attributes.
	if props.HideMediaSizes {
		var i int
		for _, src := range props.Sources {
			if src.Width == nil {
				props.Sources[i] = src
				i++
			}
		}
		props.Sources = props.Sources[:i]
	}

	return props
}

// Render renders a <picture> element to the provided writer.
func (p PictureProps) Render(ctx context.Context, w io.Writer) error {
	return templates.Render(ctx, w, "picture.html", p)
}

// Image transforms the PictureProps into an ImageProps type.
//
// This is useful when you want to render a single image element, instead
// of the entire picture.
func (p PictureProps) Image() ImageProps {
	return ImageProps{
		URL:        p.URL,
		Alt:        p.Alt,
		IsSource:   false,
		Media:      "", // Default image should not output a media query.
		MimeType:   "",
		Loading:    p.Loading,
		Width:      p.Width,
		Height:     p.Height,
		Attributes: p.Attributes,
	}
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

// PictureWithClasses appends any CSS classes to the picture.
func PictureWithClasses(classes ...string) PictureOptions {
	return func(p *PictureProps) {
		for _, v := range classes {
			p.Classes = append(p.Classes, v+" ")
		}
	}
}

// PictureWithHiddenMediaSources modifies the picture so sizes where sources with size
// attributes are hidden.
func PictureWithHiddenMediaSources() PictureOptions {
	return func(p *PictureProps) {
		p.HideMediaSizes = true
	}
}

// PictureWithAttribute attaches a custom attribute to the image that
// will be rendered in the HTML, for example data-id="main-image".
func PictureWithAttribute(key, value string) PictureOptions {
	return func(p *PictureProps) {
		if p.Attributes == nil {
			p.Attributes = make(Attributes)
		}
		p.Attributes[key] = value
	}
}

// PictureWithSize filters the sources to only include those that
// contain any of the specified size strings in their name.
//
// When multiple sizes are provided:
// - The last matching size will be used as the base source for URL, Width, and Height.
// - Only the last size applies filtering (excluding exact matches).
// - Earlier sizes include all matching sources.
func PictureWithSize(sizes ...string) PictureOptions {
	return func(p *PictureProps) {
		// If no sizes specified or no sources, return early.
		if len(sizes) == 0 || len(p.Sources) == 0 {
			return
		}

		// Find the base source that exactly matches the last matching size.
		var baseSource *ImageProps
		for i := len(sizes) - 1; i >= 0; i-- {
			size := sizes[i]
			for _, v := range p.Sources {
				if v.Name == size {
					source := v
					baseSource = &source
					goto found
				}
			}
		}
	found:
		// If we didn't find any base source, return early
		if baseSource == nil {
			return
		}

		// Create filtered sources slice
		filteredSources := make([]ImageProps, 0)
		lastSize := sizes[len(sizes)-1]

		// Add sources based on position in sizes list
		for _, source := range p.Sources {
			for _, size := range sizes {
				if strings.Contains(source.Name, size) {
					// For the last size, exclude exact matches
					if size == lastSize && source.Name == size {
						continue
					}
					// For all other sizes, include all matches
					filteredSources = append(filteredSources, source)
					break // Break inner loop once we've added this source
				}
			}
		}

		// Update the props
		p.URL = baseSource.URL
		p.Width = baseSource.Width
		p.Height = baseSource.Height
		p.Sources = filteredSources
	}
}
