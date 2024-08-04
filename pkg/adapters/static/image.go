package static

import (
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ainsleydev/webkit/pkg/markup"
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
	sources, err := i.imageSources()
	if err != nil {
		slog.Error(err.Error(), "image", i)
		return markup.PictureProps{URL: i.String()}
	}
	return markup.PictureProps{
		URL:     i.String(),
		Sources: sources,
	}
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
	imageExtensionsToMime = map[string]markup.ImageMimeType{
		".avif": markup.ImageMimeTypeAVIF,
		".jpg":  markup.ImageMimeTypeJPG,
		".jpeg": markup.ImageMimeTypeJPG,
		".png":  markup.ImageMimeTypeAPNG,
		".webp": markup.ImageMimeTypeWebP,
	}
)

// imageSources obtains all of the <source> elements within a static directory.
func (i Image) imageSources() ([]markup.ImageProps, error) {
	img := i.String()

	var images []markup.ImageProps

	// Obtain all of the images that have been resized, for example:
	// gopher-desktop-1440-1920.avif
	for _, size := range imageSizes {
		for _, ext := range append(imageExtensions, filepath.Ext(img)) {
			source, err := i.getMatches(ext, size, true)
			if err != nil {
				return nil, err
			}

			props, err := getImageProperties(filepath.Base(source))
			if err != nil {
				return nil, err
			}

			images = append(images, markup.ImageProps{
				URL:      filepath.Join(filepath.Dir(img), filepath.Base(source)),
				IsSource: true,
				MimeType: props.Mime,
				Width:    &props.Width,
				Height:   &props.Height,
			})
		}
	}

	// Obtain all of the original images that have been not resized but adjusted
	// based on file extension, for example:
	// gopher.avif & gopher.web
	for _, ext := range imageExtensions {
		source, err := i.getMatches(ext, "", false)
		if err != nil {
			return nil, err
		}

		ext := filepath.Ext(source)
		mime, ok := imageExtensionsToMime[ext]
		if !ok {
			return nil, errors.New("no mime type found for extension: " + ext)
		}

		images = append(images, markup.ImageProps{
			URL:        filepath.Join(filepath.Dir(img), filepath.Base(source)),
			IsSource:   true,
			MimeType:   mime,
			Attributes: nil,
		})
	}

	return images, nil
}

// removeFileExtension removes the file extension from a given filename.
func removeFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}

// getMatches obtains all of the images associated with a file extension and
// optional size relative to the img root.
func (i Image) getMatches(extension, size string, useSize bool) (string, error) {
	img := i.String()
	distDir := getDistPath(img)
	dir := path.Dir(distDir)
	base := removeFileExtension(filepath.Base(distDir))
	glob := filepath.Join(dir, base+extension)

	// I.e. this is a source element with a name similar to: gopher-desktop-1440-1920.avif
	// Otherwise we're just looking for gopher.avif
	if useSize {
		glob = filepath.Join(dir, base+"-"+size+"-*"+extension)
	}

	matches, err := filepath.Glob(glob)
	if err != nil {
		return "", errors.New("obtaining glob for image: " + err.Error() + ", glob: " + glob)
	}

	if len(matches) == 0 {
		return "", errors.New("no matches found for glob: " + glob)
	}

	return matches[0], nil
}

type imageProperties struct {
	Width  int
	Height int
	Mime   markup.ImageMimeType
}

var imageWidthHeightRegex = regexp.MustCompile(`.*-(\d+)x(\d+)\.[a-zA-Z0-9]+$`)

func getImageProperties(baseName string) (imageProperties, error) {
	matches := imageWidthHeightRegex.FindStringSubmatch(baseName)
	if matches == nil || len(matches) != 3 {
		return imageProperties{}, errors.New("no regex matches found for: " + baseName)
	}

	ext := filepath.Ext(baseName)
	mime, ok := imageExtensionsToMime[ext]
	if !ok {
		return imageProperties{}, errors.New("no mime type found for extension: " + ext)
	}

	width, err := strconv.Atoi(matches[1])
	if err != nil {
		return imageProperties{}, errors.New("converting width to int: " + err.Error() + " for image: " + baseName)
	}

	height, err := strconv.Atoi(matches[2])
	if err != nil {
		return imageProperties{}, errors.New("converting height to int: " + err.Error() + " for image: " + baseName)
	}

	return imageProperties{
		Width:  width,
		Height: height,
		Mime:   mime,
	}, nil
}
