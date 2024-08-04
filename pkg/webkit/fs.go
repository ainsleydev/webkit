package webkit

import (
	"net/http"
	"path/filepath"
)

// TODO: Unit Test

// Static registers a new route with path prefix to serve static files from
// the provided root directory.
func (k *Kit) Static(pattern string, staticDir string, plugs ...Plug) {
	fs := http.FileServer(noDirListingFileSystem{http.Dir(staticDir)})
	handler := http.StripPrefix(pattern, fs)
	k.Get(pattern+"*", WrapHandler(handler), plugs...)
}

// noDirListingFileSystem is a custom file system wrapper that prevents
// directory listings.
type noDirListingFileSystem struct {
	fs http.FileSystem
}

// Open implements http/fs/FileSystem.Open to open the named file or dir.
// If the requested path is a directory, it checks for an index.html file and serves it if present.
// If no index.html is found, it returns an error to prevent directory listing.
func (nfs noDirListingFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, _ := f.Stat()
	if !s.IsDir() {
		return f, nil
	}

	index := filepath.Join(path, "index.html")
	if _, err := nfs.fs.Open(index); err != nil {
		closeErr := f.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	return f, nil
}
