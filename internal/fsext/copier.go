package fsext

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyAllEmbed copies everything from the root of the embedded
// FS to destDir.
func CopyAllEmbed(efs embed.FS, destDir string) error {
	return CopyFromEmbed(efs, ".", destDir)
}

// CopyFromEmbed recursively copies all files from an embed.FS
// directory to a destination.
func CopyFromEmbed(fsx fs.FS, srcDir, destDir string) error {
	return fs.WalkDir(fsx, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, rel)

		if d.IsDir() {
			return os.MkdirAll(target, os.ModePerm)
		}

		data, err := fs.ReadFile(fsx, path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, 0o644)
	})
}
