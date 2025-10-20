#!/usr/bin/env bash
set -e

echo "🔧 Setting up development environment..."

# Tools defined in HomeBrew
declare -A tools=(
  [rg]="ripgrep"
  [jq]="jq"
  [sops]="sops"
  [age]="age"
  [repomix]="repomix"
  [terraform]="terraform"
)

# Install each one
for cmd in "${!tools[@]}"; do
  formula=${tools[$cmd]}
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "⚠️  $cmd not found. Installing via Homebrew..."
    if ! command -v brew >/dev/null 2>&1; then
      echo "❌ Homebrew not found. Please install Homebrew first: https://brew.sh/"
      exit 1
    fi
    brew install "$formula"
  else
    echo "✅ $cmd already installed."
  fi
done

# Install npm packages globally (or save-dev)
echo "📦 Installing npm packages..."
if ! command -v npm >/dev/null 2>&1; then
  echo "❌ npm not found. Please install Node.js first: https://nodejs.org/"
  exit 1
fi

echo "✅ Setup complete!"
