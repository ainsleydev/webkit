package static

import (
	"fmt"
	"github.com/ainsleydev/webkit/pkg/markup"
	"image"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Image string

const ImageDistPath = "dist"

// NewImage creates a new image type.
func NewImage(s string) Image {
	return Image(s)
}

// String implements the Stringer interface to cast the image path
// to a string.
func (i Image) String() string {
	return string(i)
}

// ImageMarkup implements the markup.ImageProvider interface and transforms the static image
// into a markup.ImageProps type ready for rendering an <img> to the DOM.
func (i Image) ImageMarkup() markup.ImageProps {
	return markup.ImageProps{
		URL:      i.String(),
		IsSource: false,
	}
}

// PictureMarkup implements the markup.PictureProvider interface and transforms the static image
// into a markup.PictureProps type ready for rendering a <picture> the DOM.
func (i Image) PictureMarkup() markup.PictureProps {
	m := markup.PictureProps{
		URL:     i.String(),
		Sources: make([]markup.ImageProps, 0, 10),
	}

	pwd, err := os.Getwd()
	if err != nil {
		slog.Error("Error getting current working directory: " + err.Error())
		return m
	}

	for _, path := range i.ImagePaths() {
		path = filepath.Join(pwd, path)
		fmt.Println(path)

		// Check if the image file exists
		if _, err := os.Stat(path); err == nil {
			m.Sources = append(m.Sources, path.Join(httpPath, imagePath))
		}
	}

	return markup.PictureProps{}
}

// RootPath returns the path relative to the application directory structure.
// I.e. transforms: /assets/images/hello.jpg to /dist/images.hello.jpg
func (i Image) RootPath() string {
	parts := strings.Split(strings.TrimPrefix(i.String(), "/"), "/")
	if len(parts) > 1 {
		parts[0] = ImageDistPath
		return "/" + strings.Join(parts, "/")
	}
	return i.String()
}

func (i Image) ImagePaths() []string {
	img := i.String()
	baseName := removeFileExtension(filepath.Base(img))
	dir := filepath.Dir(i.RootPath())

	var s []string

	for _, size := range imageSizes {
		for _, ext := range append(imageExtensions, filepath.Ext(img)) {
			s = append(s, filepath.Join(dir, baseName+"-"+size+ext))
		}
	}

	for _, ext := range imageExtensions {
		s = append(s, filepath.Join(dir, baseName+ext))
	}

	return s
}

// removeFileExtension removes the file extension from a given filename.
func removeFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}

type imageProperties struct {
	Width  int
	Height int
}

func getProperties(r io.Reader) imageProperties {
	img, _, err := image.Decode(r)
	if err != nil {
		slog.Error(err.Error())
	}

	bounds := img.Bounds()
	return imageProperties{
		Width:  bounds.Max.X,
		Height: bounds.Max.Y,
	}
}

var (
	imageSizes = []string{
		"thumbnail",
		"mobile",
		"tablet",
		"desktop",
	}
	imageExtensions = []string{
		".avif",
		".webp",
	}
)
