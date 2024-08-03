package static

import (
	"fmt"
	"github.com/ainsleydev/webkit/pkg/markup"
	"github.com/gabriel-vasile/mimetype"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Image defines a static image to be rendered onto the DON, it can either
// be a <picture> element, or an <img> element.
type Image string

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
		Sources: make([]markup.ImageProps, 0, 12), // Approx no of images after conversion.
	}

	for _, path := range i.imagePaths() {
		source, err := imageSourceMarkup(path)
		if err != nil {
			slog.Error("Obtaining image: "+err.Error(), "path", path)
			continue
		}
		m.Sources = append(m.Sources, source)
	}

	return m
}

var getDistPath = func(path string) string {
	// Assume that where the executable is run, is where the dist folder is.
	wd, err := os.Executable()
	if err != nil {
		slog.Error("Error getting current working directory: " + err.Error())
		return ""
	}
	return filepath.Join(wd, AssetToBasePath(path))
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

type sourceImage struct {
	path   string
	width  int
	height int
}

func (i Image) imagePaths() []string {
	img := i.String()
	baseName := removeFileExtension(filepath.Base(img))
	dir := filepath.Dir(img)

	var s []string

	for _, size := range imageSizes {
		s = append(s, filepath.Join(dir, baseName+"-"+size+filepath.Ext(img)))

		for _, ext := range append(imageExtensions) {
			s = append(s, filepath.Join(dir, baseName+"-"+size+ext))
		}
	}

	for _, ext := range imageExtensions {
		s = append(s, filepath.Join(dir, baseName+ext))
	}

	return s
}

func imageSourceMarkup(path string) (markup.ImageProps, error) {
	def := markup.ImageProps{
		URL:      path,
		IsSource: true,
	}

	distPath := getDistPath(path)
	if _, err := os.Stat(distPath); err != nil {
		return def, err
	}

	b, err := os.Open(distPath)
	if err != nil {
		return def, err
	}
	defer b.Close() // Don't forget to close the file when we're done

	mime, err := mimetype.DetectReader(b)
	if err != nil {
		return markup.ImageProps{}, err
	}

	// Reset the file pointer to the beginning
	_, err = b.Seek(0, io.SeekStart)
	if err != nil {
		slog.Error("Resetting file pointer: " + distPath)
		return def, err
	}

	if mime.String() == "image/jpeg" {
		fmt.Println(getProperties(b))
	}

	//props, err := getProperties(b)
	//if err != nil {
	//	return def, err
	//}

	return markup.ImageProps{
		URL:      path,
		Alt:      "",
		IsSource: false,
		Media:    "",
		MimeType: markup.ImageMimeType(mime.String()),
		//Width:      &props.Width,
		//Height:     &props.Height,
		Attributes: markup.Attributes{},
	}, nil
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

func getProperties(r io.Reader) (imageProperties, error) {
	def := imageProperties{}
	img, _, err := image.Decode(r)
	if err != nil {
		return def, err
	}
	if img == nil {
		return def, err
	}
	bounds := img.Bounds()
	return imageProperties{
		Width:  bounds.Max.X,
		Height: bounds.Max.Y,
	}, nil
}
