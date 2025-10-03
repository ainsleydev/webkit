# Terraform

Introduction and details needed.

## State Management

WebKit uses Backblaze B2 for Terraform state storage.

### Initial Setup

Before running `webkit infra apply`, configure your state backend:
```bash
webkit init --configure-state