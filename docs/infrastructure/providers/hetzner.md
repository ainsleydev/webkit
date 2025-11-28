# Hetzner

Hetzner Cloud provides cost-effective virtual machines and storage in European and US data centres. It's an excellent choice for VM-based deployments where you need more control than managed platforms offer.

## Authentication

Set your Hetzner Cloud API token as an environment variable:

```bash
export HCLOUD_TOKEN="your-token"
```

Generate a token at [console.hetzner.cloud](https://console.hetzner.cloud/) under your project's Security → API Tokens.

## Virtual machines

Hetzner VMs (called "servers") offer excellent price-to-performance ratio for self-managed deployments.

### Configuration

```json
{
  "apps": [
    {
      "name": "api",
      "type": "golang",
      "path": "./apps/api",
      "infrastructure": {
        "provider": "hetzner",
        "type": "vm",
        "config": {
          "server_type": "cx22",
          "location": "nbg1",
          "image": "ubuntu-22.04"
        }
      }
    }
  ]
}
```

### Server types

Hetzner offers shared (CX) and dedicated (CCX/CPX) CPU options:

#### Shared vCPU (CX series)

| Type | vCPUs | Memory | Storage | Monthly cost |
|------|-------|--------|---------|--------------|
| `cx22` | 2 | 4GB | 40GB | ~€4 |
| `cx32` | 4 | 8GB | 80GB | ~€7 |
| `cx42` | 8 | 16GB | 160GB | ~€14 |
| `cx52` | 16 | 32GB | 320GB | ~€29 |

#### Dedicated vCPU (CCX series)

| Type | vCPUs | Memory | Storage | Monthly cost |
|------|-------|--------|---------|--------------|
| `ccx13` | 2 | 8GB | 80GB | ~€13 |
| `ccx23` | 4 | 16GB | 160GB | ~€25 |
| `ccx33` | 8 | 32GB | 240GB | ~€50 |
| `ccx43` | 16 | 64GB | 360GB | ~€100 |

#### ARM64 (CAX series)

| Type | vCPUs | Memory | Storage | Monthly cost |
|------|-------|--------|---------|--------------|
| `cax11` | 2 | 4GB | 40GB | ~€4 |
| `cax21` | 4 | 8GB | 80GB | ~€6 |
| `cax31` | 8 | 16GB | 160GB | ~€12 |
| `cax41` | 16 | 32GB | 320GB | ~€24 |

ARM servers offer excellent value for compatible workloads.

### Locations

| Code | Location | Region |
|------|----------|--------|
| `nbg1` | Nuremberg | Germany |
| `fsn1` | Falkenstein | Germany |
| `hel1` | Helsinki | Finland |
| `ash` | Ashburn | USA East |
| `hil` | Hillsboro | USA West |

### Images

| Image | Description |
|-------|-------------|
| `ubuntu-22.04` | Ubuntu 22.04 LTS |
| `ubuntu-24.04` | Ubuntu 24.04 LTS |
| `debian-12` | Debian 12 (Bookworm) |
| `rocky-9` | Rocky Linux 9 |
| `fedora-40` | Fedora 40 |

### SSH configuration

WebKit configures SSH access to Hetzner VMs. Ensure your SSH public key is available:

```json
{
  "infrastructure": {
    "provider": "hetzner",
    "type": "vm",
    "config": {
      "ssh_keys": ["your-key-name"]
    }
  }
}
```

SSH keys are managed at the Hetzner project level.

## Volumes

Hetzner Volumes provide additional block storage for your VMs.

### Configuration

```json
{
  "resources": [
    {
      "name": "data",
      "type": "volume",
      "provider": "hetzner",
      "config": {
        "size": 100,
        "location": "nbg1",
        "format": "ext4"
      }
    }
  ]
}
```

### Volume sizes

Volumes range from 10GB to 10TB. Pricing is ~€0.05/GB/month.

| Size | Monthly cost |
|------|--------------|
| 10GB | ~€0.50 |
| 50GB | ~€2.50 |
| 100GB | ~€5 |
| 500GB | ~€25 |
| 1TB | ~€50 |

### Mount points

Volumes are automatically attached to the associated VM. Configure the mount point:

```json
{
  "config": {
    "mount_point": "/data"
  }
}
```

WebKit configures the VM to mount the volume at boot.

### Outputs

| Output | Description |
|--------|-------------|
| `data.id` | Volume ID |
| `data.size` | Volume size in GB |
| `data.mount_point` | Mount path on VM |

## Networking

### Private networks

Enable private networking between VMs:

```json
{
  "infrastructure": {
    "provider": "hetzner",
    "type": "vm",
    "config": {
      "network": {
        "enabled": true,
        "ip_range": "10.0.0.0/24"
      }
    }
  }
}
```

VMs in the same network can communicate over private IPs.

### Firewalls

Configure firewall rules:

```json
{
  "config": {
    "firewall": {
      "rules": [
        {
          "direction": "in",
          "protocol": "tcp",
          "port": "22",
          "source_ips": ["0.0.0.0/0"]
        },
        {
          "direction": "in",
          "protocol": "tcp",
          "port": "80",
          "source_ips": ["0.0.0.0/0"]
        },
        {
          "direction": "in",
          "protocol": "tcp",
          "port": "443",
          "source_ips": ["0.0.0.0/0"]
        }
      ]
    }
  }
}
```

## Deployment

Unlike managed platforms, Hetzner VMs require deployment configuration. WebKit generates:

- Ansible playbooks for initial setup
- GitHub Actions workflows for deployment
- Docker configuration for containerised apps

### Ansible provisioning

WebKit generates Ansible playbooks to configure your VM:

- Install Docker and required packages
- Configure firewall rules
- Set up SSL certificates (Let's Encrypt)
- Deploy your application

Run provisioning:

```bash
webkit infra apply  # Creates VM
# Ansible runs automatically via GitHub Actions
```

### Manual deployment

For manual deployments, SSH into your VM:

```bash
ssh root@your-vm-ip
cd /opt/app
docker compose pull
docker compose up -d
```

## Example: Self-hosted application

A Go API deployed on Hetzner with attached storage:

```json
{
  "project": {
    "name": "my-api",
    "title": "My API",
    "repo": "github.com/myorg/my-api"
  },
  "apps": [
    {
      "name": "api",
      "type": "golang",
      "path": "./apps/api",
      "build": {
        "dockerfile": true,
        "port": 8080
      },
      "infrastructure": {
        "provider": "hetzner",
        "type": "vm",
        "config": {
          "server_type": "cx22",
          "location": "nbg1",
          "image": "ubuntu-22.04"
        }
      },
      "domains": {
        "primary": "api.example.com"
      }
    }
  ],
  "resources": [
    {
      "name": "uploads",
      "type": "volume",
      "provider": "hetzner",
      "config": {
        "size": 50,
        "location": "nbg1",
        "mount_point": "/data/uploads"
      }
    }
  ]
}
```

## Comparison with DigitalOcean

| Feature | Hetzner | DigitalOcean |
|---------|---------|--------------|
| Managed apps | No | Yes (App Platform) |
| VM pricing | Lower | Higher |
| Managed DB | No | Yes |
| European locations | Excellent | Good |
| US locations | 2 | Many |
| ARM servers | Yes | No |
| Setup complexity | Higher | Lower |

Choose Hetzner when:
- Cost is a priority
- You need ARM servers
- You're comfortable with VM management
- Your users are primarily in Europe

Choose DigitalOcean when:
- You want managed services
- Simplicity is important
- You need managed databases
- Your users are globally distributed

## Further reading

- [Hetzner Cloud documentation](https://docs.hetzner.com/cloud)
- [Server types](https://docs.hetzner.com/cloud/servers/overview)
- [Volumes documentation](https://docs.hetzner.com/cloud/volumes/overview)
- [Pricing calculator](https://www.hetzner.com/cloud)
