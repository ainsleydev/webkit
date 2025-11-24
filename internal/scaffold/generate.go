package scaffold

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"

	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/pkg/enforce"
)

type (
	// Generator is used for scaffolding files to a WebKit project.
	Generator interface {
		Bytes(path string, data []byte, opts ...Option) error
		Copy(from, to string, opts ...Option) error
		CopyFromEmbed(efs embed.FS, from, to string, opts ...Option) error
		Template(path string, tpl *template.Template, data any, opts ...Option) error
		JSON(path string, content any, opts ...Option) error
		YAML(path string, content any, opts ...Option) error
		Code(path string, content string, opts ...Option) error
	}
	// FileGenerator handles file generation on a given filesystem.
	FileGenerator struct {
		Printer  *printer.Console
		fs       afero.Fs
		manifest *manifest.Tracker
	}
)

// New creates a new FileGenerator with the provided afero.Fs.
func New(fs afero.Fs, manifest *manifest.Tracker, printer *printer.Console) *FileGenerator {
	enforce.NotNil(fs, "file system is required")
	enforce.NotNil(manifest, "manifest definition is required")
	enforce.NotNil(printer, "printer is required")

	return &FileGenerator{
		Printer:  printer,
		fs:       fs,
		manifest: manifest,
	}
}

// WriteMode determines how files are written
type WriteMode int

const (
	// ModeGenerate always writes the file, overwriting if it exists
	ModeGenerate WriteMode = iota
	// ModeScaffold only writes if the file doesn't exist
	ModeScaffold
)

// Bytes writes bytes to the filesystem and ensure directories exist.
func (f FileGenerator) Bytes(path string, data []byte, opts ...Option) error {
	options := applyOptions(opts...)

	// Add to the manifest at to begin with, otherwise
	// scaffolded files won't be appended.
	if options.tracking.enabled {
		f.manifest.Add(manifest.FileEntry{
			Path:         path,
			Generator:    options.tracking.generator,
			Source:       options.tracking.source,
			ScaffoldMode: options.mode == ModeScaffold,
			Hash:         manifest.HashContent(data),
			GeneratedAt:  time.Now(),
		})
	}

	if f.shouldSkipScaffold(path, options.mode) {
		return nil
	}

	if !options.suppressNotice {
		notice := []byte(noticeForFile(path))
		data = append(notice, data...)
	}

	if err := fsext.EnsureDir(f.fs, path); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}

	exists := fsext.Exists(f.fs, path)
	if exists {
		// f.Printer.Println("Updated: " + path)
	}

	if err := afero.WriteFile(f.fs, path, data, os.ModePerm); err != nil {
		return fmt.Errorf("writing file %s: %w", path, err)
	}

	if !exists {
		// f.Printer.Println("Created: " + path)
	}

	return nil
}

// Copy simply copies a file to the destination using the scaffolder.
func (f FileGenerator) Copy(from, to string, opts ...Option) error {
	file, err := afero.ReadFile(f.fs, from)
	if err != nil {
		return errors.Wrap(err, "unable to copy file")
	}
	return f.Bytes(to, file, opts...)
}

// CopyFromEmbed copies a file from an embedded FS to the generator's FS.
func (f FileGenerator) CopyFromEmbed(efs embed.FS, from, to string, opts ...Option) error {
	file, err := efs.ReadFile(from)
	if err != nil {
		return errors.Wrap(err, "unable to copy embedded file")
	}
	return f.Bytes(to, file, opts...)
}

// Template writes a template file with the given mode.
func (f FileGenerator) Template(path string, tpl *template.Template, data any, opts ...Option) error {
	options := applyOptions(opts...)

	buf := &bytes.Buffer{}
	if !options.suppressNotice {
		buf.WriteString(noticeForFile(path))
	}

	if err := tpl.Execute(buf, data); err != nil {
		return fmt.Errorf("executing template %s: %w", tpl.Name(), err)
	}

	opts = append(opts, WithoutNotice())

	return f.Bytes(path, buf.Bytes(), opts...)
}

// JSON writes JSON content with the given mode.
func (f FileGenerator) JSON(path string, content any, opts ...Option) error {
	buf := &bytes.Buffer{}

	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "\t")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("encoding %s: %w", path, err)
	}

	opts = append(opts, WithoutNotice())

	return f.Bytes(path, buf.Bytes(), opts...)
}

// YAML writes YAML content with the given mode.
func (f FileGenerator) YAML(path string, content any, opts ...Option) error {
	buf := &bytes.Buffer{}
	buf.WriteString(noticeForFile(path))

	encoder := yaml.NewEncoder(buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("encoding %s: %w", path, err)
	}

	opts = append(opts, WithoutNotice())

	return f.Bytes(path, buf.Bytes(), opts...)
}

// Code writes Go code to a file with the given mode.
// The content should be valid Go code and will have the WebKit notice prepended.
func (f FileGenerator) Code(path string, content string, opts ...Option) error {
	return f.Bytes(path, []byte(content), opts...)
}

func (f FileGenerator) shouldSkipScaffold(path string, mode WriteMode) bool {
	if mode != ModeScaffold {
		return false
	}

	exists := fsext.Exists(f.fs, path)
	if !exists {
		return false
	}

	f.Printer.Println("• skipped scaffolding " + path + " - already exists")
	return true
}
