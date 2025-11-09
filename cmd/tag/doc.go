// Package main provides the tag command for managing Git version tags.
//
// This interactive CLI tool helps create and delete semantic version tags for releases.
// It handles version bumping (patch/minor/major), generates version files, creates commits,
// and pushes tags to trigger the GoReleaser pipeline.
//
// # Usage
//
// Run the interactive tag manager:
//
//	go run cmd/tag/main.go
//
// The tool provides a menu for:
//   - Creating new tags (patch, minor, or major version bumps)
//   - Deleting existing tags (local and remote)
//   - Viewing available tags
package main
