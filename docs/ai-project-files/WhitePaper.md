# WebKit


## Secrets and Environment Strategy

This whole section is to be considered and not confirmed.

- Secrets are not inlined in `app.json` Instead: They are referenced via a from field (`sops:`, `github-secrets:`,
  `vault:`).
- SOPS is the default backend: plaintext YAML → encrypted YAML checked into secrets. `webkit encrypt` and
  `webkit decrypt` manage the lifecycle.
- GitHub Action hook ensures no plaintext secrets are committed.

## Implementation Conventions

- Language: Go.
- Templates: Embedded via embed.FS, overridable by local templates/.
- Template context: includes Project, Apps, Resources, Env, Shared.
- Testing: `webkit validate` ensures manifests conform to schema.
- Generated file tracking: internally tracked by WebKit, not exposed in manifest.

---

