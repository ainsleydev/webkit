package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/gen"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

const (
	// docsContentPath is the path to the WebKit-specific content file.
	docsContentPath = "docs/AGENTS.md"

	// outputFile is the name of the generated root AGENTS.md file.
	outputFile = "AGENTS.md"
)

func main() {
	var output string
	flag.StringVar(&output, "output", outputFile, "Output file path for generated AGENTS.md")
	flag.Parse()

	p := printer.New(os.Stdout)

	p.Info("Generating WebKit AGENTS.md...")
	p.LineBreak()

	if err := run(afero.NewOsFs(), output); err != nil {
		p.Error(err.Error())
		p.LineBreak()
		os.Exit(1)
	}

	p.Success("AGENTS.md generated successfully")
}

func run(fs afero.Fs, output string) error {
	// Load WebKit-specific content from docs/AGENTS.md
	contentBytes, err := afero.ReadFile(fs, docsContentPath)
	if err != nil {
		return errors.Wrap(err, "reading docs/AGENTS.md")
	}

	data := map[string]any{
		"Content":   string(contentBytes),
		"CodeStyle": fsext.MustReadFromEmbed(gen.Embed, "docs/CODE_STYLE.md"),
	}

	var buf bytes.Buffer
	if err = templates.MustLoadTemplate("AGENTS.md").Execute(&buf, data); err != nil {
		return errors.Wrap(err, "executing template")
	}

	// Prepend WebKit notice as HTML comment
	notice := fmt.Sprintf("<!-- %s -->\n", scaffold.WebKitNotice)
	finalContent := append([]byte(notice), buf.Bytes()...)

	return afero.WriteFile(fs, output, finalContent, 0o644)
}
