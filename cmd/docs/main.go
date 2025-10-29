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

	"github.com/ainsleydev/webkit/internal/printer"
)

const (
	// guidelinesURL is the endpoint for fetching ainsley.dev guidelines.
	guidelinesURL = "https://ainsley.dev/guidelines/index.json"

	// outputDir is the directory where generated files are written.
	outputDir = "internal/gen/docs"

	codeStyleTemplate = "CODE_STYLE.md"
	payloadTemplate   = "PAYLOAD.md"
	svelteKitTemplate = "SVELTEKIT.md"
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

func main() {
	var output string
	flag.StringVar(&output, "output", outputDir, "Output directory for generated files")
	flag.Parse()

	ctx := context.Background()
	fs := afero.NewOsFs()
	p := printer.New(os.Stdout)

	p.Info("Fetching guidelines from ainsley.dev...")
	p.LineBreak()

	if err := run(ctx, fs, output); err != nil {
		p.Error(err.Error())
		p.LineBreak()
		os.Exit(1)
	}

	p.Success("Documentation files generated successfully")
	p.LineBreak()
}

// run executes the main documentation generation logic.
func run(ctx context.Context, fs afero.Fs, output string) error {
	// Fetch guidelines from ainsley.dev
	guidelines, err := fetchGuidelines(ctx)
	if err != nil {
		return errors.Wrap(err, "fetching guidelines")
	}

	// Group guidelines by section
	grouped := groupBySection(guidelines)

	if err = generateCodeStyleFile(fs, output, grouped); err != nil {
		return errors.Wrap(err, string("generating "+codeStyleTemplate))
	}

	if err = generatePayloadFile(fs, output, grouped); err != nil {
		return errors.Wrap(err, string("generating "+payloadTemplate))
	}

	if err = generateSvelteKitFile(fs, output, grouped); err != nil {
		return errors.Wrap(err, string("generating "+svelteKitTemplate))
	}

	return nil
}

// fetchGuidelines retrieves guidelines from the ainsley.dev site.
// The site is built with Hugo so it outputs a nice index.json
// file that can be latched on to.
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
	if err = json.NewDecoder(resp.Body).Decode(&guidelines); err != nil {
		return nil, errors.Wrap(err, "decoding guidelines")
	}

	return guidelines, nil
}

// groupBySection organizes guidelines by their section.
func groupBySection(guidelines []Guideline) map[string][]Guideline {
	grouped := make(map[string][]Guideline)

	for _, g := range guidelines {
		grouped[g.Section] = append(grouped[g.Section], g)
	}

	return grouped
}

// generateCodeStyleFile creates CODE_STYLE.md from HTML, SCSS, Go, JS, and Git sections.
func generateCodeStyleFile(fs afero.Fs, output string, grouped map[string][]Guideline) error {
	var buf bytes.Buffer

	sections := []string{"HTML", "SCSS", "Go", "JS", "Git"}

	for _, section := range sections {
		guidelines, exists := grouped[section]
		if !exists {
			continue
		}

		buf.WriteString(fmt.Sprintf("## %s\n\n", section))

		for _, g := range guidelines {
			adjusted := adjustHeadingLevels(g.Markdown)
			buf.WriteString(adjusted)
			buf.WriteString("\n\n")
		}
	}

	return writeFile(fs, output, codeStyleTemplate, buf.Bytes())
}

func generatePayloadFile(fs afero.Fs, output string, grouped map[string][]Guideline) error {
	var buf bytes.Buffer

	guidelines, exists := grouped["Payload"]
	if !exists {
		return nil
	}

	buf.WriteString("## Payload\n\n")

	for _, g := range guidelines {
		adjusted := adjustHeadingLevels(g.Markdown)
		buf.WriteString(adjusted)
		buf.WriteString("\n\n")
	}

	return writeFile(fs, output, payloadTemplate, buf.Bytes())
}

func generateSvelteKitFile(fs afero.Fs, output string, grouped map[string][]Guideline) error {
	var buf bytes.Buffer

	guidelines, exists := grouped["SvelteKit"]
	if !exists {
		return nil
	}

	buf.WriteString("## SvelteKit\n\n")

	for _, g := range guidelines {
		adjusted := adjustHeadingLevels(g.Markdown)
		buf.WriteString(adjusted)
		buf.WriteString("\n\n")
	}

	return writeFile(fs, output, svelteKitTemplate, buf.Bytes())
}

// writeFile writes content to a file in the specified output directory.
func writeFile(fs afero.Fs, outputDir string, template string, content []byte) error {
	if err := fs.MkdirAll(outputDir, 0o755); err != nil {
		return errors.Wrap(err, "creating output directory")
	}

	path := filepath.Join(outputDir, template)
	if err := afero.WriteFile(fs, path, content, 0o644); err != nil {
		return errors.Wrap(err, "writing file")
	}

	return nil
}

// adjustHeadingLevels increases all heading levels by one (H2 -> H3, H3 -> H4, etc.).
func adjustHeadingLevels(markdown string) string {
	// Split into lines to handle code blocks properly
	lines := strings.Split(markdown, "\n")
	var result []string
	inCodeBlock := false

	for _, line := range lines {
		// Track code blocks to avoid modifying headings inside them
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			result = append(result, line)
			continue
		}

		// Only adjust headings outside code blocks
		if !inCodeBlock && strings.HasPrefix(strings.TrimSpace(line), "#") {
			// Match headings (##, ###, etc.) at the start of lines
			re := regexp.MustCompile(`^(\s*)(#{2,6})(\s+.*)$`)
			line = re.ReplaceAllString(line, "${1}#${2}${3}")
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
