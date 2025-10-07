# Makefile for Terraform scripts

# Default .tfvars file
TFVARS := .tfvars
BACKEND := backend.hcl

# Setup's local machine
setup:
	brew install tflint

# Lint Terraform files
lint:
	tflint

# Initialize Terraform with backend
init:
	terraform init -backend-config=$(BACKEND)

# Initialize and upgrade Terraform modules/providers
init-upgrade:
	terraform init -backend-config=$(BACKEND) -upgrade

# Plan Terraform changes
plan:
	terraform plan -var-file=$(TFVARS)

# Apply Terraform changes
apply:
	terraform apply -var-file=$(TFVARS)

# Refresh Terraform state
refresh:
	terraform refresh -var-file=$(TFVARS)

# Destroy Terraform-managed infrastructure
destroy:
	terraform destroy -var-file=$(TFVARS)

