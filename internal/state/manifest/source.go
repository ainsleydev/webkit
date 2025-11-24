package manifest

import "fmt"

// SourceProject returns the identifier string for the project-level
// manifest source. Used for base files with nothing attached.
func SourceProject() string {
	return "project"
}

// SourceApp returns a namespaced identifier for application sources.
// Format: "app:<name>"
func SourceApp(name string) string {
	return fmt.Sprintf("app:%s", name)
}

// SourceResource returns a namespaced identifier for resource sources.
// Format: "resource:<name>".
func SourceResource(name string) string {
	return fmt.Sprintf("resource:%s", name)
}
