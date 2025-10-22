package tfembed

import (
	"embed"
)

//go:embed base modules providers
var Templates embed.FS
