# Installation

WebKit is distributed as a single binary with no external dependencies beyond the tools you're already using (Git, Docker, and optionally Terraform).

## Quick install

The fastest way to install WebKit is using our install script:

```bash
curl -sSL https://raw.githubusercontent.com/ainsleydev/webkit/main/bin/install.sh | sh
```

This script detects your operating system and architecture, downloads the appropriate binary, and installs it to `/usr/local/bin`.

## Manual installation

### Download pre-built binaries

1. Visit the [releases page](https://github.com/ainsleydev/webkit/releases/latest)
2. Download the archive for your platform:
   - **macOS (Intel)**: `webkit_darwin_amd64.tar.gz`
   - **macOS (Apple Silicon)**: `webkit_darwin_arm64.tar.gz`
   - **Linux (64-bit)**: `webkit_linux_amd64.tar.gz`
   - **Linux (ARM)**: `webkit_linux_arm64.tar.gz`
   - **Windows**: `webkit_windows_amd64.zip`

3. Extract the archive:
   ```bash
   tar -xzf webkit_*.tar.gz
   ```

4. Move the binary to your PATH:
   ```bash
   sudo mv webkit /usr/local/bin/
   ```

5. Verify the installation:
   ```bash
   webkit version
   ```

### Install from source

If you have Go 1.23 or higher installed:

```bash
git clone https://github.com/ainsleydev/webkit.git
cd webkit
make build
sudo mv bin/webkit /usr/local/bin/
```

## Verify installation

Check that WebKit is installed correctly:

```bash
webkit version
```

You should see output similar to:

```
webkit version 1.0.0 (commit: abc123, built: 2025-01-15)
```

## Prerequisites

Before using WebKit, ensure you have these tools installed:

### Required

- **Git** - Version control for your project
- **Docker** - For local development with Docker Compose
- **age** - Encryption tool for secrets management
  ```bash
  # macOS
  brew install age
  
  # Linux
  apt-get install age  # Debian/Ubuntu
  yum install age      # RHEL/CentOS
  ```

### Optional (for infrastructure deployment)

- **Terraform** - For deploying infrastructure (WebKit can work with Terraform 1.0+)
  ```bash
  # macOS
  brew install terraform
  
  # Linux - download from https://www.terraform.io/downloads
  ```

- **GitHub CLI** (`gh`) - For enhanced GitHub integration
  ```bash
  # macOS
  brew install gh
  
  # Linux
  curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
  sudo apt update
  sudo apt install gh
  ```

## Shell completion (optional)

WebKit supports shell completion for bash, zsh, and fish. To enable it:

### Bash

```bash
webkit completion bash > /etc/bash_completion.d/webkit
```

### Zsh

```bash
webkit completion zsh > "${fpath[1]}/_webkit"
```

### Fish

```bash
webkit completion fish > ~/.config/fish/completions/webkit.fish
```

## Next steps

Now that WebKit is installed, you're ready to:

- [Create your first project](/getting-started/quick-start)
- [Learn core concepts](/core-concepts/overview)
- [Explore the CLI reference](/cli/overview)
