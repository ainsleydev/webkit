// Package config provides utilities for reading and writing files
// in the WebKit configuration directory (~/.config/webkit).
//
// This directory stores user-specific configuration such as:
//   - Age encryption keys for SOPS
//   - CLI preferences and settings
//   - Cached data for improved performance
//
// All functions in this package automatically create the configuration
// directory if it doesn't exist, ensuring operations never fail due to
// missing directories.
package config
