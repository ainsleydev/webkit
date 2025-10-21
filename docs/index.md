---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "WebKit"
  text: "Infrastructure as code for full-stack web applications"
  tagline: One manifest. Complete infrastructure. Zero boilerplate.
  actions:
    - theme: brand
      text: Get started
      link: /getting-started/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/ainsleydev/webkit

features:
  - title: Manifest-driven development
    details: Define your entire project—apps, databases, storage, environment variables—in a single app.json file. WebKit handles the rest.
  - title: Complete automation
    details: Generates Terraform configurations, GitHub Actions workflows, Docker setups, and all project scaffolding automatically.
  - title: Built-in secrets management
    details: SOPS and Age encryption integrated out of the box. Manage secrets across environments without the complexity.
  - title: Idempotent updates
    details: Run webkit update anytime. It intelligently updates only what changed, preserving your customisations and cleaning up orphaned files.
  - title: Multi-provider support
    details: Deploy to DigitalOcean, Backblaze B2, and more. WebKit abstracts provider specifics into simple, portable resource definitions.
  - title: Developer-first experience
    details: Docker Compose for local development, monorepo tooling support, and intelligent defaults that get out of your way.
---

## Why WebKit?

Building modern web applications shouldn't require hours of boilerplate setup. You shouldn't have to manually wire together Terraform modules, GitHub Actions workflows, environment variable management, and Docker configurations for every new project.

WebKit solves this by centralising your entire project definition in a single `app.json` manifest. From this manifest, it generates and maintains all the infrastructure and tooling your project needs.

**No more:**
- Copy-pasting Terraform modules between projects
- Manually keeping CI/CD workflows in sync with your infrastructure
- Managing environment variables across multiple `.env` files
- Setting up Docker Compose configurations from scratch

**Instead:**
- Define your apps and resources once
- Run `webkit update` to generate everything
- Change your manifest and update again—WebKit handles the diff

```json
{
  "project": {
    "name": "my-website",
    "title": "My Website"
  },
  "resources": [
    {
      "name": "db",
      "type": "postgres",
      "provider": "digitalocean"
    }
  ],
  "apps": [
    {
      "name": "cms",
      "type": "payload",
      "path": "services/cms"
    }
  ]
}
```

That's it. WebKit generates the Terraform, workflows, Docker configs, and environment files automatically.

## What makes WebKit different?

### It's opinionated, but flexible

WebKit makes sensible decisions about project structure, tooling, and conventions. This reduces decision fatigue and ensures consistency across projects. But you're not locked in—generated files can be customised, and WebKit preserves your changes on subsequent updates.

### It treats infrastructure as a product

Infrastructure shouldn't be an afterthought. WebKit puts infrastructure configuration alongside application code, versioned in git, and reviewable in pull requests. Your infrastructure evolves with your application.

### It's built for real projects

WebKit was created to manage production web applications at [ainsley.dev](https://ainsley.dev). It handles complex scenarios: multiple apps sharing resources, environment-specific configurations, secrets management, and multi-stage deployments.

## Ready to get started?

<div class="vp-doc" style="text-align: center; margin-top: 2rem;">
  <a href="/getting-started/installation" class="vp-button brand" style="display: inline-block; margin: 0 0.5rem;">Get started →</a>
  <a href="/core-concepts/overview" class="vp-button alt" style="display: inline-block; margin: 0 0.5rem;">Learn the concepts</a>
</div>
