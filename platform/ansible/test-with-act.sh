#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if act is installed
if ! command -v act &> /dev/null; then
    echo -e "${RED}Error: act is not installed${NC}"
    echo -e "${YELLOW}Install with: brew install act${NC}"
    echo -e "${YELLOW}Or visit: https://github.com/nektos/act${NC}"
    exit 1
fi

WEBKIT_REPO_PATH=""
DRY_RUN=false

# Help message
show_help() {
    cat << EOF
Usage: ./test-with-act.sh <webkit-repo-path> [OPTIONS]

Test Ansible playbooks using act (GitHub Actions locally) with your actual
webkit-enabled repository configuration.

ARGUMENTS:
    webkit-repo-path    Path to your webkit-enabled repository
                        (e.g., ~/projects/playground)

OPTIONS:
    -n, --dry-run       Show what would run without executing
    -h, --help          Show this help message

EXAMPLES:
    # Test with your playground repo
    ./test-with-act.sh ~/projects/playground

    # Dry run to see what would execute
    ./test-with-act.sh ~/projects/playground --dry-run

SETUP:
    1. Install act: brew install act
    2. Point script at your webkit-enabled repo
    3. The script will run ansible playbook in a GitHub Actions environment

This tests your ansible in the same environment as your actual CI/CD pipeline.

EOF
}

# Parse arguments
if [ $# -eq 0 ]; then
    echo -e "${RED}Error: webkit-repo-path is required${NC}"
    show_help
    exit 1
fi

# First argument should be the webkit repo path
if [[ "$1" != -* ]]; then
    WEBKIT_REPO_PATH="$1"
    shift
else
    echo -e "${RED}Error: webkit-repo-path is required as first argument${NC}"
    show_help
    exit 1
fi

# Parse remaining options
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Validate webkit repo path
if [ ! -d "$WEBKIT_REPO_PATH" ]; then
    echo -e "${RED}Error: Directory not found: $WEBKIT_REPO_PATH${NC}"
    exit 1
fi

# Convert to absolute path
WEBKIT_REPO_PATH="$(cd "$WEBKIT_REPO_PATH" && pwd)"

# Check for app.json
if [ ! -f "$WEBKIT_REPO_PATH/app.json" ]; then
    echo -e "${RED}Error: app.json not found in $WEBKIT_REPO_PATH${NC}"
    echo -e "${YELLOW}Make sure you're pointing to a webkit-enabled repository${NC}"
    exit 1
fi

# Check for .github/workflows in the repo
if [ ! -d "$WEBKIT_REPO_PATH/.github/workflows" ]; then
    echo -e "${YELLOW}Warning: No .github/workflows directory found in $WEBKIT_REPO_PATH${NC}"
    echo -e "${YELLOW}act will look for workflow files there${NC}"
    echo ""
    echo -e "${BLUE}If you need to test ansible deployment, your repo should have:${NC}"
    echo -e "${BLUE}.github/workflows/deploy.yml (or similar)${NC}"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo -e "${GREEN}Using webkit repo: $WEBKIT_REPO_PATH${NC}"

# Run act
echo -e "${BLUE}Running GitHub Actions locally with act...${NC}"
echo ""

cd "$WEBKIT_REPO_PATH"

if [ "$DRY_RUN" = true ]; then
    echo -e "${YELLOW}Dry run mode - showing what would execute:${NC}"
    act --list
else
    # Run act with reasonable defaults
    # -j specifies job name (deployment/release/etc)
    # --container-architecture linux/amd64 forces amd64
    # -P ubuntu-latest=catthehacker/ubuntu:act-latest uses a better image

    echo -e "${YELLOW}Note: act will ask which image to use on first run${NC}"
    echo -e "${YELLOW}Recommend: Medium (catthehacker/ubuntu:act-latest)${NC}"
    echo ""

    act \
        --container-architecture linux/amd64 \
        -P ubuntu-latest=catthehacker/ubuntu:act-latest \
        --verbose
fi
