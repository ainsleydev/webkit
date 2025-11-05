#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="webkit-ansible-test"
IMAGE_NAME="webkit-ansible-test"
ANSIBLE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Flags
KEEP_RUNNING=false
CLEAN=false
SKIP_BUILD=false
VERBOSE=false
WEBKIT_REPO_PATH=""

# Help message
show_help() {
    cat << EOF
Usage: ./test-local.sh <webkit-repo-path> [OPTIONS]

Test Ansible playbooks locally using a real webkit-enabled repository.

ARGUMENTS:
    webkit-repo-path    Path to your webkit-enabled repository
                        (e.g., ../my-webkit-app or ~/projects/playground)

OPTIONS:
    -k, --keep-running    Keep container running after test
    -c, --clean          Remove container and image before starting
    -s, --skip-build     Skip Docker image build
    -v, --verbose        Show verbose Ansible output
    -h, --help           Show this help message

EXAMPLES:
    # Test with your playground repo
    ./test-local.sh ~/projects/playground

    # Keep container running for debugging
    ./test-local.sh ~/projects/playground -k

    # Verbose output
    ./test-local.sh ~/projects/playground -v

    # Clean rebuild and test
    ./test-local.sh ~/projects/playground -c

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
        -k|--keep-running)
            KEEP_RUNNING=true
            shift
            ;;
        -c|--clean)
            CLEAN=true
            shift
            ;;
        -s|--skip-build)
            SKIP_BUILD=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
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

echo -e "${GREEN}Using webkit repo: $WEBKIT_REPO_PATH${NC}"

# Cleanup function
cleanup() {
    if [ "$KEEP_RUNNING" = false ]; then
        echo -e "${YELLOW}Cleaning up...${NC}"
        docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
    else
        echo -e "${GREEN}Container '$CONTAINER_NAME' is still running for debugging${NC}"
        echo -e "${BLUE}To connect: docker exec -it $CONTAINER_NAME /bin/bash${NC}"
        echo -e "${BLUE}To view webkit config: docker exec $CONTAINER_NAME ls -la /etc/webkit${NC}"
        echo -e "${BLUE}To stop: docker rm -f $CONTAINER_NAME${NC}"
    fi
}

# Set trap for cleanup
trap cleanup EXIT

# Clean if requested
if [ "$CLEAN" = true ]; then
    echo -e "${YELLOW}Cleaning up existing container and image...${NC}"
    docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
    docker rmi "$IMAGE_NAME" 2>/dev/null || true
fi

# Build Docker image
if [ "$SKIP_BUILD" = false ]; then
    echo -e "${BLUE}Building test Docker image...${NC}"
    docker build --platform linux/amd64 -f "$ANSIBLE_DIR/Dockerfile.test" -t "$IMAGE_NAME" "$ANSIBLE_DIR"
fi

# Stop existing container if running
docker rm -f "$CONTAINER_NAME" 2>/dev/null || true

# Start container
echo -e "${BLUE}Starting test container...${NC}"
docker run -d \
    --name "$CONTAINER_NAME" \
    --platform linux/amd64 \
    -v "$WEBKIT_REPO_PATH:/webkit-repo:ro" \
    "$IMAGE_NAME"

# Wait for container to be ready
echo -e "${BLUE}Waiting for container to be ready...${NC}"
sleep 2

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}Container failed to start!${NC}"
    docker logs "$CONTAINER_NAME"
    exit 1
fi

echo -e "${GREEN}Container is running${NC}"

# Install Ansible in container
echo -e "${BLUE}Installing Ansible in container...${NC}"
docker exec "$CONTAINER_NAME" bash -c "
    apt-get update > /dev/null && \
    apt-get install -y ansible python3-pip jq > /dev/null && \
    pip3 install community.docker > /dev/null 2>&1 || true
"

# Copy Ansible files into container
echo -e "${BLUE}Copying Ansible files to container...${NC}"
docker cp "$ANSIBLE_DIR" "$CONTAINER_NAME:/tmp/ansible"

# Setup webkit config in container
echo -e "${BLUE}Setting up webkit configuration...${NC}"
docker exec "$CONTAINER_NAME" bash -c "
    mkdir -p /etc/webkit
    cp /webkit-repo/app.json /etc/webkit/
    [ -d /webkit-repo/resources ] && cp -r /webkit-repo/resources /etc/webkit/ || true
"

# Extract variables from app.json
echo -e "${BLUE}Reading configuration from app.json...${NC}"
APP_JSON_CONTENT=$(docker exec "$CONTAINER_NAME" cat /etc/webkit/app.json)

# Use jq to extract values (with defaults)
DOMAIN=$(echo "$APP_JSON_CONTENT" | docker exec -i "$CONTAINER_NAME" jq -r '.apps[0].primary_domain // "test.local"')
APP_NAME=$(echo "$APP_JSON_CONTENT" | docker exec -i "$CONTAINER_NAME" jq -r '.apps[0].name // "testapp"')
DOCKER_PORT=$(echo "$APP_JSON_CONTENT" | docker exec -i "$CONTAINER_NAME" jq -r '.apps[0].build.port // 3000')
GITHUB_USER=$(echo "$APP_JSON_CONTENT" | docker exec -i "$CONTAINER_NAME" jq -r '.repository.owner // "testuser"')
REPO_NAME=$(echo "$APP_JSON_CONTENT" | docker exec -i "$CONTAINER_NAME" jq -r '.repository.name // "testrepo"')
ENABLE_HTTPS=$(echo "$APP_JSON_CONTENT" | docker exec -i "$CONTAINER_NAME" jq -r '.apps[0].infra.config.https // true')
ADMIN_EMAIL=$(echo "$APP_JSON_CONTENT" | docker exec -i "$CONTAINER_NAME" jq -r '.apps[0].infra.config.admin_email // "test@example.com"')

# Display config
echo -e "${GREEN}Configuration:${NC}"
echo -e "  Domain: $DOMAIN"
echo -e "  App: $APP_NAME"
echo -e "  Port: $DOCKER_PORT"
echo -e "  GitHub: $GITHUB_USER/$REPO_NAME"
echo ""

# Create inventory file
docker exec "$CONTAINER_NAME" bash -c "echo 'localhost ansible_connection=local' > /tmp/ansible/inventory"

# Build ansible-playbook command
ANSIBLE_CMD="ansible-playbook /tmp/ansible/playbooks/server.yaml -i /tmp/ansible/inventory"

# Add vars from app.json
ANSIBLE_CMD="$ANSIBLE_CMD -e domain=$DOMAIN"
ANSIBLE_CMD="$ANSIBLE_CMD -e github_user=$GITHUB_USER"
ANSIBLE_CMD="$ANSIBLE_CMD -e docker_image=${REPO_NAME}-${APP_NAME}"
ANSIBLE_CMD="$ANSIBLE_CMD -e docker_image_tag=test-latest"
ANSIBLE_CMD="$ANSIBLE_CMD -e docker_port=$DOCKER_PORT"
ANSIBLE_CMD="$ANSIBLE_CMD -e app_name=$APP_NAME"
ANSIBLE_CMD="$ANSIBLE_CMD -e env_name=development"
ANSIBLE_CMD="$ANSIBLE_CMD -e age_secret_key=fake-key-for-testing"
ANSIBLE_CMD="$ANSIBLE_CMD -e github_token=fake-token"
ANSIBLE_CMD="$ANSIBLE_CMD -e admin_email=$ADMIN_EMAIL"
ANSIBLE_CMD="$ANSIBLE_CMD -e enable_https=false"
ANSIBLE_CMD="$ANSIBLE_CMD -e skip_reboot=true"

# Add verbose flag if requested
if [ "$VERBOSE" = true ]; then
    ANSIBLE_CMD="$ANSIBLE_CMD -vvv"
fi

# Run Ansible playbook
echo -e "${BLUE}Running Ansible playbook...${NC}"
echo ""

if docker exec "$CONTAINER_NAME" bash -c "cd /tmp/ansible && $ANSIBLE_CMD"; then
    echo ""
    echo -e "${GREEN}✓ Ansible playbook completed successfully!${NC}"
    EXIT_CODE=0
else
    echo ""
    echo -e "${RED}✗ Ansible playbook failed!${NC}"
    EXIT_CODE=1
fi

exit $EXIT_CODE
