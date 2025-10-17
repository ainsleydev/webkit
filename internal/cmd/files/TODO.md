# TODO (CI)

Would be good to include some sort of file matching for PRs:

```yaml
on:
  pull_request:
    paths:
      - '**/*.ts'
      - '**/*.js'
      - '**/*.svelte'
      - '**/.eslintrc*'
      - 'pnpm-lock.yaml'
```

Consider more specific cache keys:

```yaml
- name: Set up Node
  uses: actions/setup-node@v4
  with:
    node-version: '22'
    cache: 'pnpm'
    cache-dependency-path: '{{ .Path }}/pnpm-lock.yaml'
```


If you're using Turborepo (based on your templates), you could optimize:

```yaml
- name: Lint
  working-directory: {{ .Path }}
  run: pnpm turbo lint --filter={{ .Name }}
```
