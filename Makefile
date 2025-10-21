setup-mac: # Install dev tools on macOS
	@echo "🔧 Setting up development environment for macOS..."
	@if ! command -v brew >/dev/null 2>&1; then \
		echo "❌ Homebrew not found. Please install from https://brew.sh/"; \
		exit 1; \
	fi
	@brew install ripgrep jq sops age terraform || true
	@$(MAKE) install-action-validator
	@echo "✅ macOS setup complete!"
.PHONY: setup-mac

setup-linux: # Install dev tools on Linux (Ubuntu/Debian)
	@echo "🔧 Setting up development environment for Linux..."
	@if command -v apt-get >/dev/null 2>&1; then \
		sudo apt-get update -y && \
		sudo apt-get install -y ripgrep jq sops age terraform; \
	elif command -v dnf >/dev/null 2>&1; then \
		sudo dnf install -y ripgrep jq sops age terraform; \
	else \
		echo "❌ No supported package manager found (need apt or dnf)"; \
		exit 1; \
	fi
	@$(MAKE) install-action-validator
	@echo "✅ Linux setup complete!"
.PHONY: setup-linux

install-action-validator: # Install action-validator via npm or cargo
	@echo "🧩 Installing action-validator..."
	@if command -v action-validator >/dev/null 2>&1; then \
		echo "✅ action-validator already installed."; \
	elif command -v npm >/dev/null 2>&1; then \
		echo "📦 Installing via npm..."; \
		npm install -g @action-validator/core @action-validator/cli; \
	elif command -v cargo >/dev/null 2>&1; then \
		echo "⚙️ Installing via cargo..."; \
		cargo install action-validator; \
	else \
		echo "❌ Neither npm nor cargo found. Please install Node.js or Rust first."; \
		exit 1; \
	fi
.PHONY: install-action-validator

install-go-tools: ## Install Go-based tools used in the project
	@echo "🐹 Installing Go tools..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "❌ Go not found. Please install Go first: https://go.dev/dl/"; \
		exit 1; \
	fi
	go install go.uber.org/mock/mockgen@latest
	@echo "✅ Go tools installed!"
.PHONY: install-go-tools

help: # Show available commands
	@echo ""
	@echo "Available make commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""

.PHONY: help
