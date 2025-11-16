<p align="center">
  {{- if .LogoURL }}
  <img src="{{ .LogoURL }}" height="96">
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
{{- range .AppTypeBadges }}
![{{ .Name }}](https://img.shields.io/badge/{{ .Name }}-active-{{ .Color }})
{{- end }}
{{- range .DomainBadges }}
[![{{ .Name }}](https://img.shields.io/website?url=https://{{ .Name }})](https://{{ .Name }})
{{- end }}

</div>

---

## {{ .Definition.Project.Title }}

{{ .Definition.Project.Description }}

Built with [WebKit {{ .Definition.WebkitVersion }}](https://github.com/ainsleydev/webkit).

{{- if .Content }}

{{ .Content }}
{{- end }}

---

## Repository Structure

This repository contains **{{ len .Definition.Apps }} application{{ if ne (len .Definition.Apps) 1 }}s{{ end }}**{{ if .Definition.Resources }} and **{{ len .Definition.Resources }} resource{{ if ne (len .Definition.Resources) 1 }}s{{ end }}**{{ end }}.

### Applications
{{- range .Definition.Apps }}

#### {{ .Title }}

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

---

## Resources
{{- range .Definition.Resources }}

### {{ .Name }}

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

---

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
   webkit install
   ```

3. Start development:
   ```bash
   webkit up
   ```


4. Access applications:
{{- range .Definition.Apps }}
{{- if gt .Build.Port 0 }}
   - {{ .Title }}: http://localhost:{{ .Build.Port }}
{{- end }}
{{- end }}

---

## Deployment

Deployment is managed by WebKit using:
- **Infrastructure:** Terraform
- **CI/CD:** GitHub Actions
- **Environments:** Development, Staging, Production

{{- if .ProviderGroups }}

**Hosting Providers:**
{{- range $provider, $items := .ProviderGroups }}
- **{{ $provider }}:** {{ $items }}
{{- end }}
{{- end }}

---

## License

Code Copyright {{ .CurrentYear }} {{ .Definition.Project.Repo.Owner }}. Code released under the [BSD-3 Clause](LICENSE).

---

<p align="center">
  Built with <a href="https://github.com/ainsleydev/webkit">WebKit</a>
</p>
