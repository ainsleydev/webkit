package gen

import (
	"embed"
	_ "embed"
)

//go:embed *
var Embed embed.FS
