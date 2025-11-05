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

## Local Testing

Test Ansible changes locally before deploying to production. Two approaches available:

### Option 1: Quick Syntax Check (Recommended)

Fast validation without installing packages or using Docker:

```bash
cd platform/ansible

# Syntax check only
./test-local-simple.sh ~/projects/your-webkit-repo

# Full dry-run validation (checks logic and templates)
./test-local-simple.sh ~/projects/your-webkit-repo --check
```

**Requirements:** `brew install ansible jq`

**What it validates:**
- Playbook syntax
- Variable usage from your app.json
- Template rendering
- Task logic

**What it doesn't do:**
- Install packages
- Start services
- Modify your system

### Option 2: Test with act (GitHub Actions locally)

Test your complete deployment workflow if your webkit-enabled repo has GitHub Actions:

```bash
# Install act
brew install act

# Run your workflow locally
./test-with-act.sh ~/projects/your-webkit-repo
```

**Note:** This tests your user repo's workflow, which uses a tagged webkit version. To test local webkit ansible changes, use Option 1.

### When to Use Which

**Use Option 1 (syntax check)** when:
- You changed ansible playbooks/roles in webkit
- You want fast feedback
- You want to validate before committing

**Use Option 2 (act)** when:
- You want to test the full CI/CD pipeline
- Your webkit-enabled repo has GitHub Actions
- You want to test with real deployment conditions

### Examples

```bash
# Test webkit ansible changes before committing
cd platform/ansible
./test-local-simple.sh ~/projects/playground

# Validate a production config (safely)
./test-local-simple.sh ~/projects/my-production-app --check

# Test full workflow with act
./test-with-act.sh ~/projects/playground
```
