package static

import (
	"github.com/ainsleydev/webkit/pkg/markup"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Image string

const ImageDistPath = "/dist"

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
	return markup.PictureProps{}
}

// RootPath returns the path relative to the application directory structure.
// I.e. transforms: /assets/images/hello.jpg to /dist/images.hello.jpg
func (i Image) RootPath() string {
	parts := strings.Split(i.String(), "/")
	if len(parts) > 1 {
		parts[0] = ImageDistPath
		return "/" + strings.Join(parts, "/")
	}
	return i.String()
}

// removeFileExtension removes the file extension from a given filename.
func removeFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}

var (
	sizes = []string{
		"thumbnail",
		"mobile",
		"tablet",
		"",
	}
)

// getImages returns a list of image files associated with the original file in the same directory.
func getImages(original string) []string {
	pwd, err := os.Getwd()
	if err != nil {
		slog.Error("Error getting current working directory: " + err.Error())
		return nil
	}

	// Adjust the base directory path for checking and appending images
	baseName := removeFileExtension(filepath.Base(original))
	httpPath := filepath.Dir(original)
	osPath := strings.ReplaceAll(httpPath, "/assets/", "/dist/")

	extensions := []string{".avif", ".webp", filepath.Ext(original)}

	var images []string

	for _, size := range sizes {
		for _, ext := range extensions {
			var imagePath string
			if size == "" {
				imagePath = filepath.Join(baseName + ext)
			} else {
				imagePath = filepath.Join(baseName + "-" + size + ext)
			}

			// Check if the image file exists
			if _, err := os.Stat(filepath.Join(pwd, osPath, imagePath)); err == nil {
				images = append(images, path.Join(httpPath, imagePath))
			}
		}
	}

	return images
}
