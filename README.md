# cli-tools

[![Release](https://img.shields.io/github/v/release/schmoli/cli-tools)](https://github.com/schmoli/cli-tools/releases)
[![License](https://img.shields.io/github/license/schmoli/cli-tools)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/schmoli/cli-tools/release-please.yml?branch=main)](https://github.com/schmoli/cli-tools/actions)

> CLI tools for Portainer, nginx-proxy-manager, Transmission, Proxmox VE, and Keycloak APIs

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
| `KEYCLOAK_URL` | keycloak-cli | Keycloak server URL |
| `KEYCLOAK_REALM` | keycloak-cli | Auth realm (usually "master") |
| `KEYCLOAK_CLIENT_ID` | keycloak-cli | Service account client ID |
| `KEYCLOAK_CLIENT_SECRET` | keycloak-cli | Service account client secret |
| `KEYCLOAK_TARGET_REALM` | keycloak-cli | Target realm for queries |

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

## keycloak-cli

### Realms

```bash
# List all realms
keycloak-cli realms list

# Show realm details
keycloak-cli realms get master
```

### Users

```bash
# List users in a realm
keycloak-cli users list --target-realm myrealm

# Show user details
keycloak-cli users get <user-id> --target-realm myrealm

# Show user sessions
keycloak-cli users sessions <user-id> --target-realm myrealm
```

### Clients

```bash
# List clients
keycloak-cli clients list --target-realm myrealm

# Show client details
keycloak-cli clients get <client-id> --target-realm myrealm

# Show client sessions
keycloak-cli clients sessions <client-uuid> --target-realm myrealm
```

### Roles

```bash
# List realm roles
keycloak-cli roles list --target-realm myrealm

# List client roles
keycloak-cli roles list --target-realm myrealm --client <client-uuid>

# Show role details
keycloak-cli roles get admin --target-realm myrealm
```

### Groups

```bash
# List groups
keycloak-cli groups list --target-realm myrealm

# Show group details
keycloak-cli groups get <group-id> --target-realm myrealm

# List group members
keycloak-cli groups members <group-id> --target-realm myrealm
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
source <(keycloak-cli completion bash)
```

### Zsh

```bash
# Add to ~/.zshrc
source <(portainer-cli completion zsh)
source <(nproxy-cli completion zsh)
source <(trans-cli completion zsh)
source <(pve-cli completion zsh)
source <(keycloak-cli completion zsh)
```

If you get "command not found: compdef", add before the source lines:
```bash
autoload -Uz compinit && compinit
```

## Uninstall

```bash
rm ~/.local/bin/{portainer-cli,nproxy-cli,trans-cli,pve-cli,keycloak-cli}
```
