package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	docsutil "github.com/ainsleydev/webkit/internal/docs"
	"github.com/ainsleydev/webkit/internal/printer"
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

	fs := afero.NewOsFs()
	p := printer.New(os.Stdout)

	p.Info("Generating WebKit AGENTS.md...")
	p.LineBreak()

	if err := run(fs, output); err != nil {
		p.Error(err.Error())
		p.LineBreak()
		os.Exit(1)
	}

	p.Success("AGENTS.md generated successfully")
	p.LineBreak()
}

// run executes the main AGENTS.md generation logic.
func run(fs afero.Fs, output string) error {
	// Load WebKit-specific content from docs/AGENTS.md
	contentBytes, err := afero.ReadFile(fs, docsContentPath)
	if err != nil {
		return errors.Wrap(err, "reading docs/AGENTS.md")
	}
	content := string(contentBytes)

	// Load generated CODE_STYLE.md from internal/gen/docs
	codeStyle := docsutil.MustLoadGenFile(fs, docsutil.CodeStyleTemplate)

	// Load the AGENTS.md template
	tmpl := templates.MustLoadTemplate("AGENTS.md")

	// Prepare template data
	data := map[string]any{
		"Content":   content,
		"CodeStyle": codeStyle,
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return errors.Wrap(err, "executing template")
	}

	// Write to output file
	if err := afero.WriteFile(fs, output, buf.Bytes(), 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("writing %s", output))
	}

	return nil
}
