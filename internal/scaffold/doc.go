// Package scaffold handles file generation and scaffolding for WebKit projects.
//
// It provides a Generator interface for creating files from templates, JSON, YAML,
// and raw bytes with support for scaffold mode (don't overwrite) and generate mode
// (always overwrite). All generated files are tracked in the manifest for drift detection.
package scaffold
