package env

import "github.com/ainsleydev/webkit/internal/appdef"

// hasResourceReferences checks if the definition contains any environment
// variables with source="resource" that need Terraform outputs.
func hasResourceReferences(def *appdef.Definition) bool {
	// Check shared environment.
	hasResource := false
	def.Shared.Env.Walk(func(entry appdef.EnvWalkEntry) {
		if entry.Source == appdef.EnvSourceResource {
			hasResource = true
		}
	})
	if hasResource {
		return true
	}

	// Check app environments.
	for _, app := range def.Apps {
		app.Env.Walk(func(entry appdef.EnvWalkEntry) {
			if entry.Source == appdef.EnvSourceResource {
				hasResource = true
			}
		})
		if hasResource {
			return true
		}
	}

	return false
}
