# Monitoring with Uptime Kuma

WebKit automatically configures Uptime Kuma monitoring for applications and resources.

## Overview

Monitoring is **enabled by default** (opt-out) for all apps and resources with Terraform management. WebKit creates three types of monitors:

1. **HTTP monitors** - Monitor application uptime via HTTP/HTTPS requests
2. **Postgres monitors** - Monitor database connection health
3. **Push monitors** - Heartbeat tracking for CI backup jobs (TODO)

## Configuration

### Enabling/Disabling Monitoring

Monitoring is enabled by default. To disable for a specific app or resource:

```json
{
  "apps": [
    {
      "name": "web",
      "monitoring": {
        "enabled": false
      }
    }
  ],
  "resources": [
    {
      "name": "db",
      "type": "postgres",
      "monitoring": {
        "enabled": false
      }
    }
  ]
}
```

### Uptime Kuma Configuration

Add the following to your `.env.production.enc` file (encrypted with SOPS):

```bash
TF_VAR_uptime_kuma_username=admin
TF_VAR_uptime_kuma_password=<your-password>
```

The Uptime Kuma URL is configured in `platform/terraform/base/variables.tf`:

```hcl
variable "uptime_kuma_url" {
  default = "https://uptime.ainsley.dev"
}
```

## What Gets Monitored

### Applications

- **All domains** (primary + aliases) are monitored via HTTP/HTTPS
- Health check endpoint: Uses `health_check_path` from app infra config (default: `/`)
- Check interval: Every 60 seconds
- Expected status: 200 OK
- Retry interval: 60 seconds
- Max retries: 3

Example monitor created for app with multiple domains:

```json
{
  "apps": [
    {
      "name": "web",
      "domains": [
        { "name": "example.com", "type": "primary" },
        { "name": "www.example.com", "type": "alias" }
      ]
    }
  ]
}
```

This creates two monitors:
- `{project-name}-web-example-com`
- `{project-name}-web-www-example-com`

### Resources

**Postgres databases** are monitored via connection health checks:

- Check interval: Every 5 minutes (300 seconds)
- Retry interval: 60 seconds
- Max retries: 3
- Connection URL: Sourced from Terraform module outputs

**Other resource types** (S3, SQLite) are not yet supported.

### Backup Jobs (TODO - Phase 5)

Heartbeat monitoring for CI backup jobs is planned but not yet implemented. The infrastructure is in place:

- Push monitors are created for each resource with backups enabled
- Expected interval: Auto-calculated from cron schedule + 10% buffer
- Daily backups (2am): 26.4 hour heartbeat window

**TODO:**
- Export push monitor URLs from Terraform
- Store URLs as GitHub secrets
- Update backup workflow to ping heartbeat URL on success

## Monitor Details

### HTTP Monitor Configuration

```
Name:            {project-name}-{app-name}-{domain}
Type:            http
URL:             https://{domain}{health_check_path}
Method:          GET
Expected Status: [200]
Interval:        60s
Retry Interval:  60s
Max Retries:     3
TLS Validation:  Enabled
```

### Postgres Monitor Configuration

```
Name:            {project-name}-{resource-name}-{environment}
Type:            postgres
Database URL:    (from Terraform outputs)
Interval:        300s
Retry Interval:  60s
Max Retries:     3
```

### Push Monitor Configuration (Planned)

```
Name:             {project-name}-backup-{resource-name}
Type:             push
Expected Interval: 95040s (26.4 hours for daily backups)
Max Retries:      2
```

## Notifications

All monitors can send alerts to configured notification channels in Uptime Kuma.

**TODO:** Automate Slack notification configuration via Terraform (currently manual in UI).

## Implementation Details

### Architecture

1. **Schema Layer** (`schema.json`):
   - Simple `MonitoringConfig` type with single `enabled` field

2. **Appdef Layer** (`internal/appdef/monitor.go`):
   - `Monitor` struct with sophisticated defaults
   - Monitor generation logic for apps and resources
   - Smart interval calculation based on monitor type

3. **Terraform Layer** (`internal/infra/tf_vars.go`):
   - Transforms appdef.Monitor to tfMonitor
   - Generates monitor list for Terraform variables

4. **Infrastructure Layer** (`platform/terraform/modules/monitoring`):
   - Creates Uptime Kuma monitors via Terraform provider
   - Separates HTTP, Postgres, and Push monitors
   - Outputs monitor IDs and details

### Monitor Generation Flow

```
app.json (user config)
  ↓
App/Resource structs (defaults applied)
  ↓
GenerateMonitors() methods
  ↓
appdef.Monitor structs (with smart defaults)
  ↓
tfMonitor transformation
  ↓
Terraform variables (webkit.auto.tfvars.json)
  ↓
monitoring module
  ↓
Uptime Kuma monitors created
```

## Troubleshooting

### No monitors created

Check that monitoring is enabled (default) and apps have domains or resources are Postgres:

```bash
# Check generated Terraform variables
cat .webkit/terraform/base/webkit.auto.tfvars.json | jq '.monitors'
```

### Uptime Kuma authentication fails

Verify credentials in `.env.production.enc`:

```bash
# Decrypt and check
sops -d .env.production.enc | grep UPTIME_KUMA
```

### Terraform provider errors

The Uptime Kuma provider requires the Web API adapter. Verify:

1. Uptime Kuma Web API is running at configured URL
2. Credentials are correct
3. Provider version is compatible

## Future Enhancements

1. **Expose more configuration options** in `app.json`:
   - Custom check intervals
   - Expected status codes
   - Custom health check paths

2. **Support more resource types**:
   - S3 bucket monitoring (DNS/HTTP checks)
   - SQLite/Turso monitoring (HTTP API checks)
   - Redis monitoring

3. **Automated notification management**:
   - Create Slack notifications via Terraform
   - Auto-link project Slack channel

4. **Status pages**:
   - Auto-create public status pages
   - Group monitors by project

5. **Backup heartbeat integration** (Phase 5):
   - Complete GitHub Actions workflow integration
   - Automatic heartbeat pinging on successful backups
