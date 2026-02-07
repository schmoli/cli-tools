# cli-tools

[![Release](https://img.shields.io/github/v/release/schmoli/cli-tools)](https://github.com/schmoli/cli-tools/releases)
[![License](https://img.shields.io/github/license/schmoli/cli-tools)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/schmoli/cli-tools/release-please.yml?branch=main)](https://github.com/schmoli/cli-tools/actions)

> CLI tools for Portainer, nginx-proxy-manager, Transmission, Proxmox VE, and Audiobookshelf APIs

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/schmoli/cli-tools/main/install.sh | bash
```

Installs to `~/.local/bin`.

## Configuration

| Variable | Tools | Description |
|----------|-------|-------------|
| `PORTAINER_URL` | portainer-cli | Portainer server URL |
| `PORTAINER_TOKEN` | portainer-cli | Portainer API token |
| `NPROXY_URL` | nproxy-cli | nginx-proxy-manager URL |
| `NPROXY_TOKEN` | nproxy-cli | nginx-proxy-manager JWT |
| `TRANSMISSION_URL` | trans-cli | Transmission RPC URL |
| `TRANSMISSION_USER` | trans-cli | Transmission username (optional) |
| `TRANSMISSION_PASS` | trans-cli | Transmission password (optional) |
| `PVE_URL` | pve-cli | Proxmox VE URL |
| `PVE_TOKEN_ID` | pve-cli | Proxmox API token ID (user@realm!tokenname) |
| `PVE_TOKEN_SECRET` | pve-cli | Proxmox API token secret |
| `ABS_URL` | abs-cli | Audiobookshelf URL |
| `ABS_TOKEN` | abs-cli | Audiobookshelf API token |
| `SONARR_URL` | sonarr-cli | Sonarr URL |
| `SONARR_API_KEY` | sonarr-cli | Sonarr API key |
| `RADARR_URL` | radarr-cli | Radarr URL |
| `RADARR_API_KEY` | radarr-cli | Radarr API key |

### Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--url` | | Server URL (overrides env var) |
| `--token` | | API token (overrides env var) |
| `--insecure` | `-k` | Skip TLS certificate verification |
| `--help` | `-h` | Show help |
| `--version` | `-v`/`-V` | Show version |

## portainer-cli

### Stacks

```bash
# List all stacks
portainer-cli stacks list

# Show stack details (includes compose file)
portainer-cli stacks show 1
```

### Endpoints

```bash
# List all endpoints
portainer-cli endpoints list

# Show endpoint details
portainer-cli endpoints show 1
```

## nproxy-cli

### Login

Get a token by authenticating with email/password:

```bash
nproxy-cli login
# Email: admin@example.com
# Password: ****
# eyJhbGciOiJS...
```

Save the output to `NPROXY_TOKEN`.

### Hosts

```bash
# List all proxy hosts
nproxy-cli hosts list

# Show proxy host details
nproxy-cli hosts show 1
```

### Certificates

```bash
# List all certificates
nproxy-cli certificates list
nproxy-cli certs list  # alias

# Show certificate details
nproxy-cli certificates show 1
```

## trans-cli

### List Torrents

```bash
# List all torrents
trans-cli list

# Filter by status
trans-cli downloading
trans-cli seeding
trans-cli stopped
```

### Show Torrent Details

```bash
trans-cli show 1
```

### Add Torrents

```bash
# Add magnet link
trans-cli add "magnet:?xt=urn:btih:..."

# Add .torrent file
trans-cli add /path/to/file.torrent
```

### Control Torrents

```bash
trans-cli start 1
trans-cli stop 1
```

## pve-cli

### List VMs and LXCs

```bash
pve-cli list
```

### Start Guest

```bash
# Start VM or LXC by ID
pve-cli start 100
```

### Stop Guest

```bash
# Stop VM or LXC by ID
pve-cli stop 100
```

## abs-cli

### Libraries

```bash
# List all libraries
abs-cli libraries list
```

### Books

```bash
# List audiobooks (uses first library by default)
abs-cli books list
abs-cli books list --library <id> --limit 100

# Show audiobook details
abs-cli books show <id>
```

### Progress

```bash
# Show listening progress (continue listening)
abs-cli progress
```

### Search

```bash
# Search audiobooks
abs-cli search "primal hunter"
abs-cli search --library <id> "author name"
```

### Scan

```bash
# Trigger library scan
abs-cli scan                    # scan all libraries
abs-cli scan --library <id>     # scan specific library
```

### Getting Your API Token

1. Log into Audiobookshelf web UI as admin
2. Go to Settings → Users
3. Click on your account
4. Copy the API token

## sonarr-cli

### Series

```bash
# List all series
sonarr-cli series list

# Show series details
sonarr-cli series show <id>
```

### Calendar

```bash
# Show upcoming episodes (default: 7 days)
sonarr-cli calendar
sonarr-cli calendar --days 14
```

### Queue

```bash
# Show download queue
sonarr-cli queue
```

### Wanted

```bash
# Show missing episodes
sonarr-cli wanted
sonarr-cli wanted --limit 50
```

### Search

```bash
# Search for new series
sonarr-cli search "breaking bad"
```

## radarr-cli

### Movies

```bash
# List all movies
radarr-cli movies list

# Show movie details
radarr-cli movies show <id>
```

### Calendar

```bash
# Show upcoming releases (default: 30 days)
radarr-cli calendar
radarr-cli calendar --days 60
```

### Queue

```bash
# Show download queue
radarr-cli queue
```

### Wanted

```bash
# Show missing movies
radarr-cli wanted
radarr-cli wanted --limit 50
```

### Search

```bash
# Search for new movies
radarr-cli search "dune"
```

## Output Format

All output is YAML. Errors go to stderr:

```yaml
error:
  code: AUTH_FAILED
  message: Invalid or expired token
```

Exit codes: 1=config, 2=auth, 3=not found, 4=network, 5=api error

## Shell Completions

### Bash

```bash
# Add to ~/.bashrc
source <(portainer-cli completion bash)
source <(nproxy-cli completion bash)
source <(trans-cli completion bash)
source <(pve-cli completion bash)
source <(abs-cli completion bash)
source <(sonarr-cli completion bash)
source <(radarr-cli completion bash)
```

### Zsh

```bash
# Add to ~/.zshrc
source <(portainer-cli completion zsh)
source <(nproxy-cli completion zsh)
source <(trans-cli completion zsh)
source <(pve-cli completion zsh)
source <(abs-cli completion zsh)
source <(sonarr-cli completion zsh)
source <(radarr-cli completion zsh)
```

If you get "command not found: compdef", add before the source lines:
```bash
autoload -Uz compinit && compinit
```

## Uninstall

```bash
rm ~/.local/bin/{portainer-cli,nproxy-cli,trans-cli,pve-cli,abs-cli,sonarr-cli,radarr-cli}
```
