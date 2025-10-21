setup: # Install all dev tools (Homebrew packages + Go tools + action-validator)

	# Brew
	@echo "Setting up development environment..."
	@if ! command -v brew >/dev/null 2>&1; then \
		echo "Homebrew not found. Installing Homebrew..."; \
		/bin/bash -c "$$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"; \
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
		echo "Installing via npm..."; \
		npm install -g @action-validator/core @action-validator/cli; \
	else \
		echo "npm not found. Please install Node.js first: https://nodejs.org/"; \
		exit 1; \
	fi

	@echo "Setup complete."
.PHONY: setup

help: # Show available commands
	@echo ""
	@echo "Available make commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  %-25s %s\n", $$1, $$2}'
	@echo ""
.PHONY: help
