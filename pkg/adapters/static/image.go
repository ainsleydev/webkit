package static

import (
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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
