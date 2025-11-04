#!/bin/bash
# Integration test for webkit Ansible role
# This test verifies that the webkit role correctly installs and configures webkit

set -e

TEST_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROLE_DIR="$(dirname "$TEST_DIR")"
ANSIBLE_DIR="$(dirname "$(dirname "$ROLE_DIR")")"

echo "==> Running webkit Ansible role integration test"

# Check prerequisites
if ! command -v ansible-playbook &> /dev/null; then
    echo "ERROR: ansible-playbook not found. Please install Ansible."
    exit 1
fi

if ! command -v docker &> /dev/null; then
    echo "ERROR: docker not found. This test requires Docker to run."
    exit 1
fi

# Create a test container to act as a "VM"
echo "==> Creating test container"
CONTAINER_NAME="webkit-ansible-test-$(date +%s)"
docker run -d --name "$CONTAINER_NAME" \
    --privileged \
    ubuntu:22.04 \
    sleep 3600

# Ensure cleanup on exit
trap "docker rm -f $CONTAINER_NAME > /dev/null 2>&1" EXIT

# Install Python in the container (required for Ansible)
echo "==> Installing Python in test container"
docker exec "$CONTAINER_NAME" apt-get update -qq
docker exec "$CONTAINER_NAME" apt-get install -y -qq python3 python3-pip curl sudo

# Create a temporary inventory file
TMP_INVENTORY=$(mktemp)
cat > "$TMP_INVENTORY" <<EOF
[test]
$CONTAINER_NAME ansible_connection=docker ansible_user=root
EOF

# Create a test playbook
TMP_PLAYBOOK=$(mktemp)
cat > "$TMP_PLAYBOOK" <<'EOF'
---
- name: Test webkit role
  hosts: test
  become: true
  vars:
    webkit_version: latest
    age_secret_key: "AGE-SECRET-KEY-1TEST123456789"
    app_definition_path: "{{ playbook_dir }}/test_fixtures/app.json"
    secrets_path: "{{ playbook_dir }}/test_fixtures/secrets"

  roles:
    - webkit

  post_tasks:
    - name: Verify webkit is installed
      command: /usr/local/bin/webkit version
      register: webkit_version_output
      changed_when: false

    - name: Debug webkit version
      debug:
        msg: "{{ webkit_version_output.stdout }}"

    - name: Verify app.json was copied
      stat:
        path: /etc/webkit/app.json
      register: app_json_stat

    - name: Assert app.json exists
      assert:
        that:
          - app_json_stat.stat.exists
        fail_msg: "app.json was not copied to /etc/webkit/"

    - name: Verify AGE key was written
      stat:
        path: /root/.config/webkit/age.key
      register: age_key_stat

    - name: Assert AGE key exists and has correct permissions
      assert:
        that:
          - age_key_stat.stat.exists
          - age_key_stat.stat.mode == '0600'
        fail_msg: "AGE key was not created or has wrong permissions"
EOF

# Create test fixtures
TMP_FIXTURES_DIR=$(mktemp -d)
mkdir -p "$TMP_FIXTURES_DIR/secrets"

cat > "$TMP_FIXTURES_DIR/app.json" <<'EOF'
{
  "webkit_version": "v0.0.15",
  "apps": [
    {
      "name": "test-app",
      "path": "./test-app",
      "env": {
        "production": {
          "FOO": {
            "source": "value",
            "value": "bar"
          }
        }
      }
    }
  ]
}
EOF

cat > "$TMP_FIXTURES_DIR/secrets/production.yaml" <<'EOF'
TEST_SECRET: test_value
EOF

# Copy fixtures to test location
TMP_TEST_DIR=$(mktemp -d)
cp -r "$TMP_FIXTURES_DIR" "$TMP_TEST_DIR/test_fixtures"

# Run the Ansible playbook
echo "==> Running Ansible playbook"
cd "$TMP_TEST_DIR"
ANSIBLE_ROLES_PATH="$ANSIBLE_DIR/roles" \
    ansible-playbook \
    -i "$TMP_INVENTORY" \
    "$TMP_PLAYBOOK" \
    -v

# Additional verification in the container
echo "==> Verifying webkit installation"
docker exec "$CONTAINER_NAME" /usr/local/bin/webkit version

echo "==> Verifying file permissions"
docker exec "$CONTAINER_NAME" stat -c "%a %n" /root/.config/webkit/age.key

echo "==> Test passed! ✓"

# Cleanup
rm -f "$TMP_INVENTORY" "$TMP_PLAYBOOK"
rm -rf "$TMP_FIXTURES_DIR" "$TMP_TEST_DIR"
