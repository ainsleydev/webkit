package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/manifest"
)

const (
	// guidelinesURL is the endpoint for fetching ainsley.dev guidelines
	guidelinesURL = "https://ainsley.dev/guidelines/index.json"

	// agentsFilename is the name of the generated file
	agentsFilename = "AGENTS.md"

	// docsDir is the directory for WebKit-specific template content
	docsDir = "docs"
)

// Guideline represents a single guideline entry from ainsley.dev.
type Guideline struct {
	Date         time.Time `json:"date"`
	Description  string    `json:"description"`
	Draft        bool      `json:"draft"`
	Heading      string    `json:"heading"`
	LastMod      time.Time `json:"lastmod"`
	Markdown     string    `json:"markdown"`
	Permalink    string    `json:"permalink"`
	PlainContent string    `json:"plainContent"`
	PublishDate  time.Time `json:"publishdate"`
	Section      string    `json:"section"`
	Subsection   string    `json:"subsection"`
	Summary      string    `json:"summary"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Weight       int       `json:"weight"`
}

// GuidelineSet represents a collection of guidelines grouped by section.
type GuidelineSet struct {
	All              []Guideline
	BySectionHeading map[string][]Guideline
}

func main() {
	var (
		manifestPath string
	)

	flag.StringVar(&manifestPath, "manifest", "app.json", "Path to manifest file")
	flag.Parse()

	ctx := context.Background()
	fs := afero.NewOsFs()

	if err := run(ctx, fs, manifestPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main documentation generation logic.
func run(ctx context.Context, fs afero.Fs, manifestPath string) error {
	// Fetch guidelines from ainsley.dev
	guidelines, err := fetchGuidelines(ctx)
	if err != nil {
		return errors.Wrap(err, "fetching guidelines")
	}

	// Group guidelines by section
	guidelineSet := groupGuidelines(guidelines)

	// Load manifest to determine app types
	def, err := loadManifest(fs, manifestPath)
	if err != nil {
		return errors.Wrap(err, "loading manifest")
	}

	// Generate root AGENTS.md (excluding Payload and SvelteKit)
	if err := generateRootAgents(fs, guidelineSet); err != nil {
		return errors.Wrap(err, "generating root AGENTS.md")
	}

	// Generate app-specific AGENTS.md files
	if def != nil {
		if err := generateAppSpecificAgents(fs, guidelineSet, def); err != nil {
			return errors.Wrap(err, "generating app-specific AGENTS.md")
		}
	}

	return nil
}

// fetchGuidelines retrieves guidelines from ainsley.dev.
func fetchGuidelines(ctx context.Context) ([]Guideline, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, guidelinesURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fetching guidelines")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	var guidelines []Guideline
	if err := json.NewDecoder(resp.Body).Decode(&guidelines); err != nil {
		return nil, errors.Wrap(err, "decoding guidelines")
	}

	return guidelines, nil
}

// groupGuidelines organizes guidelines by section and heading.
func groupGuidelines(guidelines []Guideline) GuidelineSet {
	set := GuidelineSet{
		All:              guidelines,
		BySectionHeading: make(map[string][]Guideline),
	}

	for _, g := range guidelines {
		key := fmt.Sprintf("%s:%s", g.Section, g.Heading)
		set.BySectionHeading[key] = append(set.BySectionHeading[key], g)
	}

	return set
}

// loadManifest loads the app.json manifest file if it exists.
func loadManifest(fs afero.Fs, path string) (*appdef.Definition, error) {
	exists, err := afero.Exists(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "checking manifest existence")
	}

	// If no manifest exists (e.g., in WebKit repo itself), return nil
	if !exists {
		return nil, nil
	}

	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, errors.Wrap(err, "reading manifest")
	}

	var def appdef.Definition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, errors.Wrap(err, "unmarshaling manifest")
	}

	return &def, nil
}

// generateRootAgents creates the root AGENTS.md file with general guidelines.
func generateRootAgents(fs afero.Fs, guidelines GuidelineSet) error {
	// Sections to include in root AGENTS.md (excluding Payload and SvelteKit)
	includeSections := map[string]bool{
		"HTML": true,
		"SCSS": true,
		"Go":   true,
		"JS":   true,
		"Git":  true,
	}

	// Load custom WebKit-specific content
	customContent, err := loadCustomContent(fs)
	if err != nil {
		return errors.Wrap(err, "loading custom content")
	}

	var buf bytes.Buffer

	// Write header
	buf.WriteString("# Agent Guidelines\n\n")
	buf.WriteString("This document provides guidelines for AI agents working on the WebKit codebase.\n\n")
	buf.WriteString("## Note For Humans\n\n")
	buf.WriteString("This is a living document that will improve as more people/agents use it over time. Every effort has\n")
	buf.WriteString("been made to keep the guidance in here as generic and reusable as possible. Please keep this in mind\n")
	buf.WriteString("with any future edits.\n\n")

	// Add note about updating docs
	buf.WriteString("## Updating Documentation\n\n")
	buf.WriteString("If you need to update developer guidelines, clone and edit the [ainsley.dev/website](https://github.com/ainsleydev/website) repository. ")
	buf.WriteString("These guidelines are automatically synced from there.\n\n")

	// Add custom WebKit-specific content if it exists
	if customContent != "" {
		buf.WriteString(customContent)
		buf.WriteString("\n\n")
	}

	// Add guidelines sections
	for _, guideline := range guidelines.All {
		if !includeSections[guideline.Section] {
			continue
		}

		// Adjust heading levels (H2 -> H3, H3 -> H4, etc.)
		adjustedMarkdown := adjustHeadingLevels(guideline.Markdown)

		// Add section header if this is the first entry for this section
		sectionKey := fmt.Sprintf("%s:%s", guideline.Section, guideline.Heading)
		if len(guidelines.BySectionHeading[sectionKey]) > 0 {
			buf.WriteString(fmt.Sprintf("## %s\n\n", guideline.Section))
			buf.WriteString(adjustedMarkdown)
			buf.WriteString("\n\n")
		}
	}

	content := buf.String()

	// Write to root AGENTS.md
	if err := afero.WriteFile(fs, agentsFilename, []byte(content), 0644); err != nil {
		return errors.Wrap(err, "writing AGENTS.md")
	}

	return nil
}

// generateAppSpecificAgents creates AGENTS.md files in app subdirectories.
func generateAppSpecificAgents(fs afero.Fs, guidelines GuidelineSet, def *appdef.Definition) error {
	for _, app := range def.Apps {
		var sectionToInclude string

		switch app.Type {
		case appdef.AppTypePayload:
			sectionToInclude = "Payload"
		case appdef.AppTypeSvelteKit:
			sectionToInclude = "SvelteKit"
		default:
			// Skip other app types
			continue
		}

		// Generate content for this app type
		var buf bytes.Buffer

		buf.WriteString(fmt.Sprintf("# %s Guidelines\n\n", sectionToInclude))
		buf.WriteString(fmt.Sprintf("This document provides %s-specific guidelines for AI agents.\n\n", sectionToInclude))

		// Add relevant guidelines
		for _, guideline := range guidelines.All {
			if guideline.Section != sectionToInclude {
				continue
			}

			// Adjust heading levels
			adjustedMarkdown := adjustHeadingLevels(guideline.Markdown)
			buf.WriteString(adjustedMarkdown)
			buf.WriteString("\n\n")
		}

		// Write to app-specific AGENTS.md
		appAgentsPath := filepath.Join(app.Path, agentsFilename)
		if err := afero.WriteFile(fs, appAgentsPath, buf.Bytes(), 0644); err != nil {
			return errors.Wrap(err, fmt.Sprintf("writing %s", appAgentsPath))
		}
	}

	return nil
}

// loadCustomContent loads WebKit-specific content from docs/ directory.
func loadCustomContent(fs afero.Fs) (string, error) {
	// Try docs/AGENTS.md first
	agentsPath := filepath.Join(docsDir, agentsFilename)
	if exists, _ := afero.Exists(fs, agentsPath); exists {
		content, err := afero.ReadFile(fs, agentsPath)
		if err != nil {
			return "", errors.Wrap(err, "reading custom content")
		}
		return string(content), nil
	}

	return "", nil
}

// adjustHeadingLevels decreases all heading levels by one (H2 -> H3, etc.).
func adjustHeadingLevels(markdown string) string {
	// Match headings (##, ###, etc.) at the start of lines
	re := regexp.MustCompile(`(?m)^(#{2,6})\s`)

	return re.ReplaceAllStringFunc(markdown, func(match string) string {
		// Add one more # to decrease the level
		return "#" + match
	})
}
