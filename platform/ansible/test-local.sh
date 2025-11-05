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

# Default test variables (override with environment variables)
TEST_DOMAIN="${TEST_DOMAIN:-test.local}"
TEST_GITHUB_USER="${TEST_GITHUB_USER:-testuser}"
TEST_DOCKER_IMAGE="${TEST_DOCKER_IMAGE:-test-app}"
TEST_DOCKER_IMAGE_TAG="${TEST_DOCKER_IMAGE_TAG:-latest}"
TEST_DOCKER_PORT="${TEST_DOCKER_PORT:-3000}"
TEST_APP_NAME="${TEST_APP_NAME:-testapp}"
TEST_ENV_NAME="${TEST_ENV_NAME:-development}"
TEST_AGE_SECRET_KEY="${TEST_AGE_SECRET_KEY:-fake-key-for-testing}"
TEST_GITHUB_TOKEN="${TEST_GITHUB_TOKEN:-fake-token}"
TEST_ADMIN_EMAIL="${TEST_ADMIN_EMAIL:-test@example.com}"

# Flags
KEEP_RUNNING=false
CLEAN=false
SKIP_BUILD=false
VERBOSE=false

# Help message
show_help() {
    cat << EOF
Usage: ./test-local.sh [OPTIONS]

Test Ansible playbooks locally using Docker.

OPTIONS:
    -k, --keep-running    Keep container running after test
    -c, --clean          Remove container and image before starting
    -s, --skip-build     Skip Docker image build
    -v, --verbose        Show verbose Ansible output
    -h, --help           Show this help message

ENVIRONMENT VARIABLES:
    You can override test variables by setting these before running:
    TEST_DOMAIN, TEST_GITHUB_USER, TEST_DOCKER_IMAGE, TEST_DOCKER_IMAGE_TAG,
    TEST_DOCKER_PORT, TEST_APP_NAME, TEST_ENV_NAME, TEST_AGE_SECRET_KEY,
    TEST_GITHUB_TOKEN, TEST_ADMIN_EMAIL

EXAMPLES:
    # Basic test run
    ./test-local.sh

    # Keep container running for debugging
    ./test-local.sh -k

    # Clean rebuild
    ./test-local.sh -c

    # Custom domain
    TEST_DOMAIN=myapp.test ./test-local.sh

EOF
}

# Parse arguments
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

# Cleanup function
cleanup() {
    if [ "$KEEP_RUNNING" = false ]; then
        echo -e "${YELLOW}Cleaning up...${NC}"
        docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
    else
        echo -e "${GREEN}Container '$CONTAINER_NAME' is still running for debugging${NC}"
        echo -e "${BLUE}To connect: docker exec -it $CONTAINER_NAME /bin/bash${NC}"
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
    docker build -f "$ANSIBLE_DIR/Dockerfile.test" -t "$IMAGE_NAME" "$ANSIBLE_DIR"
fi

# Stop existing container if running
docker rm -f "$CONTAINER_NAME" 2>/dev/null || true

# Start container with systemd
echo -e "${BLUE}Starting test container...${NC}"
docker run -d \
    --name "$CONTAINER_NAME" \
    --privileged \
    --tmpfs /tmp \
    --tmpfs /run \
    --tmpfs /run/lock \
    -v /sys/fs/cgroup:/sys/fs/cgroup:ro \
    "$IMAGE_NAME"

# Wait for container to be ready
echo -e "${BLUE}Waiting for container to be ready...${NC}"
sleep 3

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo -e "${RED}Container failed to start!${NC}"
    docker logs "$CONTAINER_NAME"
    exit 1
fi

echo -e "${GREEN}Container is running${NC}"

# Install Ansible in container (faster than SSH setup)
echo -e "${BLUE}Installing Ansible in container...${NC}"
docker exec "$CONTAINER_NAME" bash -c "
    apt-get update > /dev/null && \
    apt-get install -y ansible python3-pip > /dev/null && \
    pip3 install community.docker > /dev/null 2>&1 || true
"

# Copy Ansible files into container
echo -e "${BLUE}Copying Ansible files to container...${NC}"
docker cp "$ANSIBLE_DIR" "$CONTAINER_NAME:/tmp/ansible"

# Create inventory file
echo -e "${BLUE}Creating inventory file...${NC}"
docker exec "$CONTAINER_NAME" bash -c "echo 'localhost ansible_connection=local' > /tmp/ansible/inventory"

# Build ansible-playbook command
ANSIBLE_CMD="ansible-playbook /tmp/ansible/playbooks/server.yaml -i /tmp/ansible/inventory"

# Add extra vars
ANSIBLE_CMD="$ANSIBLE_CMD -e domain=$TEST_DOMAIN"
ANSIBLE_CMD="$ANSIBLE_CMD -e github_user=$TEST_GITHUB_USER"
ANSIBLE_CMD="$ANSIBLE_CMD -e docker_image=$TEST_DOCKER_IMAGE"
ANSIBLE_CMD="$ANSIBLE_CMD -e docker_image_tag=$TEST_DOCKER_IMAGE_TAG"
ANSIBLE_CMD="$ANSIBLE_CMD -e docker_port=$TEST_DOCKER_PORT"
ANSIBLE_CMD="$ANSIBLE_CMD -e app_name=$TEST_APP_NAME"
ANSIBLE_CMD="$ANSIBLE_CMD -e env_name=$TEST_ENV_NAME"
ANSIBLE_CMD="$ANSIBLE_CMD -e age_secret_key=$TEST_AGE_SECRET_KEY"
ANSIBLE_CMD="$ANSIBLE_CMD -e github_token=$TEST_GITHUB_TOKEN"
ANSIBLE_CMD="$ANSIBLE_CMD -e admin_email=$TEST_ADMIN_EMAIL"
ANSIBLE_CMD="$ANSIBLE_CMD -e enable_https=false"
ANSIBLE_CMD="$ANSIBLE_CMD -e skip_reboot=true"

# Add verbose flag if requested
if [ "$VERBOSE" = true ]; then
    ANSIBLE_CMD="$ANSIBLE_CMD -vvv"
fi

# Run Ansible playbook
echo -e "${BLUE}Running Ansible playbook...${NC}"
echo -e "${YELLOW}Command: $ANSIBLE_CMD${NC}"
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
