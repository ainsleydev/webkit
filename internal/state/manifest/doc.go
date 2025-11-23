// Package manifest tracks WebKit-generated files for drift detection and cleanup.
//
// It maintains a manifest of all generated files with their source, generator,
// and content hashes. This enables detecting manual changes (drift) and removing
// orphaned files when configuration changes.
package manifest
