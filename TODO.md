# TODO

## Validation

We need to utilise one of the following packages to transform our `internal/appdef/definition.go` to
a JSON schema.

- https://github.com/invopop/jsonschema
- https://github.com/swaggest/jsonschema-go

Leaning towards the latter as it's a bit more verbose.

**Flow**:

- Add JSON schema decorations to the structures.
- Create a new `webkit validate` command which will ensure the `app.json` is valid and true.
- Add the same validation to `wekit update`, so it's validate every time a user updates.
- Ensure proper testing.

**Required Validation**:

- Validate domains in app specs, they should not contain https.
- Validate .Path on App and ensure it exists.
- Validate that terraform-managed VM apps (.Infra.Type == "vm" (or app) && .IsTerraformManaged())
  must have at least one domain in .Domains array.
- Validate that domain names in .Domains should not contain protocol prefixes (e.g., "https://").
- Validate these issues with env.
- 
```
Run ./webkit env generate \
Fetching Terraform outputs...
resolving app "cms" env: terraform output not found for environment 'production', resource 'https://ams3', output 'digitaloceanspaces.com' (referenced by key 'S3_ENDPOINT')
Generated .env file for cms
****
```

## Documentation

Create and update the `docs` folder with coherent documentation for WebKit.

## README Generation

Create beautiful looking README's from the `app.json` data.

## Misc

- BetterStack/OneUptime Providers for Infra.
- Improve Coverage.
- Improve path matching for GitHub. Why should we run test and lint if it’s not?
- Create an infra plan —destroy command. So we can see whats destroyed?
- Add .dockerignore to all apps.
- If the app is Payload, add or make a public folder so Next doesnt fail builds.
