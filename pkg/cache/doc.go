// Package cache provides a unified interface for caching with multiple backend implementations.
//
// It supports various cache stores including in-memory, file-based, Redis, and OS-specific caches.
// The Store interface allows switching between backends without changing application code.
//
// All cache values are automatically marshaled/unmarshaled, and support tagging for grouped
// invalidation and configurable expiration times.
package cache
