# Ansible Deployment

Ansible playbooks and roles for deploying WebKit applications to DigitalOcean VMs.

## Usage

This directory is automatically copied to user repositories during VM deployment workflows. When a project with a VM-based app is released, the GitHub Actions workflow:

1. Reads `webkit_version` from the user's `app.json`
2. Checks out this WebKit repository at that version
3. Copies this `platform/ansible` directory into the workspace
4. Runs the playbook with app-specific configuration

This approach keeps user repositories clean while ensuring ansible files are available during deployment.

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
