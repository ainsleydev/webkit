# Ansible Deployment

Ansible playbooks and roles for deploying WebKit applications to DigitalOcean VMs.

## Usage

This directory is automatically copied to user repositories during VM deployment workflows. When a project with a VM-based app is released, the GitHub Actions workflow:

1. Reads `webkit_version` from the user's `app.json`
2. Checks out this WebKit repository at that version
3. Copies this `platform/ansible` directory into the workspace
4. Runs the playbook with app-specific configuration

## Structure

```
platform/ansible/
├── ansible.cfg          # Ansible configuration
├── playbooks/
│   └── server.yaml     # Main deployment playbook
└── roles/              # Ansible roles
    ├── certbot/        # SSL certificate management
    ├── docker/         # Docker installation and setup
    ├── fail2ban/       # Security and intrusion prevention
    ├── nginx/          # Reverse proxy configuration
    ├── tools/          # System utilities
    ├── ufw/            # Firewall configuration
    └── webkit/         # WebKit app deployment
```

## The Playbook

The `server.yaml` playbook configures a production-ready server with:
- System updates and security hardening (UFW, fail2ban)
- Docker and container orchestration
- Nginx reverse proxy with SSL (certbot)
- Application deployment with environment variable decryption (SOPS/Age)

All configuration is passed via variables from the workflow, sourced from the user's `app.json`.

## Local Testing with Docker

Test Ansible changes locally before deploying to production using the included Docker testing environment.

### Quick Start

```bash
cd platform/ansible
./test-local.sh
```

This will:
1. Build a Ubuntu 22.04 Docker container
2. Install Ansible and dependencies
3. Run the playbook with test variables
4. Show results and clean up

### Options

```bash
# Keep container running for debugging
./test-local.sh --keep-running

# Clean rebuild (removes existing container/image)
./test-local.sh --clean

# Skip Docker image rebuild (faster iteration)
./test-local.sh --skip-build

# Verbose Ansible output
./test-local.sh --verbose

# Combine options
./test-local.sh -k -v
```

### Custom Test Variables

Override default test values with environment variables:

```bash
TEST_DOMAIN=myapp.test \
TEST_APP_NAME=myapp \
TEST_DOCKER_PORT=8080 \
./test-local.sh
```

Available variables:
- `TEST_DOMAIN` (default: test.local)
- `TEST_GITHUB_USER` (default: testuser)
- `TEST_DOCKER_IMAGE` (default: test-app)
- `TEST_DOCKER_IMAGE_TAG` (default: latest)
- `TEST_DOCKER_PORT` (default: 3000)
- `TEST_APP_NAME` (default: testapp)
- `TEST_ENV_NAME` (default: development)
- `TEST_ADMIN_EMAIL` (default: test@example.com)

### Debugging

When using `--keep-running`, the container stays active after the test:

```bash
# Connect to the container
docker exec -it webkit-ansible-test /bin/bash

# Check service status
docker exec webkit-ansible-test systemctl status nginx

# View logs
docker exec webkit-ansible-test journalctl -u nginx

# Stop when done
docker rm -f webkit-ansible-test
```

### Benefits

- **No production risk**: Test on isolated containers
- **Fast iteration**: Rebuild and test in seconds
- **Catch errors early**: Find provisioning issues before deployment
- **Full integration testing**: Tests actual package installation and service configuration

### Limitations

Some tasks will fail in the test environment (by design):
- Certbot/HTTPS configuration (DNS validation not possible)
- Docker image pulling from GHCR (uses fake credentials)
- Webkit env generation (requires real app.json and secrets)

These failures are expected and can be ignored during role development. The script sets `enable_https=false` and `skip_reboot=true` automatically.
