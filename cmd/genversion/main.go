package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/version"
)

func main() {
	var (
		versionFlag = flag.String("version", "", "Version string (e.g., v1.0.0)")
		commitFlag  = flag.String("commit", "", "Git commit hash")
		dateFlag    = flag.String("date", "", "Build date")
		builtByFlag = flag.String("builtby", "local", "Built by (e.g., goreleaser, local)")
	)

	flag.Parse()

	// Use defaults if not provided
	versionInfo := version.VersionInfo{
		Version: getOrDefault(*versionFlag, "v0.0.1-dev"),
		Commit:  getOrDefault(*commitFlag, "none"),
		Date:    getOrDefault(*dateFlag, time.Now().Format(time.RFC3339)),
		BuiltBy: getOrDefault(*builtByFlag, "local"),
	}

	// Create filesystem and generator
	fs := afero.NewOsFs()
	console := printer.New(os.Stdout)
	tracker := manifest.NewTracker()
	gen := scaffold.New(fs, tracker, console)

	// Generate version file
	if err := version.GenerateVersionFile(fs, gen, versionInfo); err != nil {
		console.Error(fmt.Sprintf("Failed to generate version file: %s", err))
		os.Exit(1)
	}

	console.Success(fmt.Sprintf("Generated version.gen.go with version: %s", versionInfo.Version))
}

func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
