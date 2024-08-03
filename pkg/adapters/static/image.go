package static

import (
	"errors"
	"github.com/ainsleydev/webkit/pkg/markup"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
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
	return markup.PictureProps{
		URL:     i.String(),
		Sources: i.imageSources(),
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
)

func (i Image) imageSources() []markup.ImageProps {
	img := i.String()
	//baseName := removeFileExtension(filepath.Base(img))

	var images []markup.ImageProps

	for _, size := range imageSizes {
		for _, ext := range append(imageExtensions, filepath.Ext(img)) {
			distDir := getDistPath(img)
			dir := path.Dir(distDir)
			base := removeFileExtension(filepath.Base(distDir))
			glob := filepath.Join(dir, base+"-"+size+"-*"+ext)

			matches, err := filepath.Glob(glob)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			if len(matches) == 0 {
				slog.Error("No matches found for image: " + glob)
				continue
			}

			//fmt.Println(matches[0], base)
			props, err := getImageProperties(filepath.Base(matches[0]))
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			images = append(images, markup.ImageProps{
				URL:        filepath.Join(filepath.Dir(img), filepath.Base(matches[0])),
				IsSource:   true,
				MimeType:   props.Mime,
				Width:      &props.Width,
				Height:     &props.Height,
				Attributes: nil,
			})
		}
	}

	//for _, ext := range imageExtensions {
	//	s = append(s, filepath.Join(filepath.Dir(img), baseName+ext))
	//}

	return images
}

// removeFileExtension removes the file extension from a given filename.
func removeFileExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	return strings.TrimSuffix(fileName, ext)
}

type imageProperties struct {
	Width  int
	Height int
	Mime   markup.ImageMimeType
}

var imageExtensionsToMime = map[string]markup.ImageMimeType{
	".avif": markup.ImageMimeTypeAVIF,
	".jpg":  markup.ImageMimeTypeJPG,
	".jpeg": markup.ImageMimeTypeJPG,
	".png":  markup.ImageMimeTypeAPNG,
	".webp": markup.ImageMimeTypeWebP,
}

var imageWidthHeightRegex = regexp.MustCompile(`.*-(\d+)x(\d+)\.[a-zA-Z0-9]+$`)

func getImageProperties(baseName string) (imageProperties, error) {
	matches := imageWidthHeightRegex.FindStringSubmatch(baseName)
	if matches == nil || len(matches) != 3 {
		return imageProperties{}, errors.New("No matches TODO")
	}
	mime, ok := imageExtensionsToMime[filepath.Ext(baseName)]
	if !ok {
		slog.Error("No mime found: " + baseName)
	}
	return imageProperties{
		Width:  mustConvertStringToInt(matches[1]),
		Height: mustConvertStringToInt(matches[2]),
		Mime:   mime,
	}, nil
}

func mustConvertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		slog.Error("Converting string " + s + " to int for image props")
		return 0
	}
	return i
}
