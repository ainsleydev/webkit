package fsext

import "embed"

// ReadFromEmbed reads a file from an embedded filesystem and
// returns its contents as a string.
//
// Returns an error if the file cannot be read.
func ReadFromEmbed(fs embed.FS, name string) (string, error) {
	file, err := fs.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

// MustReadFromEmbed reads a file from an embedded filesystem and
// returns its contents as a string.
//
// Panics if the file cannot be read.
func MustReadFromEmbed(fs embed.FS, name string) string {
	content, err := ReadFromEmbed(fs, name)
	if err != nil {
		panic(err)
	}
	return content
}
