#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

ANSIBLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEBKIT_REPO_PATH=""
CHECK_MODE=false

show_help() {
    cat << EOF
Usage: ./test-local-simple.sh <webkit-repo-path> [OPTIONS]

Test Ansible playbooks locally using ansible-playbook --check mode.
This validates syntax, variable usage, and logic without actually changing anything.

ARGUMENTS:
    webkit-repo-path    Path to your webkit-enabled repository

OPTIONS:
    --check             Run in check mode (dry-run, no changes)
    -h, --help          Show this help message

EXAMPLES:
    # Syntax check only (safe, no changes)
    ./test-local-simple.sh ~/projects/playground

    # Full validation
    ./test-local-simple.sh ~/projects/playground --check

REQUIREMENTS:
    brew install ansible

WHAT THIS DOES:
    - Validates playbook syntax
    - Checks variable usage from your app.json
    - Tests template rendering
    - Shows what would be executed
    - Does NOT actually install packages or change system

EOF
}

# Parse arguments
if [ $# -eq 0 ]; then
    echo -e "${RED}Error: webkit-repo-path is required${NC}"
    show_help
    exit 1
fi

if [[ "$1" != -* ]]; then
    WEBKIT_REPO_PATH="$1"
    shift
else
    echo -e "${RED}Error: webkit-repo-path is required${NC}"
    show_help
    exit 1
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        --check)
            CHECK_MODE=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Validate
if [ ! -d "$WEBKIT_REPO_PATH" ]; then
    echo -e "${RED}Error: Directory not found: $WEBKIT_REPO_PATH${NC}"
    exit 1
fi

WEBKIT_REPO_PATH="$(cd "$WEBKIT_REPO_PATH" && pwd)"

if [ ! -f "$WEBKIT_REPO_PATH/app.json" ]; then
    echo -e "${RED}Error: app.json not found${NC}"
    exit 1
fi

# Check for ansible
if ! command -v ansible-playbook &> /dev/null; then
    echo -e "${RED}Error: ansible-playbook not found${NC}"
    echo -e "${YELLOW}Install with: brew install ansible${NC}"
    exit 1
fi

echo -e "${GREEN}Using webkit repo: $WEBKIT_REPO_PATH${NC}"

# Extract config from app.json
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq not found (needed to parse app.json)${NC}"
    echo -e "${YELLOW}Install with: brew install jq${NC}"
    exit 1
fi

DOMAIN=$(jq -r '.apps[0].primary_domain // "test.local"' "$WEBKIT_REPO_PATH/app.json")
APP_NAME=$(jq -r '.apps[0].name // "testapp"' "$WEBKIT_REPO_PATH/app.json")
DOCKER_PORT=$(jq -r '.apps[0].build.port // 3000' "$WEBKIT_REPO_PATH/app.json")
GITHUB_USER=$(jq -r '.repository.owner // "testuser"' "$WEBKIT_REPO_PATH/app.json")
REPO_NAME=$(jq -r '.repository.name // "testrepo"' "$WEBKIT_REPO_PATH/app.json")
ADMIN_EMAIL=$(jq -r '.apps[0].infra.config.admin_email // "test@example.com"' "$WEBKIT_REPO_PATH/app.json")

echo -e "${GREEN}Configuration:${NC}"
echo -e "  Domain: $DOMAIN"
echo -e "  App: $APP_NAME"
echo -e "  Port: $DOCKER_PORT"
echo -e "  GitHub: $GITHUB_USER/$REPO_NAME"
echo ""

# Run ansible-playbook
CMD="ansible-playbook $ANSIBLE_DIR/playbooks/server.yaml"
CMD="$CMD -i localhost,"
CMD="$CMD --connection=local"
CMD="$CMD -e domain=$DOMAIN"
CMD="$CMD -e github_user=$GITHUB_USER"
CMD="$CMD -e docker_image=${REPO_NAME}-${APP_NAME}"
CMD="$CMD -e docker_image_tag=test-latest"
CMD="$CMD -e docker_port=$DOCKER_PORT"
CMD="$CMD -e app_name=$APP_NAME"
CMD="$CMD -e env_name=development"
CMD="$CMD -e age_secret_key=fake-key"
CMD="$CMD -e github_token=fake-token"
CMD="$CMD -e admin_email=$ADMIN_EMAIL"
CMD="$CMD -e enable_https=false"
CMD="$CMD -e skip_reboot=true"
CMD="$CMD -e webkit_config_dir=$WEBKIT_REPO_PATH"

if [ "$CHECK_MODE" = true ]; then
    echo -e "${BLUE}Running in CHECK mode (no changes will be made)${NC}"
    CMD="$CMD --check"
else
    echo -e "${BLUE}Running SYNTAX CHECK only${NC}"
    CMD="$CMD --syntax-check"
fi

echo -e "${YELLOW}Command: $CMD${NC}"
echo ""

if eval "$CMD"; then
    echo ""
    echo -e "${GREEN}✓ Validation successful!${NC}"
    exit 0
else
    echo ""
    echo -e "${RED}✗ Validation failed!${NC}"
    exit 1
fi
