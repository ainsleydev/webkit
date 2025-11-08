# Manual Domain Management

This Terraform configuration supports manual domain management in DigitalOcean projects while still managing other infrastructure resources (apps, databases, buckets) through Terraform.

## Problem

When you manually add domains to a DigitalOcean project through the DO UI, Terraform would normally try to remove them on the next apply since they're not defined in the Terraform configuration. This creates a conflict between infrastructure-as-code and manual domain management.

## Solution

The configuration uses an external data source to query the DigitalOcean API and fetch all manually-added domain URNs from the project. These domain URNs are then merged with Terraform-managed resources (apps, databases, buckets) when updating the project.

### How It Works

1. **Project Lookup**: The script receives the `project_title` from Terraform variables
2. **Find Project ID**: If no `project_id` is provided, the script queries `/v2/projects` to find the project by title/name
3. **Fetch Domains**: Using the project ID, the script queries `/v2/projects/{id}/resources` to get domain URNs
4. **Merge Logic**: Domain URNs are merged with Terraform-managed resource URNs in `main.tf`
5. **Project Update**: The `digitalocean_project` resource receives the combined list

This approach eliminates the need to manually set `digitalocean_project_id` - the script automatically finds the project by its title.

## Testing

Before using this in Terraform, you can test the domain fetching functionality:

```bash
cd platform/terraform
DO_API_KEY="your-api-key" PROJECT_ID="1f726d4d-1d77-4ee0-a4b6-a1a66720209a" make test-domains
```

Example output:
```
searchspares.co.uk
searchspares.com
```

## Setup Instructions

No special setup required! The script automatically finds your project by its title (from `project_title` variable).

**Important:** The DigitalOcean project must be created first (via initial `webkit infra apply`) before you can manually add domains to it.

### Workflow

1. **First Apply**: Run `webkit infra apply` to create the DigitalOcean project and infrastructure
2. **Add Domains**: Once the project exists, you can manually add domains via the DigitalOcean UI
3. **Subsequent Applies**: Run `webkit infra plan/apply` - your manually-added domains will be preserved

### Managing Domains

**To add a domain manually (after initial terraform apply):**
1. Ensure the DigitalOcean project has been created by Terraform
2. Go to DigitalOcean UI → Projects → Your Project
3. Click "Add Resources"
4. Select your domain(s)
5. Run `webkit infra apply` - your domains will be preserved

**To remove a domain manually:**
1. Go to DigitalOcean UI → Projects → Your Project
2. Remove the domain
3. Run `terraform apply` - the removal will be preserved

## Requirements

- **jq**: The script requires `jq` for JSON processing
- **curl**: Used to query the DigitalOcean API
- **bash**: The script is written in bash

Install dependencies:
```bash
# macOS
brew install jq curl

# Ubuntu/Debian
apt-get install jq curl

# CentOS/RHEL
yum install jq curl
```

## How the Script Works

The `get_project_domains.sh` script:

1. Receives `project_id` (optional), `project_title`, and `do_token` as JSON input from Terraform
2. If `project_id` is empty, queries `GET /v2/projects` to find the project by title
3. Queries the DO API: `GET /v2/projects/{project_id}/resources` to get all project resources
4. Filters resources for domain URNs (format: `do:domain:example.com`)
5. Returns comma-separated domain URNs to Terraform

## Troubleshooting

### Testing the Script

Use the test command to verify everything works:

```bash
cd platform/terraform
DO_API_KEY="your-api-key" PROJECT_ID="1f726d4d-..." make test-domains
```

If the test fails, check:
- DO_API_KEY is set correctly
- PROJECT_ID is the correct UUID (get it from `terraform output digitalocean_project_id`)
- jq and curl are installed
- API token has project read permissions

### Script Execution Error

If you see errors about the script not being executable:
```bash
chmod +x platform/terraform/base/scripts/get_project_domains.sh
```

### jq Not Found

Install `jq`:
```bash
# macOS
brew install jq

# Ubuntu/Debian
apt-get install jq

# CentOS/RHEL
yum install jq
```

### Domains Being Removed

If domains are still being removed:
1. Verify your `project_title` variable matches the exact project name in DigitalOcean
2. Check that the script is executable (`chmod +x platform/terraform/base/scripts/get_project_domains.sh`)
3. Verify `jq` is installed (`brew install jq` on macOS)
4. Check DO API token has project read permissions
5. Run the test script to verify domain fetching works

### API Permission Errors

Ensure your DigitalOcean API token has:
- Read access to projects
- Read access to project resources

You can test this with:
```bash
curl -X GET \
  -H "Authorization: Bearer $DO_API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.digitalocean.com/v2/projects"
```

## Benefits

- ✅ **Flexible**: Manage domains manually while Terraform manages infrastructure
- ✅ **No Conflicts**: Terraform won't remove manually-added domains
- ✅ **Infrastructure as Code**: Apps, databases, and buckets still managed by Terraform
- ✅ **Simple**: One-time setup after first apply
- ✅ **API-Driven**: Uses official DigitalOcean API for reliability
- ✅ **Testable**: Test script validates functionality before Terraform apply

## Limitations

- Requires `jq`, `curl`, and `bash` in the execution environment
- Adds a small API call overhead on each Terraform plan/apply
- First apply must complete before manual domain management works
- Project ID must be set manually after first apply

## Files

- `platform/terraform/base/scripts/get_project_domains.sh` - Main script called by Terraform
- `platform/terraform/base/main.tf` - Terraform configuration with external data source
- `platform/terraform/base/variables.tf` - Variable definition for `digitalocean_project_id`
- `platform/terraform/base/outputs.tf` - Output for project ID after first apply
- `Makefile` - Contains `test-domains` target for testing domain fetching
