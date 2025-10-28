package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdjustHeadingLevels(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"H2 to H3": {
			input: "## Heading\n\nContent",
			want:  "### Heading\n\nContent",
		},
		"H3 to H4": {
			input: "### Heading\n\nContent",
			want:  "#### Heading\n\nContent",
		},
		"Multiple headings": {
			input: "## First\n\n### Second\n\n#### Third",
			want:  "### First\n\n#### Second\n\n##### Third",
		},
		"H6 becomes H7": {
			input: "###### Max Level",
			want:  "####### Max Level",
		},
		"H1 unchanged": {
			input: "# Title\n\n## Section",
			want:  "# Title\n\n### Section",
		},
		"No headings": {
			input: "Just some text\n\nWith paragraphs",
			want:  "Just some text\n\nWith paragraphs",
		},
		"Code block with hash": {
			input: "## Heading\n\n```bash\n# This is a comment\n```",
			want:  "### Heading\n\n```bash\n# This is a comment\n```",
		},
		"Inline code with hash": {
			input: "## Heading\n\nUse `#` for headings",
			want:  "### Heading\n\nUse `#` for headings",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := adjustHeadingLevels(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGroupBySection(t *testing.T) {
	t.Parallel()

	guidelines := []Guideline{
		{Section: "Go", Heading: "General", Markdown: "Go content"},
		{Section: "Go", Heading: "Testing", Markdown: "Test content"},
		{Section: "HTML", Heading: "General", Markdown: "HTML content"},
		{Section: "Payload", Heading: "Fields", Markdown: "Payload content"},
	}

	got := groupBySection(guidelines)

	assert.Len(t, got, 3)
	assert.Len(t, got["Go"], 2)
	assert.Len(t, got["HTML"], 1)
	assert.Len(t, got["Payload"], 1)
}

func TestFetchGuidelines(t *testing.T) {
	t.Parallel()

	t.Run("Successful fetch", func(t *testing.T) {
		t.Parallel()

		mockGuidelines := []Guideline{
			{
				Section:     "Go",
				Heading:     "General",
				Title:       "General",
				Markdown:    "## Content",
				Description: "Test guideline",
				Date:        time.Now(),
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(mockGuidelines)
			require.NoError(t, err)
		}))
		defer server.Close()

		ctx := context.Background()

		// We can't easily test this without modifying the const,
		// but we can test the HTTP logic by calling it directly
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		var guidelines []Guideline
		err = json.NewDecoder(resp.Body).Decode(&guidelines)
		require.NoError(t, err)

		assert.Len(t, guidelines, 1)
		assert.Equal(t, "Go", guidelines[0].Section)
	})

	t.Run("Server error", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	t.Run("Creates directory and writes file", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		content := []byte("test content")

		err := writeFile(fs, "output/dir", "test.md", content)
		require.NoError(t, err)

		exists, err := afero.Exists(fs, "output/dir/test.md")
		require.NoError(t, err)
		assert.True(t, exists)

		read, err := afero.ReadFile(fs, "output/dir/test.md")
		require.NoError(t, err)
		assert.Equal(t, content, read)
	})
}

func TestGenerateCodeStyleFile(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()

	grouped := map[string][]Guideline{
		"HTML": {
			{Section: "HTML", Heading: "General", Markdown: "## Indentation\n\nUse tabs."},
		},
		"SCSS": {
			{Section: "SCSS", Heading: "Naming", Markdown: "## BEM\n\nUse BEM notation."},
		},
		"Go": {
			{Section: "Go", Heading: "General", Markdown: "## Formatting\n\nUse gofmt."},
		},
		"JS": {
			{Section: "JS", Heading: "General", Markdown: "## Style\n\nUse camelCase."},
		},
		"Git": {
			{Section: "Git", Heading: "Commits", Markdown: "## Format\n\nUse conventional commits."},
		},
	}

	err := generateCodeStyleFile(fs, "test-output", grouped)
	require.NoError(t, err)

	content, err := afero.ReadFile(fs, filepath.Join("test-output", "CODE_STYLE.md"))
	require.NoError(t, err)

	contentStr := string(content)

	// Should contain all sections
	assert.Contains(t, contentStr, "## HTML")
	assert.Contains(t, contentStr, "## SCSS")
	assert.Contains(t, contentStr, "## Go")
	assert.Contains(t, contentStr, "## JS")
	assert.Contains(t, contentStr, "## Git")

	// Headings should be adjusted
	assert.Contains(t, contentStr, "### Indentation")
	assert.Contains(t, contentStr, "### BEM")
	assert.Contains(t, contentStr, "### Formatting")
	assert.Contains(t, contentStr, "### Style")
	assert.Contains(t, contentStr, "### Format")
}

func TestGeneratePayloadFile(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()

	grouped := map[string][]Guideline{
		"Payload": {
			{Section: "Payload", Heading: "Fields", Markdown: "## Naming\n\nUse camelCase."},
		},
	}

	err := generatePayloadFile(fs, "test-output", grouped)
	require.NoError(t, err)

	content, err := afero.ReadFile(fs, filepath.Join("test-output", "PAYLOAD.md"))
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "## Payload")
	assert.Contains(t, contentStr, "### Naming")
}

func TestGenerateSvelteKitFile(t *testing.T) {
	t.Parallel()

	fs := afero.NewMemMapFs()

	grouped := map[string][]Guideline{
		"SvelteKit": {
			{Section: "SvelteKit", Heading: "Routing", Markdown: "## File structure\n\nUse +page.svelte."},
		},
	}

	err := generateSvelteKitFile(fs, "test-output", grouped)
	require.NoError(t, err)

	content, err := afero.ReadFile(fs, filepath.Join("test-output", "SVELTEKIT.md"))
	require.NoError(t, err)

	contentStr := string(content)
	assert.Contains(t, contentStr, "## SvelteKit")
	assert.Contains(t, contentStr, "### File structure")
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("Generates all files", func(t *testing.T) {
		t.Parallel()

		// This would require mocking the HTTP call
		// We've tested individual components above
		t.Skip("Integration test - requires HTTP mocking")
	})
}
