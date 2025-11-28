# Installation

This guide covers how to install WebKit on your system.

## Prerequisites

Before installing WebKit, ensure you have the following:

- **Go 1.25+** - Required for building from source or using `go install`
- **Git** - For version control operations
- **SOPS** (optional) - Required only if you're using encrypted secrets
- **Age** (optional) - Required only if you're using SOPS encryption

## Quick install

The fastest way to install WebKit is using our install script:

```bash
curl -sSL https://raw.githubusercontent.com/ainsleydev/webkit/main/bin/install.sh | sh
```

This automatically detects your operating system and architecture, downloads the appropriate binary, and installs it to `/usr/local/bin`.

## Install with Go

If you have Go installed, you can install WebKit directly:

```bash
go install github.com/ainsleydev/webkit@latest
```

This installs the `webkit` binary to your `$GOPATH/bin` directory. Ensure this directory is in your `PATH`.

## Download binary

You can download pre-built binaries from the [GitHub releases page](https://github.com/ainsleydev/webkit/releases/latest).

Available platforms:
- **macOS** (Darwin) - Intel (`x86_64`) and Apple Silicon (`arm64`)
- **Linux** - Intel (`x86_64`) and ARM (`arm64`)
- **Windows** - Intel (`x86_64`)

After downloading, extract the archive and move the binary to a directory in your `PATH`:

```bash
# macOS/Linux example
tar -xzf webkit_darwin_arm64.tar.gz
sudo mv webkit /usr/local/bin/
```

## Verify installation

After installation, verify WebKit is working correctly:

```bash
webkit version
```

You should see output similar to:

```
WebKit v0.9.2
```

## Optional dependencies

### SOPS and Age

If you plan to use encrypted secrets with WebKit, you'll need to install SOPS and Age:

**macOS (Homebrew):**
```bash
brew install sops age
```

**Linux:**
```bash
# SOPS
curl -LO https://github.com/getsops/sops/releases/download/v3.9.0/sops-v3.9.0.linux.amd64
sudo mv sops-v3.9.0.linux.amd64 /usr/local/bin/sops
sudo chmod +x /usr/local/bin/sops

# Age
curl -LO https://github.com/FiloSottile/age/releases/download/v1.2.0/age-v1.2.0-linux-amd64.tar.gz
tar -xzf age-v1.2.0-linux-amd64.tar.gz
sudo mv age/age /usr/local/bin/
sudo mv age/age-keygen /usr/local/bin/
```

### Terraform

WebKit uses Terraform for infrastructure provisioning. Install Terraform if you plan to deploy infrastructure:

**macOS (Homebrew):**
```bash
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```

**Linux:**
```bash
# Add HashiCorp GPG key
wget -O- https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg

# Add repository
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list

# Install
sudo apt update && sudo apt install terraform
```

## Next steps

Once WebKit is installed, you're ready to:

- Follow the [quick start guide](/getting-started/quick-start) for a 5-minute introduction
- Build your [first project](/getting-started/your-first-project) with a complete tutorial
- Learn about [core concepts](/getting-started/core-concepts) behind WebKit
