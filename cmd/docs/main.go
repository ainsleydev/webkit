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

	// outputPath is where the generated AGENTS.md file will be written.
	outputPath = "AGENTS.md"
)

func main() {
	fs := afero.NewOsFs()
	console := printer.New(os.Stdout)
	tracker := manifest.NewTracker()
	gen := scaffold.New(fs, tracker, console)

	opts := docs.GenerateAgentsOptions{
		FS:                fs,
		Generator:         gen,
		CustomContentPath: customContentPath,
		OutputPath:        outputPath,
		TemplateData:      nil,
		TrackingSource:    manifest.SourceProject(),
	}

	if err := docs.GenerateAgents(opts); err != nil {
		console.Error(fmt.Sprintf("Failed to generate AGENTS.md: %s", err))
		os.Exit(1)
	}

	console.Success(fmt.Sprintf("Generated %s successfully", outputPath))
}
