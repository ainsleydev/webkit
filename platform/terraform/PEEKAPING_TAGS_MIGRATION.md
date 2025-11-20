# Peekaping Tags Migration Guide

## Problem

Previously, every webkit repo tried to create "WebKit" and environment tags (e.g., "Production", "Staging") in Peekaping. Since tag names must be unique across the entire Peekaping instance, only the first repo would succeed, and all subsequent repos would fail with:

```
Error: create tag failed
tag with this name already exists
```

## Solution

We've changed the monitoring module to use **data sources** to look up shared tags instead of creating them:

- **Shared tags** (WebKit, Production, Staging, etc.) are looked up via data sources
- **Project-specific tags** are still created as resources (unique per project)

## Migration Steps

### 1. Create Shared Tags (One-Time Setup)

The shared tags must exist in Peekaping before any repos can reference them. You have two options:

#### Option A: Manual Creation in Peekaping UI
1. Log into your Peekaping instance
2. Navigate to Tags
3. Create these tags:
   - Name: `WebKit`, Color: `#10B981` (Green)
   - Name: `Production`, Color: `#EF4444` (Red)
   - Name: `Staging`, Color: `#F59E0B` (Orange)
   - Name: `Development`, Color: `#3B82F6` (Blue)
   - Name: `Test`, Color: `#8B5CF6` (Purple)

#### Option B: Use the Setup Terraform File
```bash
cd platform/terraform

# Configure provider (or use environment variables)
export PEEKAPING_ENDPOINT="https://your-peekaping-instance.com"
export PEEKAPING_API_KEY="your-api-key"

# Create the shared tags
terraform init
terraform apply -target=peekaping_tag.shared_webkit \
                -target=peekaping_tag.shared_production \
                -target=peekaping_tag.shared_staging \
                -target=peekaping_tag.shared_development \
                -target=peekaping_tag.shared_test

# After tags are created, you can delete setup-shared-tags.tf
rm setup-shared-tags.tf
```

### 2. Update Existing Repos

For repos that already have Peekaping monitoring configured:

#### If tags were successfully created before:
```bash
# Remove the old tag resources from state
terraform state rm module.monitoring.peekaping_tag.webkit
terraform state rm module.monitoring.peekaping_tag.environment

# Re-initialize to pick up data sources
terraform init -upgrade
terraform plan

# You should see no changes (data sources will find existing tags)
```

#### If tag creation failed before:
```bash
# Just re-run terraform after shared tags are created
terraform init -upgrade
terraform plan
terraform apply
```

### 3. New Repos

New repos will automatically work once shared tags exist - no additional steps needed!

## What Changed

### Before (`platform/terraform/modules/monitoring/main.tf`)
```hcl
resource "peekaping_tag" "webkit" {
  name        = "WebKit"
  color       = "#10B981"
  description = "Managed by WebKit infrastructure"
}

resource "peekaping_tag" "environment" {
  name        = title(var.environment)
  color       = local.tag_color
  description = "${title(var.environment)} environment"
}

resource "peekaping_tag" "project" {
  name        = var.project_title
  color       = local.tag_color
  description = "Monitor for ${var.project_title}"
}
```

### After
```hcl
# Shared tags - looked up via data source
data "peekaping_tag" "webkit" {
  name = "WebKit"
}

data "peekaping_tag" "environment" {
  name = title(var.environment)
}

# Project-specific tag - still created as resource
resource "peekaping_tag" "project" {
  name        = var.project_title
  color       = local.tag_color
  description = "Monitor for ${var.project_title}"
}
```

## Troubleshooting

### Error: "not found: no tag matched the given criteria"
This means the shared tags haven't been created yet. Follow Step 1 above to create them.

### Error: "tag with this name already exists"
If you still see this error for the `project` tag, it means two repos have the same `project_title`. Ensure each repo has a unique project title.

## Technical Details

- **Database constraint**: Peekaping enforces a UNIQUE constraint on tag names at the database level
- **API validation**: The Peekaping API checks for duplicate names before creating tags
- **Data source lookup**: Uses case-insensitive name matching to find existing tags
- **Multiple state files**: Each webkit repo has its own Terraform state, so they all need to reference shared resources via data sources

## References

- [Peekaping Tag Resource Documentation](https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/tag)
- [Peekaping Tag Data Source Documentation](https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/data-sources/tag)
- [Peekaping Database Schema](https://github.com/0xfurai/peekaping/blob/main/apps/server/cmd/bun/migrations/20250712111431_add_tags_system.tx.up.sql)
