package main

import (
	"fmt"
	"os"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/docs"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

const (
	// customContentPath defines where webkit-specific injected docs reside.
	// This can be easily changed if needed.
	customContentPath = "ai/docs"
)

func main() {
	fs := afero.NewOsFs()
	console := printer.New(os.Stdout)
	tracker := manifest.NewTracker()
	gen := scaffold.New(fs, tracker, console)

	opts := docs.GenerateDocumentOptions{
		FS:                fs,
		Generator:         gen,
		DocumentType:      docs.DocumentTypeAgents,
		CustomContentPath: customContentPath,
		Data:              nil, // No app.json manifest needed for webkit repo
		TrackingSource:    manifest.SourceProject(),
	}

	if err := docs.GenerateDocument(opts); err != nil {
		console.Error(fmt.Sprintf("Failed to generate %s: %s", opts.DocumentType, err))
		os.Exit(1)
	}

	console.Success(fmt.Sprintf("Generated %s successfully", opts.DocumentType))
}
