package static

import (
	"fmt"
	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
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
	orig := getDistPath
	defer func() {
		getDistPath = orig
	}()
	getDistPath = func(path string) string {
		wd, err := os.Getwd()
		require.NoError(t, err)
		testdata := strings.Replace(path, "assets", "testdata", 1)
		return filepath.Join(wd, testdata)
	}

	input := Image("/assets/images/gopher.jpg")
	mm := input.PictureMarkup()
	b, err := json.MarshalIndent(mm, "", "\t")
	require.NoError(t, err)
	fmt.Println(string(b))
}

func TestImage_ImagePaths(t *testing.T) {
	input := Image("/assets/images/image.jpg")
	want := []string{
		"/assets/images/image-thumbnail.avif",
		"/assets/images/image-thumbnail.webp",
		"/assets/images/image-thumbnail.jpg",
		"/assets/images/image-mobile.avif",
		"/assets/images/image-mobile.webp",
		"/assets/images/image-mobile.jpg",
		"/assets/images/image-tablet.avif",
		"/assets/images/image-tablet.webp",
		"/assets/images/image-tablet.jpg",
		"/assets/images/image-desktop.avif",
		"/assets/images/image-desktop.webp",
		"/assets/images/image-desktop.jpg",
		"/assets/images/image.avif",
		"/assets/images/image.webp",
	}
	got := input.imageSources()
	assert.Equal(t, want, got)
}
