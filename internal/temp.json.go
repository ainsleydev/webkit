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

// WriteMode determines how files are written
type WriteMode int

const (
	// ModeGenerate always writes the file, overwriting if it exists
	ModeGenerate WriteMode = iota
	// ModeScaffold only writes if the file doesn't exist
	ModeScaffold
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

// Template writes a template file with the given mode
func (g Generator) Template(path string, tpl *template.Template, data any, mode WriteMode) error {
	if mode == ModeScaffold {
		exists, _ := afero.Exists(g.fs, path)
		if exists {
			fmt.Println("• skipped scaffolding", path, "- already exists")
			return nil
		}
	}

	buf := &bytes.Buffer{}
	buf.WriteString(noticeForFile(path))

	if err := tpl.Execute(buf, data); err != nil {
		return fmt.Errorf("executing template %s: %w", tpl.Name(), err)
	}

	return g.WriteFile(path, buf.Bytes())
}

// JSON writes JSON content with the given mode
func (g Generator) JSON(path string, content any, mode WriteMode) error {
	if mode == ModeScaffold {
		exists, _ := afero.Exists(g.fs, path)
		if exists {
			fmt.Println("• skipped scaffolding", path, "- already exists")
			return nil
		}
	}

	buf := &bytes.Buffer{}

	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("marshalling %s: %w", path, err)
	}

	return g.WriteFile(path, buf.Bytes())
}

// YAML writes YAML content with the given mode
func (g Generator) YAML(path string, content any, mode WriteMode, addNotice bool) error {
	if mode == ModeScaffold {
		exists, _ := afero.Exists(g.fs, path)
		if exists {
			fmt.Println("• skipped scaffolding", path, "- already exists")
			return nil
		}
	}

	buf := &bytes.Buffer{}

	if addNotice {
		buf.WriteString(noticeForFile(path))
	}

	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("marshalling %s: %w", path, err)
	}

	return g.WriteFile(path, buf.Bytes())
}
