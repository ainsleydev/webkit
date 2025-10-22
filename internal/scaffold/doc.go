// Package scaffold provides utilities for generating files in WebKit projects,
// including support for templates, scaffolding mode (non-overwriting), and manifest tracking.
//
// # Generation Modes
//
// The package supports two generation modes:
//   - Generate mode: Always writes the file, overwriting if it exists
//   - Scaffold mode: Only writes if the file doesn't exist, preserving user modifications
//
// # File Types
//
// The Generator interface supports multiple file formats:
//   - Raw bytes
//   - Go templates with custom functions
//   - JSON with automatic formatting
//   - YAML with proper serialization
//   - Copying from embedded filesystems
//
// # Manifest Tracking
//
// All generated files are automatically tracked in the manifest with metadata
// about their source, generator, and content hash. This enables drift detection
// and cleanup of orphaned files when the app definition changes.
package scaffold
