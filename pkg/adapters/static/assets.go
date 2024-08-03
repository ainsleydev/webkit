package static

import "strings"

const AssetHTTPDir = "assets"

// AssetBaseDir is the base directory where assets are stored in the os.
const AssetBaseDir = "dist"

// AssetToBasePath returns the path relative to the application directory structure.
// I.e. transforms: /assets/images/hello.jpg to /dist/images.hello.jpg
func AssetToBasePath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) > 1 {
		parts[0] = AssetBaseDir
		return strings.Join(parts, "/")
	}
	return path
}
