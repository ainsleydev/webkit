## CLI Design

Below defines the primary commands for webkit. They will never be destructive by default and users are not allowed to
overwrite webkit changes with commands like `--force`.

- `webkit init` — Interactive creation of app.json.
- `webkit new [path]` — Scaffold new project from templates.
- `webkit update` — Refresh generated files; safe and idempotent by default.
- `webkit validate` — Validate app.json against schema.
- `webkit generate` — Generate files from templates.
- `webkit encrypt` & `webkit decrypt` — Manage secrets via SOPS`
- `webkit infra plan|apply` — Optional infra provisioning hooks (TBC)

## Templates and Generated Files

WebKit ships with a set of templates to generate standardised project files. This makes it easier for engineers to focus
on business problems instead of editor configuration, linting and formatting standards.

- Static files (like .editorconfig) can be copied directly.
- Templated files (like GitHub workflows) use Go’s text/template engine + Sprig helpers (only if needed).

**Files WebKit should generate:**

- `.editorconfig`
- `.prettierrc`
- `.prettierignore`
- `.dockerignore`
- `.github/workflows/*`
- `.github/dependabot`
- `.github/settings.yml` (repo config)
- `README.md`