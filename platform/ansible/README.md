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

Test Ansible changes locally using your actual webkit-enabled repositories (like playground or production repos) without deploying to live servers.

### Quick Start

```bash
cd platform/ansible
./test-local.sh ~/path/to/your-webkit-repo
```

The script will:
1. Read your repo's `app.json` configuration
2. Spin up a Ubuntu 22.04 Docker container
3. Run the Ansible playbook with your real config
4. Show results and clean up

**That's it!** No need to set variables - it uses your actual repository configuration.

### Examples

```bash
# Test with playground repo
./test-local.sh ~/projects/playground

# Test with production config (safely!)
./test-local.sh ~/projects/my-production-app

# Keep container running for debugging
./test-local.sh ~/projects/playground -k

# Verbose Ansible output
./test-local.sh ~/projects/playground -v

# Clean rebuild
./test-local.sh ~/projects/playground -c
```

### Options

- `-k, --keep-running` - Keep container running after test for inspection
- `-c, --clean` - Remove existing container/image before starting
- `-s, --skip-build` - Skip Docker image rebuild (faster iteration)
- `-v, --verbose` - Show verbose Ansible output

### What Gets Tested

The script automatically extracts from your `app.json`:
- Domain configuration
- App name and port
- GitHub repository details
- Admin email
- HTTPS settings

Then runs the full playbook against a fresh Ubuntu container, testing:
- All role installations (fail2ban, docker, nginx, ufw, etc.)
- Configuration file templating
- Service setup and management
- Package installations

### Debugging

When using `-k`, inspect the container after the test:

```bash
# Connect to container
docker exec -it webkit-ansible-test /bin/bash

# Check installed services
docker exec webkit-ansible-test systemctl status nginx
docker exec webkit-ansible-test systemctl status docker

# View webkit config
docker exec webkit-ansible-test cat /etc/webkit/app.json

# Stop when done
docker rm -f webkit-ansible-test
```

### Benefits

- **Use real configs**: Test with your actual `app.json` files
- **No production risk**: Runs in isolated Docker containers
- **Fast iteration**: Test changes in seconds, not minutes
- **Catch errors early**: Find issues before they hit live servers
- **Multiple repos**: Test against playground, staging, or production configs

### Expected Failures

Some tasks will fail (by design - they require live infrastructure):
- Certbot/HTTPS (requires DNS validation)
- Docker image pulls (requires authentication)
- Webkit env generation (requires decryption keys)

These are normal and can be ignored. The script automatically sets `enable_https=false` and `skip_reboot=true`.
