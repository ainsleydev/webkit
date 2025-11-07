# Slack integration

WebKit provides comprehensive Slack integration for CI/CD notifications, allowing you to monitor deployments, backups, maintenance tasks, and infrastructure alerts in one centralised location.

## Overview

WebKit integrates with Slack at two levels:

1. **CI/CD notifications** - Deployment status, backup failures, and maintenance alerts via GitHub Actions.
2. **Infrastructure alerts** - DigitalOcean App Platform service-level alerts (CPU, memory, restarts).

## Prerequisites

Before setting up Slack integration, ensure you have:

- A Slack workspace with admin access.
- A Slack app (bot) installed in your workspace.
- Organisation-level GitHub secrets configured:
  - `ORG_SLACK_BOT_TOKEN` - Bot User OAuth token (starts with `xoxb-`).
  - `ORG_SLACK_USER_TOKEN` - User OAuth token (starts with `xoxp-`).

## Setup

### Step 1: Configure environment variables

WebKit requires Slack tokens to be available as environment variables when running Terraform operations:

```bash
export SLACK_BOT_TOKEN="xoxb-your-bot-token"
export SLACK_USER_TOKEN="xoxp-your-user-token"
```

These tokens are used to:
- Create a dedicated Slack channel for your project alerts.
- Store the channel ID in GitHub secrets automatically.

### Step 2: Run Terraform

When you run `webkit infra apply`, Terraform will automatically:

1. Create a Slack channel named `{project-name}-alerts`.
2. Configure the channel topic: "CI/CD alerts and notifications for {project-title}".
3. Add you as a permanent member of the channel.
4. Store the channel ID as `TF_SLACK_CHANNEL_ID` in your GitHub repository secrets.

```bash
webkit infra apply
```

**Output example:**
```
slack_channel_id = "C1234567890"
slack_channel_name = "my-project-alerts"
```

### Step 3: Verify GitHub secrets

After Terraform completes, verify the channel ID was stored:

```bash
# Using GitHub CLI
gh secret list

# You should see:
# TF_SLACK_CHANNEL_ID  Updated 2025-11-07
```

### Step 4: Set up infrastructure alerts (optional)

To receive DigitalOcean App Platform alerts in Slack, you need to create an incoming webhook manually.

#### Why manual?

DigitalOcean requires a webhook URL for alerts, and Slack's incoming webhooks cannot be created programmatically via Terraform. This is a one-time setup per project.

#### Steps:

1. **Navigate to your Slack app settings:**
   - Go to [api.slack.com/apps](https://api.slack.com/apps).
   - Select your app (e.g., "ainsley.dev").

2. **Enable incoming webhooks:**
   - Click "Incoming Webhooks" in the sidebar.
   - Turn on "Activate Incoming Webhooks".

3. **Create a new webhook:**
   - Scroll down and click "Add New Webhook to Workspace".
   - Select the channel created by Terraform (e.g., `my-project-alerts`).
   - Authorise the webhook.

4. **Copy the webhook URL:**
   ```
   https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXX
   ```

5. **Set as environment variable:**
   ```bash
   export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/..."
   ```

6. **Run Terraform again:**
   ```bash
   webkit infra apply
   ```

Terraform will now configure DigitalOcean App Platform to send alerts to your Slack channel.

## What gets notified?

### CI/CD notifications (automatic)

Once Terraform creates the channel and stores the ID in GitHub secrets, the following notifications are automatically sent:

#### Backup failures
- **Database backup failures** - Per resource (Postgres databases).
- **Codebase backup failures** - GitHub repository backups.
- **Google Drive sync failures** - Backup synchronisation to Google Drive.

**Example message:**
```
❌ Database Backup Failed

Status: ❌ Failed
Triggered By: github-actions

Details:
The scheduled database backup has failed. This may impact disaster recovery capabilities.

[View Logs] [View Repository]
```

#### Server maintenance failures
- **VM maintenance failures** - Weekly Ansible maintenance playbook failures.

**Example message:**
```
❌ Server Maintenance Failed - Production Web

Status: ❌ Failed
Triggered By: github-actions

Details:
Weekly maintenance tasks failed for Production Web. Server updates may not have been applied.

[View Ansible Logs]
```

#### Release notifications
- **Successful deployments** - When all applications deploy successfully.
- **Failed deployments** - When any step in the release pipeline fails.

**Success example:**
```
✅ Release Successful

All applications have been successfully deployed to production.

Deployed Apps:
- Web Application
- API Service

[View Deployment] [View Commit]
```

**Failure example:**
```
❌ Release Failed

Status: ❌ Failed
Triggered By: ainsley

Impact:
One or more steps in the release pipeline failed. Applications may not be running the latest code.

[View Failure Logs] [Re-run Workflow]
```

### Infrastructure alerts (requires webhook setup)

After configuring the incoming webhook, DigitalOcean App Platform will send the following alerts:

#### CPU utilisation
Triggered when CPU usage exceeds 80% for 5 consecutive minutes.

#### Memory utilisation
Triggered when memory usage exceeds 80% for 5 consecutive minutes.

#### Restart count
Triggered when the application restarts more than 3 times in a 5-minute window.

**Note:** These alerts remain in place until Prometheus monitoring is implemented.

## Message format

All notifications use Slack's Block Kit format with rich formatting:

- **Colour-coded attachments** - Green (success), red (failure), yellow (warning), blue (info).
- **Structured fields** - Status, triggered by, commit information.
- **Action buttons** - Direct links to GitHub workflows, commits, and repositories.
- **Contextual information** - Timestamps, workflow names, and relevant metadata.

**No emojis in titles** - Emojis are added to message bodies based on status, but titles remain clean.

## Configuration reference

### Environment variables

| Variable | Required | Description |
|----------|----------|-------------|
| `SLACK_BOT_TOKEN` | Yes | Bot User OAuth token for Terraform operations |
| `SLACK_USER_TOKEN` | Yes | User OAuth token for Terraform operations |
| `SLACK_WEBHOOK_URL` | No | Incoming webhook URL for DO app alerts |

### GitHub secrets (automatically created)

| Secret | Created By | Description |
|--------|------------|-------------|
| `TF_SLACK_CHANNEL_ID` | Terraform | Channel ID for CI/CD notifications |
| `ORG_SLACK_BOT_TOKEN` | Manual | Bot token (organisation-level) |

### Terraform outputs

| Output | Description |
|--------|-------------|
| `slack_channel_id` | Channel ID for programmatic access |
| `slack_channel_name` | Channel name (e.g., "my-project-alerts") |

## Troubleshooting

### Channel not created

**Issue:** Terraform fails with "slack: authentication failed".

**Solution:** Verify your bot token has the correct permissions:
- `channels:manage` - Create and manage channels.
- `channels:read` - Read channel information.
- `chat:write` - Post messages to channels.

### Notifications not appearing

**Issue:** Workflows run successfully, but no Slack messages appear.

**Solution:**
1. Verify `TF_SLACK_CHANNEL_ID` exists in GitHub secrets.
2. Check the channel ID is correct (starts with 'C').
3. Ensure the bot is a member of the channel.
4. Check GitHub Actions logs for Slack API errors.

### Infrastructure alerts not working

**Issue:** DigitalOcean alerts don't appear in Slack.

**Solution:**
1. Verify webhook URL is set as `SLACK_WEBHOOK_URL` environment variable.
2. Run `webkit infra apply` after setting the webhook URL.
3. Check DigitalOcean App Platform settings to confirm webhook is configured.
4. Test the webhook URL manually:
   ```bash
   curl -X POST $SLACK_WEBHOOK_URL \
     -H "Content-Type: application/json" \
     -d '{"text":"Test message"}'
   ```

### Webhook URL not being used

**Issue:** Terraform doesn't configure DigitalOcean alerts despite webhook URL being set.

**Solution:**
1. Ensure `SLACK_WEBHOOK_URL` environment variable is exported before running Terraform.
2. Verify the variable is being passed through Terraform layers:
   - Base → Apps module → DigitalOcean app provider.
3. Check Terraform output for `alert_destination` block in DigitalOcean app resource.

## Best practices

1. **Use organisation-level secrets** - Store `ORG_SLACK_BOT_TOKEN` at the organisation level to share across projects.

2. **One channel per project** - Each project gets its own dedicated alerts channel for better organisation.

3. **Archive on destroy** - When running `terraform destroy`, channels are archived (not deleted) to preserve message history.

4. **Test notifications** - After setup, manually trigger a workflow to verify notifications work correctly.

5. **Monitor rate limits** - Slack has rate limits (1 message/second per channel). WebKit workflows are designed to stay within these limits.

6. **Regular webhook rotation** - For security, rotate webhook URLs periodically and update the environment variable.

## Further reading

- [Slack Block Kit Builder](https://app.slack.com/block-kit-builder/) - Design and preview message layouts.
- [Slack API documentation](https://api.slack.com/) - Complete API reference.
- [DigitalOcean App Platform alerts](https://docs.digitalocean.com/products/app-platform/how-to/manage-alerts/) - Alert configuration and troubleshooting.
