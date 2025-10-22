// Package appdef provides types and functions for parsing and working with
// the WebKit application definition file (app.json).
//
// The app definition is the central configuration file for WebKit projects,
// declaring what applications should be deployed, what infrastructure resources
// they need, and how they should be configured across different environments.
//
// # Core Concepts
//
// The Definition type represents the complete app.json structure, which includes:
//   - Project metadata (name, repository information)
//   - Shared configuration (environment variables common to all apps)
//   - Resources (databases, storage buckets, and other infrastructure)
//   - Apps (individual deployable applications with their own settings)
//
// # Application Types
//
// WebKit supports multiple application types, each with sensible defaults:
//   - golang: Go applications with Go-specific build commands
//   - svelte-kit: SvelteKit applications with npm/pnpm workflows
//   - payload: Payload CMS applications with npm/pnpm workflows
//
// # Environment Variables
//
// Environment variables can be defined at multiple levels with different sources:
//   - Shared environment: Variables available to all applications
//   - App-specific environment: Variables scoped to individual apps
//   - Source types: Static values, SOPS-encrypted secrets, or Terraform outputs
//
// The package handles merging these variable sources with proper precedence,
// where app-specific values override shared values.
//
// # Default Values
//
// The package automatically applies sensible defaults for commands, build settings,
// and resource configurations, allowing users to provide minimal configuration
// while still having full control when needed.
package appdef
