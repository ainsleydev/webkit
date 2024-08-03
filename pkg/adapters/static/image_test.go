package static

import (
	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImage_ImageMarkup(t *testing.T) {
	input := NewImage("/assets/images/cat.jpg")
	got := input.ImageMarkup()
	want := markup.ImageProps{
		URL: "/assets/images/cat.jpg",
	}
	assert.Equal(t, want, got)
}

func TestImage_PictureMarkup(t *testing.T) {
	t.Skip()
}

func TestImage_RootPath(t *testing.T) {
	input := Image("/assets/images/image.jpg")
	input.PictureMarkup()
}

func TestImage_ImagePaths(t *testing.T) {
	input := Image("/assets/images/image.jpg")
	want := []string{
		"/dist/images/image-thumbnail.avif",
		"/dist/images/image-thumbnail.webp",
		"/dist/images/image-thumbnail.jpg",
		"/dist/images/image-mobile.avif",
		"/dist/images/image-mobile.webp",
		"/dist/images/image-mobile.jpg",
		"/dist/images/image-tablet.avif",
		"/dist/images/image-tablet.webp",
		"/dist/images/image-tablet.jpg",
		"/dist/images/image-desktop.avif",
		"/dist/images/image-desktop.webp",
		"/dist/images/image-desktop.jpg",
		"/dist/images/image.avif",
		"/dist/images/image.webp",
	}
	got := input.ImagePaths()
	assert.Equal(t, want, got)
}
