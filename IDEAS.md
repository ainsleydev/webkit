# Ideas

-   Validate HTML5 input once rendered - See: https://codeberg.org/derat/validate
-   Validate Schema.org types using go get github.com/xeipuuv/gojsonschema - See: "https://schema.org/docs/jsonldcontext.json"
-   To Schema, navigation items, rename to Html package insteal of markip
-   The adapter has methods that returns html elements or structures, payload and static can implement them.
-   Have a internal/website/definition file for the service json
-   How about web as a package name?

## Infrastructure

### Make Terraform credentials optional for read-only operations

**Context:**
Currently, `infra.NewTerraform()` requires all Terraform environment variables (DO_API_KEY, BACK_BLAZE_*, GITHUB_TOKEN*, etc.) even when we only need to read Terraform outputs. This forces CI/CD pipelines to have full infrastructure management credentials just to resolve resource references like `db.connection_url`.

**Problem:**
- CI environments need all 8 Terraform credentials to run `env generate`
- This violates the principle of least privilege
- Exposes more credentials than necessary to GitHub Actions runners
- Couples environment variable generation to full Terraform access

**Proposed solution:**
1. Split `TFEnvironment` into two structs:
   - `TFBackendEnv`: Credentials for accessing Terraform state (S3/Backblaze backend)
   - `TFProviderEnv`: Credentials for managing infrastructure (DigitalOcean, Backblaze resources)

2. Add operation modes to `NewTerraform()`:
   - `TerraformModeReadOnly`: Only requires backend credentials, can read state/outputs
   - `TerraformModeManage`: Requires all credentials, can plan/apply/destroy

3. Update `env generate` to use read-only mode:
   ```go
   tf, err := infra.NewTerraform(ctx, appDef, manifest, infra.TerraformModeReadOnly)
   ```

**Benefits:**
- CI only needs backend access (S3 credentials) to read outputs
- More secure: separate credentials for reading vs managing infrastructure
- Local development still works (developers have all credentials)
- Better separation of concerns: reading state ≠ managing infrastructure

**Implementation considerations:**
- Backwards compatibility: default to current behaviour if mode not specified
- Update `ParseTFEnvironment()` to accept optional requirements
- Add documentation for minimal CI credentials setup
- Consider using Terraform Cloud/remote state for even better security

**Status:** Currently using full credentials everywhere (implemented quick fix). This is a future architectural improvement.
