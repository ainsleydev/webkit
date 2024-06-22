package tpl

import (
	"embed"
	"html/template"
)

var (
	// FS embeds all the html files within the tpl package.
	//go:embed *
	FS embed.FS

	// Templates embeds all of the html files within the tpl package.
	Templates = template.Must(template.ParseFS(FS, "*.html"))
)
