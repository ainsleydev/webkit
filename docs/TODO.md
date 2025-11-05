# TODO & Known Issues

This document tracks known issues, caveats, and future improvements for WebKit.

## DigitalOcean App Platform

### ✅ Environment Variable Drift - RESOLVED

**Issue:** Perpetual Terraform drift caused by DO's server-side encryption of environment variables.
**Solution:** Implemented smart lifecycle management using `ignore_changes` + `replace_triggered_by` pattern.
**Reference:** See `platform/terraform/providers/digital_ocean/app/main.tf` for implementation details.

Related issue: https://github.com/digitalocean/terraform-provider-digitalocean/issues/869

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
