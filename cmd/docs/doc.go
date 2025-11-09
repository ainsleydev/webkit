// Package main provides the docs command for syncing documentation from ainsley.dev.
//
// This command fetches coding guidelines from ainsley.dev and generates markdown files
// for different sections (Code Style, Payload, SvelteKit). The guidelines are fetched
// from a JSON endpoint and formatted into categorized documentation files.
//
// # Usage
//
// Generate documentation files:
//
//	go run cmd/docs/main.go
//
// Or use the pnpm script:
//
//	pnpm docs:gen
//
// Generated files are written to internal/gen/docs/ directory.
package main
