package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

func main() {
	versionFlag := flag.String("version", "", "Version string (e.g., v1.0.0)")

	flag.Parse()

	// Use defaults if not provided
	version := getOrDefault(*versionFlag, "v0.0.1-dev")

	// Create filesystem and generator
	fs := afero.NewOsFs()
	console := printer.New(os.Stdout)
	tracker := manifest.NewTracker()
	gen := scaffold.New(fs, tracker, console)

	// Generate simplified version file content
	content := fmt.Sprintf(`package version

const Version = %q
`, version)

	// Write the generated file
	if err := gen.Code("internal/version/version.go", content); err != nil {
		console.Error(fmt.Sprintf("Failed to generate version file: %s", err))
		os.Exit(1)
	}

	console.Success(fmt.Sprintf("Generated version.go with version: %s", version))
}

func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
