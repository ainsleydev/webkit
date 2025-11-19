# Monitoring with Uptime Kuma

WebKit automatically configures Uptime Kuma monitoring for applications.

## Overview

Monitoring is **enabled by default** (opt-out) for all apps with Terraform management. WebKit currently supports:

1. **HTTP monitors** - Monitor application uptime via HTTP/HTTPS requests

**Note:** Resource monitoring (databases, backup heartbeats) has been temporarily disabled to simplify the initial implementation. It can be re-enabled in the future when needed.

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
- Check interval: Every 60 seconds
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

## Monitor Details

### HTTP Monitor Configuration

```
Name:           {project-name}-{app-name}-{domain}
Type:           http
URL:            https://{domain}
Method:         GET
Interval:       60s
Retry Interval: 60s
Max Retries:    3
TLS Validation: Enabled
Upside Down:    false
Ignore TLS:     false
```

## Notifications

Notifications are configured manually in the Uptime Kuma UI. The `ehealth-co-id/uptimekuma` provider does not currently support automatic notification configuration via Terraform.

**Note:** Notification IDs cannot be linked to monitors through Terraform with the current provider version.

## Implementation Details

### Architecture

1. **Schema Layer** (`schema.json`):
   - Simple `MonitoringConfig` type with single `enabled` field

2. **Appdef Layer** (`internal/appdef/monitor.go`):
   - `Monitor` struct
   - Monitor generation logic for apps

3. **Terraform Layer** (`internal/infra/tf_vars.go`):
   - Transforms appdef.Monitor to tfMonitor
   - Generates monitor list for Terraform variables (only HTTP monitors currently)

4. **Infrastructure Layer** (`platform/terraform/modules/monitoring`):
   - Creates Uptime Kuma monitors via ehealth-co-id/uptimekuma provider
   - Only HTTP monitors currently implemented
   - Outputs monitor IDs and details

### Monitor Generation Flow

```
app.json (user config)
  ↓
App structs (monitoring defaults applied)
  ↓
App.GenerateMonitors()
  ↓
appdef.Monitor structs
  ↓
tfMonitor transformation (name, type, url, method)
  ↓
Terraform variables (webkit.auto.tfvars.json)
  ↓
monitoring module (applies interval, retries, etc.)
  ↓
Uptime Kuma HTTP monitors created
```

## Troubleshooting

### No monitors created

Check that monitoring is enabled (default) and apps have domains:

```bash
# Check generated Terraform variables
cat .webkit/terraform/base/webkit.auto.tfvars.json | jq '.monitors'
```

Monitors are only created for apps with domains. Resources (databases, etc.) are not currently monitored.

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

1. **Re-enable resource monitoring**:
   - Postgres database connection health checks
   - Backup heartbeat monitoring (push monitors)
   - Other database types (MySQL, MongoDB, Redis)

2. **Enhance provider support**:
   - Fork or contribute to ehealth-co-id/uptimekuma provider to add:
     - `notification_id_list` support
     - `accepted_status_codes` support
     - Better database monitor support

3. **Expose more configuration options** in `app.json`:
   - Custom check intervals
   - Expected status codes
   - Custom health check paths

4. **Automated notification management**:
   - Create Slack notifications via Terraform
   - Auto-link monitors to notification channels

5. **Status pages**:
   - Auto-create public status pages
   - Group monitors by project
