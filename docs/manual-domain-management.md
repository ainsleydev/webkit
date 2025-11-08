# Manual Domain Management

This Terraform configuration supports manual domain management in DigitalOcean projects while still managing other infrastructure resources (apps, databases, buckets) through Terraform.

## Problem

When you manually add domains to a DigitalOcean project through the DO UI, Terraform would normally try to remove them on the next apply since they're not defined in the Terraform configuration. This creates a conflict between infrastructure-as-code and manual domain management.

## Solution

The configuration uses an external data source to query the DigitalOcean API and fetch all manually-added domain URNs from the project. These domain URNs are then merged with Terraform-managed resources (apps, databases, buckets) when updating the project.

### How It Works

1. **External Script**: `platform/terraform/base/scripts/get_project_domains.sh` queries the DO API to fetch domain URNs
2. **Data Source**: The `external` data source in `main.tf` calls the script
3. **Merge Logic**: Domain URNs are merged with Terraform-managed resource URNs
4. **Project Update**: The `digitalocean_project` resource receives the combined list

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

### First Apply (New Project)

1. Run your first Terraform apply without setting `digitalocean_project_id`:

```bash
terraform apply
```

2. After successful apply, note the `digitalocean_project_id` output:

```bash
terraform output digitalocean_project_id
# Output: 1f726d4d-1d77-4ee0-a4b6-a1a66720209a
```

3. Set the project ID for future applies using one of these methods:

   **Option A: Environment Variable (Recommended)**
   ```bash
   export TF_VAR_digitalocean_project_id="1f726d4d-1d77-4ee0-a4b6-a1a66720209a"
   ```

   **Option B: terraform.tfvars file**
   ```hcl
   digitalocean_project_id = "1f726d4d-1d77-4ee0-a4b6-a1a66720209a"
   ```

   **Option C: Command line**
   ```bash
   terraform apply -var="digitalocean_project_id=1f726d4d-1d77-4ee0-a4b6-a1a66720209a"
   ```

### Subsequent Applies

Once the project ID is set, you can freely:
- ✅ Add domains to the project manually in DigitalOcean UI
- ✅ Remove domains from the project manually
- ✅ Add new Terraform-managed resources (databases, apps, etc.)
- ✅ Terraform will preserve your manual domains while managing its own resources

### Managing Domains

**To add a domain manually:**
1. Go to DigitalOcean UI → Projects → Your Project
2. Click "Add Resources"
3. Select your domain(s)
4. Run `terraform apply` - your domains will be preserved

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

1. Receives `project_id` and `do_token` as JSON input from Terraform
2. Queries the DO API: `GET /v2/projects/{project_id}/resources`
3. Filters resources for domain URNs (format: `do:domain:example.com`)
4. Returns comma-separated domain URNs to Terraform

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
1. Verify `digitalocean_project_id` is set correctly
2. Check that the script is executable
3. Verify `jq` is installed
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
