package templates

import (
	"bytes"
	"context"
	"embed"
	"html/template"
	"io"

	"github.com/ainsleydev/webkit/pkg/tpl"
	"github.com/ainsleydev/webkit/pkg/util/stringutil"
)

//go:embed *.html
var templatesFS embed.FS

// templates are all the templates defined in the package.
var templates = template.Must(template.New("").Funcs(tpl.Funcs).ParseFS(templatesFS,
	"head.html",
	"image.html",
	"opengraph.html",
	"picture.html",
	"twitter.html",
))

func Render(_ context.Context, w io.Writer, name string, data any) error {
	buf := &bytes.Buffer{}
	if err := templates.ExecuteTemplate(buf, name, data); err != nil {
		return err
	}
	s := stringutil.RemoveDuplicateWhitespace(buf.String())
	_, err := w.Write([]byte(stringutil.FormatHTML(s)))
	return err
}
