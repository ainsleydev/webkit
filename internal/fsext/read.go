package fsext

import "embed"

func MustReadFromEmbed(fs embed.FS, name string) string {
	file, err := fs.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(file)
}
