# Makefile for Terraform scripts

# Default .tfvars file
TFVARS ?= .tfvars
BACKEND := backend.hcl

setup: # Setup's local machine
	brew install terraform
	brew install tflint
	brew install jq
.PHONY: setup

fmt: # Format Terraform files
	terraform fmt -recursive
.PHONY: fmt

lint: # Lint Terraform files
	tflint
.PHONY: lint

init: # Initialize Terraform with backend
	terraform init -backend-config=$(BACKEND)
.PHONY: lint

init-local: # Initialize Terraform without backend (local state)
	terraform init
.PHONY: init-local

init-upgrade: # Initialize and upgrade Terraform modules/providers
	terraform init -backend-config=$(BACKEND) -upgrade
.PHONY: init-upgrade

plan: # Plan Terraform changes
	terraform plan -var-file=$(TFVARS)
.PHONY: plan

apply: # Apply Terraform changes
	terraform apply -var-file=$(TFVARS)
.PHONY: apply

refresh: # Refresh Terraform state
	terraform refresh -var-file=$(TFVARS)
.PHONY: refresh

destroy: # Destroy Terraform-managed infrastructure
	terraform destroy -var-file=$(TFVARS)
.PHONY: destroy

# Example: DO_API_KEY="dop_..." PROJECT_ID="1f726d4d-1d77-4ee0-a4b6-a1a66720209a" make test-domains
test-domains: # Test domain fetching from DigitalOcean project (requires DO_API_KEY and PROJECT_ID env vars)
	@echo '{"project_id":"$(PROJECT_ID)","do_token":"$(DO_API_KEY)"}' | ./base/scripts/get_project_domains.sh | jq -r '.domain_urns' | tr ',' '\n' | sed 's/do:domain://'
.PHONY: test-domains

pack: # Packs up the folder for AI.
	npx repopack --ignore "**/*.log,tmp/,**/fixtures/**,.env,.tfvars,**/.terraform/**,.terraform.lock.hcl" --output Terraform.txt
.PHONY: pack

todo: # Show to-do items per file
	$(Q) grep \
		--exclude=Makefile.util \
		--exclude-dir=vendor \
		--exclude-dir=.vercel \
		--exclude-dir=.gen \
		--exclude-dir=.idea \
		--exclude-dir=public \
		--exclude-dir=node_modules \
		--exclude-dir=archetypes \
		--exclude-dir=.git \
		--text \
		--color \
		-nRo \
		-E '\S*[^\.]TODO.*' \
		.
.PHONY: todo

help: # Display this help
	$(Q) awk 'BEGIN {FS = ":.*#"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?#/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help
