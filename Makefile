setup: # Install all dev tools (Homebrew packages + Go tools + action-validator)
	@echo "Setting up development environment..."

	# Brew
	@if ! command -v brew >/dev/null 2>&1; then \
		echo "Homebrew not found."; \
		exit 1; \
	fi
	@echo "Installing Brew packages..."
	@brew install ripgrep jq sops age terraform || true

	# Go Dependencies
	@echo "Installing Go tools..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "Go not found. Please install Go first: https://go.dev/dl/"; \
		exit 1; \
	fi
	go install go.uber.org/mock/mockgen@latest
	@echo "Go tools installed."

	# Action Validator
	@if command -v action-validator >/dev/null 2>&1; then \
		echo "action-validator already installed."; \
	elif command -v npm >/dev/null 2>&1; then \
		npm install -g @action-validator/core @action-validator/cli; \
	else \
		echo "npm not found. Please install Node.js first: https://nodejs.org/"; \
		exit 1; \
	fi

	@echo "Setup complete."
.PHONY: setup

test-domains: ## Test domain fetching from DigitalOcean project (requires DO_API_KEY and PROJECT_NAME env vars)
	@if [ -z "$$DO_API_KEY" ]; then \
		echo "Error: DO_API_KEY not set"; \
		echo "Usage: DO_API_KEY=\"...\" PROJECT_NAME=\"Search Spares\" make test-domains"; \
		exit 1; \
	fi
	@if [ -z "$$PROJECT_NAME" ]; then \
		echo "Error: PROJECT_NAME not set"; \
		echo "Usage: DO_API_KEY=\"...\" PROJECT_NAME=\"Search Spares\" make test-domains"; \
		exit 1; \
	fi
	@PROJECT_ID=$$(curl -s -X GET \
		-H "Authorization: Bearer $$DO_API_KEY" \
		"https://api.digitalocean.com/v2/projects" | \
		jq -r --arg name "$$PROJECT_NAME" '.projects[] | select(.name == $$name) | .id'); \
	if [ -z "$$PROJECT_ID" ] || [ "$$PROJECT_ID" = "null" ]; then \
		echo "Error: Project '$$PROJECT_NAME' not found"; \
		exit 1; \
	fi; \
	RESULT=$$(echo "{\"project_id\":\"$$PROJECT_ID\",\"do_token\":\"$$DO_API_KEY\"}" | \
		./platform/terraform/base/scripts/get_project_domains.sh); \
	DOMAIN_URNS=$$(echo "$$RESULT" | jq -r '.domain_urns'); \
	if [ -z "$$DOMAIN_URNS" ] || [ "$$DOMAIN_URNS" = "" ]; then \
		echo "No domains found"; \
	else \
		echo "$$DOMAIN_URNS" | tr ',' '\n' | sed 's/do:domain://'; \
	fi
.PHONY: test-domains

help: # Show available commands
	@echo ""
	@echo "Available make commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  %-25s %s\n", $$1, $$2}'
	@echo ""
.PHONY: help
