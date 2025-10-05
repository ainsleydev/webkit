package cgtools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

// Generator handles file generation on a given filesystem
type Generator struct {
	fs afero.Fs
}

// NewGenerator creates a new Generator with the provided afero.Fs
func NewGenerator(fs afero.Fs) *Generator {
	return &Generator{fs: fs}
}

// WriteFile ensures directories exist and writes the data to the path
func (g Generator) WriteFile(path string, data []byte) error {
	if err := g.fs.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	if err := afero.WriteFile(g.fs, path, data, 0o644); err != nil {
		return fmt.Errorf("writing file %s: %w", path, err)
	}

	fmt.Println("Created:", path)
	return nil
}

// ScaffoldTemplate writes the template only if the file does not already exist
func (g Generator) ScaffoldTemplate(path string, tpl *template.Template, data any) error {
	exists, _ := afero.Exists(g.fs, path)
	if exists {
		fmt.Println("• skipped scaffolding", path, "- already exists")
		return nil
	}

	return g.writeFileWithTemplate(path, tpl, data)
}

// GenerateTemplate writes a template file, overwriting if it already exists
func (g Generator) GenerateTemplate(path string, tpl *template.Template, data any) error {
	exists, _ := afero.Exists(g.fs, path)
	if exists {
		fmt.Println("• regenerating", path, "- already exists")
	}

	return g.writeFileWithTemplate(path, tpl, data)
}

// GenerateJSON marshals content to JSON and prepends the Webkit
// notice if requested.
func (g Generator) GenerateJSON(path string, content any) error {
	buf := &bytes.Buffer{}

	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("marshalling %s: %w", path, err)
	}

	return g.WriteFile(path, buf.Bytes())
}

// GenerateYAML marshals content to YAML and prepends the Webkit
// notice if requested.
func (g Generator) GenerateYAML(path string, content any) error {
	buf := &bytes.Buffer{}

	buf.WriteString(noticeForFile(path))

	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)

	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("marshalling %s: %w", path, err)
	}

	return g.WriteFile(path, buf.Bytes())
}

const command = "webkit"

// writeFileWithTemplate executes the template and writes the result
func (g Generator) writeFileWithTemplate(path string, tpl *template.Template, data any) error {
	buf := &bytes.Buffer{}

	// Inject notice automatically
	buf.WriteString(noticeForFile(path))

	if err := tpl.Execute(buf, data); err != nil {
		return fmt.Errorf("executing template %s: %w", tpl.Name(), err)
	}

	return g.WriteFile(path, buf.Bytes())
}
