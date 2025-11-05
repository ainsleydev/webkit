# TODO & Known Issues

This document tracks known issues, caveats, and future improvements for WebKit.

## DigitalOcean App Platform

### Environment Variable Management

**Issue:** DigitalOcean App Platform encrypts environment variables server-side, which causes Terraform state comparison issues.

**Context:**
When environment variables are set via Terraform, DigitalOcean encrypts them on their end. This means:
1. Terraform stores the plain-text value in state
2. DigitalOcean returns the encrypted value via API
3. Terraform detects a diff on every `terraform plan` (plain text vs encrypted)

This is a known issue in the DigitalOcean Terraform provider:
- https://github.com/digitalocean/terraform-provider-digitalocean/issues/869

**Current Solution:**
All environment variables (including secrets) are managed via Terraform. The deployment workflow in GitHub Actions only triggers deployments and does not modify environment variables.

To avoid perpetual Terraform drift, use lifecycle ignore rules:

```hcl
resource "digitalocean_app" "example" {
  # ... app configuration

  lifecycle {
    ignore_changes = [
      spec[0].service[0].env
    ]
  }
}
```

**Alternative Approaches Considered:**

1. **GitHub Secrets with JSON Objects** (Rejected - Too Complex)
   - Store secrets as JSON: `TF_DO_APP_PLATFORM_WEB_SECRETS`
   - Parse and inject via workflow using `yq` and `doctl`
   - Pros: Secrets not in Terraform state
   - Cons: Complex bash scripting, error-prone, JSON parsing in CI, unclear separation of concerns

2. **Individual GitHub Secrets** (Rejected - Maintenance Overhead)
   - Per-secret GitHub secrets: `DO_WEB_PAYLOAD_SECRET`, `DO_WEB_SENTRY_DSN`
   - Pros: Easy to manage in GitHub UI
   - Cons: Manual secret creation, hardcoded names in workflow, duplication

3. **Skip Secret Management in CI** (Rejected - Manual Process)
   - Manage secrets manually via DigitalOcean UI or CLI
   - Pros: Simple workflow
   - Cons: No automation, manual overhead

4. **Terraform Only** (Current Implementation - Recommended)
   - All env vars managed via Terraform
   - Secrets encrypted with SOPS/Age
   - Workflow only triggers deployment
   - Pros: Single source of truth, simple workflow, integrates with existing SOPS
   - Cons: Secrets visible in Terraform state (but state should be encrypted)

**Recommendation:**
Use Terraform for all environment variable management, including secrets. This maintains a single source of truth and keeps the deployment workflow simple. Ensure Terraform state is properly secured (encrypted backend, restricted access).

### Future Improvements

- [ ] Document Terraform lifecycle ignore patterns for DO App Platform
- [ ] Add example Terraform configuration for DO App Platform with env vars
- [ ] Consider contributing fix to DigitalOcean Terraform provider
- [ ] Investigate if DO API supports reading plain-text values for comparison

## GitHub Actions Secrets

### Token Inheritance

Secrets can be defined at repository or organisation level. The deployment workflow supports fallback:

```yaml
token: ${{ secrets.REPO_DO_ACCESS_TOKEN || secrets.ORG_DO_ACCESS_TOKEN }}
```

This allows:
- Repository-level override: Set `REPO_DO_ACCESS_TOKEN` in specific repo
- Organisation-level default: Set `ORG_DO_ACCESS_TOKEN` for all repos

**Required Secrets:**
- `REPO_DO_ACCESS_TOKEN` or `ORG_DO_ACCESS_TOKEN` - DigitalOcean API token

## Future Platform Support

### Other Cloud Providers

- [ ] AWS ECS/Fargate deployment
- [ ] GCP Cloud Run deployment
- [ ] Azure Container Apps deployment
- [ ] Fly.io deployment
- [ ] Railway deployment

Each platform will need similar research around:
- Environment variable injection
- Secret management
- State comparison issues
- API authentication
