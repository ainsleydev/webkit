// Package enforce provides runtime assertion utilities that exit on failure.
//
// Unlike testing assertions, these functions are designed for application startup
// and configuration validation where failures should terminate the program immediately.
// Each function prints a descriptive error and exits with status 1 when conditions aren't met.
package enforce
