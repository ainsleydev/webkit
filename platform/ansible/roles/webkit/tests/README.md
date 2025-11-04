# WebKit Ansible Role Integration Tests

This directory contains integration tests for the WebKit Ansible role.

## Prerequisites

- Docker
- Ansible (`ansible-playbook` command)
- Bash

## Running the tests

```bash
./test_webkit_role_integration.sh
```

## What the test does

The integration test:

1. Creates a Docker container running Ubuntu 22.04 to simulate a VM.
2. Runs the webkit Ansible role against the container.
3. Verifies that:
   - WebKit binary is installed at `/usr/local/bin/webkit`.
   - `app.json` is copied to `/etc/webkit/app.json`.
   - Secrets are copied to `/etc/webkit/resources/secrets/`.
   - AGE key is written to `/root/.config/webkit/age.key` with correct permissions (600).
4. Cleans up the test container automatically.

## CI Integration

To run this test in CI:

```yaml
- name: Test Ansible WebKit Role
  run: |
    cd platform/ansible/roles/webkit/tests
    ./test_webkit_role_integration.sh
```

## Troubleshooting

If the test fails:

1. Check that Docker is running: `docker ps`.
2. Check that Ansible is installed: `ansible-playbook --version`.
3. Review the Ansible playbook output for specific errors.
4. Manually inspect the test container: `docker exec -it webkit-ansible-test-<timestamp> bash`.
