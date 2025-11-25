<p align="center">
  {{- if .LogoURL }}
  <img src="{{ .LogoURL }}" width="200">
  {{- end }}
  <h3 align="center">{{ .Definition.Project.Title }}</h3>
</p>

<p align="center">
  {{ .Definition.Project.Description }}
</p>

{{- if .DomainLinks }}

<p align="center">
  {{ .DomainLinks }}
</p>
{{- end }}

<div align="center">

![WebKit](https://img.shields.io/badge/webkit-{{ .Definition.WebkitVersion }}-blue)
[![Backup](https://github.com/{{ .Definition.Project.Repo.Owner }}/{{ .Definition.Project.Repo.Name }}/actions/workflows/backup.yaml/badge.svg)](https://github.com/{{ .Definition.Project.Repo.Owner }}/{{ .Definition.Project.Repo.Name }}/actions/workflows/backup.yaml)
[![PR](https://github.com/{{ .Definition.Project.Repo.Owner }}/{{ .Definition.Project.Repo.Name }}/actions/workflows/pr.yaml/badge.svg)](https://github.com/{{ .Definition.Project.Repo.Owner }}/{{ .Definition.Project.Repo.Name }}/actions/workflows/pr.yaml)
[![Release](https://github.com/{{ .Definition.Project.Repo.Owner }}/{{ .Definition.Project.Repo.Name }}/actions/workflows/release.yaml/badge.svg)](https://github.com/{{ .Definition.Project.Repo.Owner }}/{{ .Definition.Project.Repo.Name }}/actions/workflows/release.yaml)
[![ainsley.dev](https://img.shields.io/badge/-ainsley.dev-black?style=flat&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABGdBTUEAALGPC/xhBQAAACBjSFJNAAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAABmJLR0QA/wD/AP+gvaeTAAAAB3RJTUUH5wEYDzUGL1b35AAABA1JREFUWMPtlttvFVUUxn977ZnZu+W0tLalqRovBAUvQag0xNQbpSIosSSIJC198YknJfHJxDf9A/DBJ0x8MbFACjVqvCASq6FYFLFBvJAaAomkFCmhHGpLO+PDzOmZzpn2nKP4pCs5ycmevb7vW99as/fA//FfD1XO5p1nzuA3NWJHx5T8cVkRBPHHQfRjd0tzyZhOOQIy27bAxET9zCuvvhY0r2kC/OiRABeAN4BL/4oDr9+3lGszPs7UVNfUE23v3Nj5koszR/8N4EXg3XJckFIFuCLUuU7GWNNtTg25cu4syJx0F+gGMuU4UJKAt1Yux1UKV6TVat1qs+OYwQESMwDQCjwKsOv4iZsnwGihwbiuEek2WjJGhMrvv0UujYKa08VFkQvuTXNgz6oVeCIo1CqrZYMRwTiaytERKn44kRQAsAFYDbBrsLgLRQU0GI919TXKiHQaUQ1GBCuCCQKqjg/MqInrM4lZrgc6A1CljHhRAZ4Ip65m77FaOmbJdehC5vzZr1RAf/T6x6NDwb3/uAVfP74GnwCjZasRuXuWXASj9XQme+3t6erqPcB0IvUuYCsUH8YFBRhRNBqvyYpsn0MeOnG6wvc/9x33MPBjSvp24Na/7cDP7Y/gKIURecZoeTBObkSwWg7UNjaOeFfGLgK9KRAPAM8Wc2FeAUaEWtddbEV2WBFtREXkCqvlghE5yOQkvucBHAR+T0BooAtYXLYDI5sewxWFJ/Kk1bI2UTlW5DMFp03+JPwJ+DQFai2wbiEXUgVUas0trmuslm4jUmGi/tuwDVmrpafBuNPVrs7N/wzQA2QTUJbwYLIlOxB0tOGJ4IhqsSJts+T54Rv0lBz1RFh9ZJA385fOAHAshaMNaAF4OcWFQgeUwhMlrlJdnqjaOLkR8Y2WvbWec9VIQeo4sJf8FZ2LmmgWJO1cmm8I7wc2a6XwosGL+v+rFfnYUYplh47Obo5dvZ8Av6TgbSZ8KxYWEGxZn/u7Dbg9t8HNnwF9S2qqzqVUn4vzQF/K+m3AC1A4jGlId0QC8l0BXKVGrahe//okNR99WZAUc6EXuJiC+zxw57wOxKp/DliRAvCFKDUkxS+YIeBwyvryCHuOC0kH6oBOCj/V/gTeA6aK0oefZj3ARGJdRdh1BQ7Eqm8HHk4B/Q7oB1B9acWFEWtDf5STjGbgqbgLcQcqCQ8NL5EUAPuBsRKqz8UVYB+F97QXcSyatSXoWJ8zvB04AFQlkoaBp4HhhaqPR1TdUsLjeVni8TjhVX0odCAkd4AdKeQAHxIwXEb1Odt+Az5IeVQVcTmhgDBWAhtTNl8G9qGAwKfU2N3SnJvi/RFGMjYCD8UFdACNKRsHgZMA6v0j5ZpAlPtNyvqSiJO/AKik60y0ALlUAAAAJXRFWHRkYXRlOmNyZWF0ZQAyMDIzLTAxLTI0VDE1OjUzOjA2KzAwOjAwm5vntAAAACV0RVh0ZGF0ZTptb2RpZnkAMjAyMy0wMS0yNFQxNTo1MzowNiswMDowMOrGXwgAAABXelRYdFJhdyBwcm9maWxlIHR5cGUgaXB0YwAAeJzj8gwIcVYoKMpPy8xJ5VIAAyMLLmMLEyMTS5MUAxMgRIA0w2QDI7NUIMvY1MjEzMQcxAfLgEigSi4A6hcRdPJCNZUAAAAASUVORK5CYII=)](https://ainsley.dev)
[![Twitter Handle](https://img.shields.io/twitter/follow/ainsleydev)](https://twitter.com/ainsleydev)

</div>

## {{ .Definition.Project.Title }}

{{ if .Content -}}
{{ .Content }}
{{- end }}

{{ .Definition.Project.Description }}

Built with [WebKit {{ .Definition.WebkitVersion }}](https://github.com/ainsleydev/webkit).

This repository contains **{{ len .Definition.Apps }} application{{ if ne (len .Definition.Apps) 1 }}s{{ end }}**{{ if .Definition.Resources }} and **{{ len .Definition.Resources }} resource{{ if ne (len .Definition.Resources) 1 }}s{{ end }}**{{ end }}.

{{- if .MonitorBadges }}

## Status

Uptime monitors for the application. Visit the [status page]({{ .StatusPageURL }}) for more details.
{{- if .DashboardURL }}

View all monitors on the [dashboard]({{ .DashboardURL }}).
{{- end }}

| Monitor | Status |
|---------|--------|
{{- range .MonitorBadges }}
| {{ .Name }} | [![Status]({{ .BadgeURL }})]({{ $.StatusPageURL }}) |
{{- end }}
{{- end }}

## Apps

{{- range .Definition.Apps }}

### {{ .Title }}

{{ if eq .Type "payload" }}
<img src="https://img.shields.io/badge/-Payload CMS-000000?style=flat&logo=payloadcms&logoColor=white"/>
{{- else if eq .Type "svelte-kit" }}
<img src="https://img.shields.io/badge/-Svelte-FF3E00?style=flat&logo=svelte&logoColor=white"/>
{{- end }}

**Type:** {{ .Type }}{{ if .Build.Port }} | **Port:** {{ .Build.Port }}{{ end }}{{ if .PrimaryDomain }} | **Domain:** {{ .PrimaryDomain }}{{ end }}

{{ .Description }}

{{- if .OrderedCommands }}

**Commands:**
| Command | Script | CI/CD |
|---------|--------|-------|
{{- range .OrderedCommands }}
| {{ .Name }} | `{{ .Cmd }}` | {{ if .SkipCI }}No{{ else }}Yes{{ end }} |
{{- end }}
{{- end }}

{{- $tools := .InstallCommands }}
{{- if $tools }}

**Tools:**
{{- range $tools }}
- `{{ . }}`
{{- end }}
{{- end }}

{{- if .Infra.Provider }}

**Infrastructure:**
- **Provider:** {{ .Infra.Provider }}
- **Type:** {{ .Infra.Type }}
{{- range $key, $value := .Infra.Config }}
- **{{ $key }}:** {{ $value }}
{{- end }}
{{- end }}
{{- end }}

{{- if .Definition.Resources }}

## Resources

{{- range .Definition.Resources }}

### {{ .Title }}

{{ if eq .Type "postgres" }}
<img src="https://img.shields.io/badge/-PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white"/>
{{- else if eq .Type "s3" }}
![DigitalOcean Spaces](https://img.shields.io/badge/DigitalOcean-Spaces-009EE0?style=flat&logo=digitalOcean&logoColor=white)
{{- else if eq .Type "sqlite" }}
<img src="https://img.shields.io/badge/-SQLite-003B57?style=flat&logo=sqlite&logoColor=white"/>
{{- end }}

**Type:** {{ .Type }} | **Provider:** {{ .Provider }} | **Backups:** {{ if .Backup.Enabled }}Enabled{{ else }}Disabled{{ end }}
{{- if .Description }}

{{ .Description }}
{{- end }}

**Configuration:**
| Setting | Value |
|---------|-------|
{{- range $key, $value := .Config }}
| {{ $key }} | {{ $value }} |
{{- end }}

**Available Outputs:**
{{- range .Type.Documentation }}
- `{{ .Name }}` - {{ .Description }}
{{- end }}
{{- end }}
{{- end }}

## Development

### Prerequisites

{{- if .Definition.ContainsGo }}
- Go (latest)
{{- end }}
{{- if .Definition.ContainsJS }}
- Node.js and pnpm
{{- end }}
- Docker and Docker Compose
- WebKit CLI

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/{{ .Definition.Project.Repo.Owner }}/{{ .Definition.Project.Repo.Name }}.git
   cd {{ .Definition.Project.Repo.Name }}
   ```

2. Install dependencies:
   ```bash
   pnpm install
   ```

3. Start development:
   ```bash
   pnpm dev
   ```

4. Access applications:
{{- range .Definition.Apps }}
{{- if gt .Build.Port 0 }}
   - {{ .Title }}: http://localhost:{{ .Build.Port }}
{{- end }}
{{- end }}

## Deployment

Deployment is managed by WebKit using:

- **Infrastructure:** Terraform
- **CI/CD:** GitHub Actions
- **Environments:** Development, Production

{{- if .ProviderGroups }}

**Hosting Providers:**
{{- range $provider, $items := .ProviderGroups }}
- **{{ $provider }}:** {{ $items }}
{{- end }}
{{- end }}

## License

Code Copyright {{ .CurrentYear }} {{ .Definition.Project.Repo.Owner }}. Code released under the [BSD-3 Clause](LICENSE).

___

<p align="center">
  Built with <a href="https://github.com/ainsleydev/webkit">WebKit</a>
</p>
