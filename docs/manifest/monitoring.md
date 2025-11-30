# Monitoring with Peekaping

WebKit automatically configures Peekaping monitoring for applications.

## Overview

Monitoring is **enabled by default** (opt-out) for all apps with Terraform management. WebKit currently supports:

1. **HTTP monitors** - Monitor application uptime via HTTP/HTTPS requests
2. **HTTP keyword monitors** - Monitor page content for specific keywords
3. **DNS monitors** - Monitor domain resolution
4. **Push monitors** - Heartbeat monitoring for backup jobs and scheduled tasks
5. **Status pages** - Public status pages showing service health
6. **Slack notifications** - Automatic alerts via Slack webhooks

## Configuration

### Enabling/Disabling Monitoring

Monitoring is enabled by default. To disable for a specific app:

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

### Peekaping Configuration

Add the following to your `.env.production.enc` file (encrypted with SOPS):

```bash
PEEKAPING_ENDPOINT=https://peekaping.example.com
PEEKAPING_API_KEY=<your-api-key>
```

The API key is used by the Terraform provider to authenticate with your Peekaping instance.

### Slack Notifications (Optional)

To enable Slack notifications, configure a webhook URL:

```bash
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/WEBHOOK/URL
```

All monitors will automatically send alerts to this Slack channel when services go down.

## What Gets Monitored

### Applications

- **All domains** (primary + aliases) are monitored via HTTP/HTTPS
- Check interval: Every 60 seconds
- Timeout: 30 seconds
- Max retries: 3
- Retry interval: 60 seconds
- Resend interval: 10 minutes (how often to resend alerts)

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

This creates two HTTP monitors:
- `{project-name}-web-example-com`
- `{project-name}-web-www-example-com`

## Monitor Details

### HTTP Monitor Configuration

```
Name:               {project-name}-{app-name}-{domain}
Type:               http
URL:                https://{domain}
Method:             GET
Encoding:           json
Accepted Status:    2XX
Auth Method:        none
Interval:           60s
Timeout:            30s
Retry Interval:     60s
Max Retries:        3
Resend Interval:    10m
```

### HTTP Keyword Monitor Configuration

HTTP keyword monitors check for the presence or absence of specific text in HTTP responses. This is useful for content validation, error detection, and ensuring important information appears on your pages.

**Configuration Fields**:
- `url` - URL to check (required)
- `method` - HTTP method (required, e.g., GET, POST)
- `keyword` - Text to search for in the response (required)
- `invert_keyword` - If true, alert when keyword IS found; if false, alert when keyword is NOT found (optional, default: false)
- `max_redirects` - Maximum redirects to follow (optional, default: 3)

**Use Cases**:
1. **Content validation** - Ensure important text appears on pages
2. **Error detection** - Alert when error messages appear
3. **Security monitoring** - Detect unauthorized content
4. **Compliance** - Verify required disclaimers are displayed

**Example - Normal Mode (Alert When Missing)**:

```json
{
  "monitoring": {
    "custom": [
      {
        "name": "Homepage Has Tagline",
        "type": "http-keyword",
        "interval": 300,
        "config": {
          "url": "https://example.com",
          "method": "GET",
          "keyword": "craftsmanship"
        }
      }
    ]
  }
}
```

This monitors alerts if "craftsmanship" is NOT found on the homepage.

**Example - Inverted Mode (Alert When Present)**:

```json
{
  "monitoring": {
    "custom": [
      {
        "name": "Error Page Detection",
        "type": "http-keyword",
        "interval": 120,
        "config": {
          "url": "https://example.com/api/health",
          "method": "GET",
          "keyword": "error",
          "invert_keyword": true
        }
      }
    ]
  }
}
```

This monitor alerts if "error" IS found in the API health check response.

### DNS Monitor Configuration

```
Name:               {project-name}-{app-name}-{domain}-dns
Type:               dns
Hostname:           {domain}
Resolver:           1.1.1.1 (Cloudflare)
DNS Type:           A
Interval:           300s (5 minutes)
Timeout:            30s
Retry Interval:     60s
Max Retries:        3
Resend Interval:    10m
```

### Push Monitor Configuration (Heartbeats)

Push monitors are used for backup and maintenance job monitoring. Jobs ping the monitor URL on success, and if no ping is received within the expected interval, an alert is sent.

**Backup Monitors**:

```
Name:               {Project Title} - {Resource Title} Backup
Type:               push
Interval:           25 hours (90000s)
Max Retries:        2
Retry Interval:     60s
Resend Interval:    10m
```

Created for each resource with `backup.enabled = true` and `monitoring.enabled = true`. The backup workflow pings this monitor after successful completion. The 25-hour interval allows for daily backups with a 1-hour buffer.

**Maintenance Monitors**:

```
Name:               {Project Title} - {App Title} Maintenance
Type:               push
Interval:           8 days (691200s)
Max Retries:        2
Retry Interval:     60s
Resend Interval:    10m
```

Created for each VM app with `monitoring.enabled = true`. The server-maintenance workflow pings this monitor after successful completion. The 8-day interval allows for weekly maintenance with a 1-day buffer.

**Ping URLs**:

Monitor ping URLs are automatically exported as GitHub repository variables:
- Backup monitors: `{ENV}_{RESOURCE_NAME}_BACKUP_PING_URL` (e.g., `PROD_DB_BACKUP_PING_URL`)
- Maintenance monitors: `{ENV}_{APP_NAME}_MAINTENANCE_PING_URL` (e.g., `PROD_WEB_MAINTENANCE_PING_URL`)

CI/CD workflows use the `.github/actions/peekaping-ping` composite action to send heartbeat pings.

## Tags and Organization

Each monitor is automatically tagged with:

- **Project Tag**: Your project title with brand color
- **Environment Tag**: `Production`, `Staging`, etc.
- **WebKit Tag**: Green tag indicating WebKit-managed infrastructure

Tags help organize and filter monitors in the Peekaping UI.

**Note**: The shared "WebKit" and environment tags (Production, Staging, Development, Test) must be created once in Peekaping before deploying any webkit repos.

## Status Pages

WebKit automatically creates a public status page for your project. The status page is accessible via:

1. **Default URL**: `https://peekaping.example.com/status/{project-name}`
2. **Custom Domain**: `status.{your-primary-domain}` (automatically configured)

For example, if your primary app domain is `example.com`, the status page will be configured for `status.example.com`.

### What the Status Page Shows

- Real-time health of **all** HTTP, DNS, and push monitors
- Historical uptime data
- Incident information
- Custom branding (logo and colors from `app.json`)

All monitors are automatically attached to the status page - no manual configuration needed.

### Status Page Branding

Configure status page branding in your `app.json`:

```json
{
  "project": {
    "title": "My Project",
    "brand": {
      "logo_url": "https://example.com/logo.png",
      "primary_color": "#3B82F6"
    }
  }
}
```

## Implementation Details

### Architecture

1. **Schema Layer** (`schema.json`):
   - Simple `MonitoringConfig` type with single `enabled` field

2. **Appdef Layer** (`internal/appdef/monitor.go`):
   - `Monitor` struct
   - Monitor generation logic for apps

3. **Terraform Layer** (`internal/infra/tf_vars.go`):
   - Transforms appdef.Monitor to tfMonitor
   - Generates monitor list for Terraform variables

4. **Infrastructure Layer** (`platform/terraform/modules/monitoring`):
   - Creates Peekaping monitors via tafaust/peekaping provider
   - Creates tags, notifications, and status pages
   - Outputs monitor IDs and details

### Provider

WebKit uses the [tafaust/peekaping](https://registry.terraform.io/providers/tafaust/peekaping) Terraform provider, which supports:

- HTTP, DNS, and Push monitors
- Slack notifications
- Public status pages
- Tags for organization
- Full API support for Peekaping instances

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
monitoring module (applies interval, retries, notifications, tags)
  ↓
Peekaping monitors, status page created
```

## Troubleshooting

### No monitors created

Check that monitoring is enabled (default) and apps have domains:

```bash
# Check generated Terraform variables
cat .webkit/terraform/base/webkit.auto.tfvars.json | jq '.monitors'
```

Monitors are only created for apps with domains. Resource monitors (databases with push heartbeats) are planned for future releases.

### Peekaping authentication fails

Verify credentials in `.env.production.enc`:

```bash
# Decrypt and check
sops -d .env.production.enc | grep PEEKAPING
```

Ensure the credentials match your Peekaping instance admin account.

### Terraform provider errors

The Peekaping provider requires:

1. Peekaping instance is running and accessible at the configured endpoint
2. API key is correct and has proper permissions
3. Provider version is compatible (`~> 0.1.1`)

Check Terraform logs for detailed error messages:

```bash
./webkit infra plan
```

### Slack notifications not working

Verify webhook URL:

```bash
# Test the webhook manually
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"Test from WebKit"}' \
  YOUR_SLACK_WEBHOOK_URL
```

If the manual test works but monitors aren't sending alerts, check the notification configuration in Peekaping UI.

## Custom Monitors

In addition to automatically generated monitors for domains, you can define custom monitors in your `app.json`:

```json
{
  "monitoring": {
    "custom": [
      {
        "name": "External API Check",
        "type": "http",
        "interval": 60,
        "config": {
          "url": "https://api.example.com/health",
          "method": "GET",
          "max_redirects": 5
        }
      },
      {
        "name": "Homepage Content Validation",
        "type": "http-keyword",
        "interval": 300,
        "config": {
          "url": "https://mysite.com",
          "method": "GET",
          "keyword": "Welcome",
          "invert_keyword": false
        }
      },
      {
        "name": "Custom DNS Check",
        "type": "dns",
        "interval": 600,
        "config": {
          "domain": "mail.example.com",
          "resolver_type": "MX"
        }
      }
    ]
  }
}
```

**Available Monitor Types**:
- `http` - HTTP/HTTPS endpoint monitoring
- `http-keyword` - HTTP content keyword monitoring
- `dns` - DNS resolution monitoring
- `postgres` - PostgreSQL connection monitoring
- `push` - Heartbeat/webhook monitoring

### Default Intervals

The `interval` field is **optional** for custom monitors. If not specified, sensible defaults are automatically applied based on the monitor type:

| Monitor Type | Default Interval | Description |
|--------------|------------------|-------------|
| `http` | 60 seconds | Standard HTTP health checks |
| `http-keyword` | 60 seconds | HTTP content validation |
| `postgres` | 60 seconds | Database connection checks |
| `dns` | 300 seconds (5 minutes) | DNS resolution checks |
| `push` | 90000 seconds (25 hours) | Heartbeat monitoring for daily jobs |

**Example with default interval:**
```json
{
  "name": "API Health Check",
  "type": "http",
  "config": {
    "url": "https://api.example.com/health",
    "method": "GET"
  }
}
```
This monitor will automatically use a 60-second interval.

**Note**: All intervals must be at least 20 seconds (Peekaping provider requirement). The default intervals are chosen to match the behavior of automatically-generated monitors and provide appropriate check frequencies for each monitor type.

## Future Enhancements

1. **Resource monitoring**:
   - Postgres database connection health checks
   - Other database types (MySQL, MongoDB, Redis)

2. **Enhanced configuration** in `app.json`:
   - Custom check intervals per app
   - Expected status codes per endpoint
   - Custom heartbeat intervals for backup and maintenance monitors
   - Authentication for health checks
   - Regular expression matching for keyword monitors

3. **Multi-region monitoring**:
   - Check endpoints from multiple geographic locations
   - Regional status pages

4. **Advanced alerting**:
   - Multiple notification channels (email, PagerDuty, etc.)
   - Alert escalation policies
   - Maintenance windows

5. **Performance metrics**:
   - Response time tracking
   - SLA monitoring and reporting
   - Historical performance graphs
