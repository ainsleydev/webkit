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

## Slack Infra

Send Slack messages on deployment:

**Example**:

```yaml
  - name: Notify Slack of Web Deploy Success
	if: success()
	uses: ./.github/actions/slack-notify
	with:
	  title: 'Web App Deployment Successful'
	  message: 'The web application has been successfully deployed to production and is now live!'
	  status: 'success'
	  url: https://searchspares.com
	  commit_sha: ${{ github.sha }}
	  slack_bot_token: ${{ secrets.SLACK_BOT_TOKEN }}
	  channel_id: ${{ secrets.SLACK_CHANNEL_ID }}

  - name: Notify Slack of Web Deploy Failure
	if: failure()
	uses: ./.github/actions/slack-notify
	with:
	  title: 'Web App Deployment Failed'
	  message: 'The web application deployment to production has failed. Please check the logs for details.'
	  status: 'failure'
	  commit_sha: ${{ github.sha }}
	  slack_bot_token: ${{ secrets.SLACK_BOT_TOKEN }}
	  channel_id: ${{ secrets.SLACK_CHANNEL_ID }}
```

**Action**: 

```yaml
name: 'Slack Notify'
description: 'Send a rich message to Slack channel using blocks'
inputs:
  slack_bot_token:
    description: 'Slack Bot Token'
    required: true
  channel_id:
    description: 'Channel to send to'
    required: true
  title:
    description: 'Message title'
    required: true
  message:
    description: 'Message content'
    required: true
  status:
    description: 'Status type (success, failure, info)'
    required: false
    default: 'info'
  url:
    description: 'Optional URL to link to'
    required: false
  commit_sha:
    description: 'Optional commit SHA'
    required: false
runs:
  using: 'composite'
  steps:
    - name: Send rich message to Slack
      shell: bash
      run: |
        # Determine color and emoji based on status
        case "${{ inputs.status }}" in
          "success")
            COLOR="good"
            EMOJI=":white_check_mark:"
            ;;
          "failure")
            COLOR="danger"
            EMOJI=":x:"
            ;;
          *)
            COLOR="#36a64f"
            EMOJI=":information_source:"
            ;;
        esac

        # Build the blocks JSON
        BLOCKS='[
          {
            "type": "section",
            "text": {
              "type": "mrkdwn",
              "text": "'$EMOJI' *${{ inputs.title }}*"
            }
          },
          {
            "type": "section",
            "text": {
              "type": "mrkdwn",
              "text": "${{ inputs.message }}"
            }
          }'

        # Add URL field if provided
        if [ -n "${{ inputs.url }}" ]; then
          BLOCKS+=',
          {
            "type": "section",
            "text": {
              "type": "mrkdwn",
              "text": ":link: <${{ inputs.url }}|View Deployment>"
            }
          }'
        fi

        # Add commit info if provided
        if [ -n "${{ inputs.commit_sha }}" ]; then
          SHORT_SHA=$(echo "${{ inputs.commit_sha }}" | cut -c1-7)
          BLOCKS+=',
          {
            "type": "context",
            "elements": [
              {
                "type": "mrkdwn",
                "text": ":git-commit: Commit: `'$SHORT_SHA'` | :calendar: '$(date -u +"%Y-%m-%d %H:%M UTC")'"
              }
            ]
          }'
        fi

        BLOCKS+=']'

        # Create the payload
        PAYLOAD=$(cat <<EOF
        {
          "channel": "${{ inputs.channel_id }}",
          "blocks": $BLOCKS,
          "attachments": [
            {
              "color": "$COLOR",
              "blocks": []
            }
          ]
        }
        EOF
        )

        # Send to Slack
        curl -X POST -H "Authorization: Bearer ${{ inputs.slack_bot_token }}" \
             -H "Content-type: application/json" \
             --data "$PAYLOAD" \
             https://slack.com/api/chat.postMessage

```
