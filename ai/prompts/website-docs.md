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

The reason why they're hosted on ainsley.dev site is so that it can be viewed by external parties,
and it's not just in markdown. WebKit, should now use these documents on ainsley.dev to create some
really useful artifacts for AI agents, for both the root WebKit repo, and WebKit enabled
repositories.

When one of these repository events gets triggered, it's the sign for WebKit to update all of it's
documentation then create a PR on the repo using the ainsley.dev bot.

As ainsley.dev's site is built in Hugo (GoLang), functionality has been added to get these guideline
entries with ease. They are ordered in exactly the same way as on the site, example below.

**Guidelines Content**:

Available at: `https://ainsley.dev/guidelines/index.json`

```json
[
	{
		"date": "2025-10-27T00:00:00Z",
		"description": "HTML formatting standards including indentation, quotes & self-closing tags",
		"draft": false,
		"heading": "General",
		"lastmod": "2025-10-27T00:00:00Z",
		"markdown": "\n## Validity\n\nAll HTML should be using the [Markup Validat.... (shortened)",
		"permalink": "https://ainsley.dev/guidelines/html/general/",
		"plainContent": " Validity All HTML should be using the Markup Validation Service before (shortened)",
		"publishdate": "2025-10-27T00:00:00Z",
		"section": "HTML",
		"subsection": "html",
		"summary": "Validity All HTML should be using the Markup Validation Service before creating a pull request or pushing to production. This will help avoid common mistakes such as closing tags, wrong attributes and many more.\nBy validating HTML it ensures that web pages are consistent across multiple devices and platforms and increases the chance of search engines to properly pass markup.\nIndentation Use tabs instead of spaces for markup. Do not mix tabs with spaces, ensure it is probably formatted.",
		"title": "General",
		"url": "/guidelines/html/general/",
		"weight": 2
	}
]
```

## TODO

- Create a go script in cmd which will be executed by the dispatch-guidelines workflow which will
  unmarshal into a Go struct aligned with the example above.
	- Reside in `cmd/docs/main.go`
	- HTTP GET https://ainsley.dev/guidelines/index.json
	- The script should unmarshal into the type and generate the documentation.
	- Then generate the markdown file, we probably need to convert all headings to be lower by one
	  as the website will have `## H2s`, but we probably need to nest it.
	- We should test this.
- Update `dispatch-guidelines.yaml` which will create a new PR in the webkit repo updating the
  documentation.
	- We should use the ainsley.dev bot for this, i.e ainsley.dev bot will be the author of the PR.
	- Example in `publish.yaml` and below.
- We should also have sub directories depending on the app type. For example, if the app is payload,
  and the path is `cms` then another `AGENTS.md` file will be generated within `cms` that contains
  Payload specific code style, the same can be true for `svelte`
	- This also means that we don't need to include `svelte` or `payload` code style in the root
	  `AGENTS.md`, saving space.
- I want to mention to Agents that if I say we should update docs, then the ainsley.dev/website repo
  should be cloned and edited.
- Come up with a plan on how we can tackle above.

### ainsley.dev Bot Example

```yaml
  -   name: Create GitHub App token
  uses: actions/create-github-app-token@v2
  id: app-token
  with:
	  app-id: 2161597
	  private-key: ${{ secrets.ORG_GITHUB_APP_PRIVATE_KEY }}

		  -   name: Create release PR or publish to npm
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

## Resuing AGENTS.md

In addition to above, we also need to refactor something:

Currently, there is a file in internal/templates/AGENTS.md that is used for services / apps (webkit
enabled) to create Agents files and inject in their own template by creating a file in their own
repo under docs (you can see this under internal/cmd/docs/cmd.go)

However, webKit's AGENTS.md shares very similar characteristics such as code formatting and other
bits. Essentially, I want WebKit to use this too.

Just like it's child repos, Webkit should define our template in docs/AGENTS.md where we can place 
webkit specific context.

## Rules

- Observe the `AGENTS.md` file before you start coding or planning.
- Always make sure a `TODO` checklist is created before any implementation is carried out.
