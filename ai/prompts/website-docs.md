# Website Docs

The main organisational website -> https://ainsley.dev, hosted at github.com/ainsleydev/website now
has a developer guidelines section located at https://ainsley.dev/guidelines. Attached is the
structure of these new guidelines.

Everytime these docs change, a dispatch repository event is triggered, which triggers an event in
the WebKit repo under `.github/workflows/dispatch-guidelines.yaml`.


- HTML
	- General
- SCSS
	- General
	- Naming
- Go
	- General
	- Comments
		- Function Patterns
	- Constructors & Funcs
	- Control Flow
	- Errors
	- Testing
- JS
	- General
	- Testing
- Git
	- Commits
	- Pre-Commit Checklist
- SvelteKit
	- General
	- Routing
- Payload
	- General
	- Fields
	- Hooks

## Context

Currently, WebKit generates it's repo documentation by a template located in webkit

- We
- Create a PR on the WebKit repo with the ainsley.dev Bot, you can see an example of this below.

```yaml
  - name: Create GitHub App token
	uses: actions/create-github-app-token@v2
	id: app-token
	with:
	  app-id: 2161597
	  private-key: ${{ secrets.ORG_GITHUB_APP_PRIVATE_KEY }}

  - name: Create release PR or publish to npm
	id: changesets
	uses: changesets/action@v1
	with:
	  createGithubReleases: false
	  version: pnpm changeset:version
	  publish: pnpm changeset:publish
	  commit: 'chore: Updating package versions'
	  title: 'chore: Updating package versions'
	env:
	  GITHUB_TOKEN: ${{ steps.app-token.outputs.token }}
	  NPM_TOKEN: ${{ secrets.ORG_NPM_TOKEN }}
```

https://ainsley.dev/guidelines/index.json

## Rules

- Observe the `AGENTS.md` file before you start coding or planning.
- Always make sure a `TODO` checklist is created before any implementation is carried out.
