# Agent Guidelines

This document provides guidelines for AI agents working on the WebKit codebase.

## Note For Humans

This is a living document that will improve as more people/agents use it over time. Every effort has
been made to keep the guidance in here as generic and reusable as possible. Please keep this in mind
with any future edits.

**Note**: Investigation summaries and debugging analysis should be displayed via UI only, not
committed to the repository.

## Updating Documentation

If you need to update developer guidelines, clone and edit the [ainsley.dev/website](https://github.com/ainsleydev/website) repository.
These guidelines are automatically synced from there.

{{ .Content -}}

{{ .CodeStyle -}}

{{ .Git -}}

{{ if .Payload }}
## Libraries

{{ .Payload -}}
{{ end }}

{{ if .SvelteKit }}
{{ if not .Payload }}## Libraries

{{ end }}{{ .SvelteKit -}}
{{ end }}
