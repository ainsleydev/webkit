# Troubleshooting

## Migration check failures

If `pnpm migrate:create` fails with "Dependencies out of sync":

```bash
pnpm install
```

This happens when you pull changes but forget to install updated dependencies. The migration check compares `pnpm-lock.yaml` against the cached lockfile in `node_modules/.pnpm/lock.yaml` to ensure consistency between local and CI environments.
