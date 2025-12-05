# README Consolidation & Completion Instructions

## Overview

Consolidate tool-specific docs into single README and add post-install shell completion instructions.

## Changes

### 1. install.sh - Post-Install Completion Instructions

After PATH notice, detect shell and print completion setup:

```bash
echo ""
echo "Shell completions (optional):"
SHELL_NAME=$(basename "$SHELL")
case "$SHELL_NAME" in
  bash)
    echo "  Add to ~/.bashrc:"
    echo "    source <(portainer-cli completion bash)"
    echo "    source <(nproxy-cli completion bash)"
    ;;
  zsh)
    echo "  Add to ~/.zshrc:"
    echo "    source <(portainer-cli completion zsh)"
    echo "    source <(nproxy-cli completion zsh)"
    ;;
  *)
    echo "  Run: portainer-cli completion --help"
    ;;
esac
```

### 2. README.md - Consolidated Structure

```markdown
# cli-tools
[badges]
> CLI tools for Portainer and nginx-proxy-manager APIs

## Install
[curl command]

## Configuration

| Variable | Tools | Description |
|----------|-------|-------------|
| PORTAINER_URL | portainer-cli | Portainer server URL |
| PORTAINER_TOKEN | portainer-cli | Portainer API token |
| NPROXY_URL | nproxy-cli | nginx-proxy-manager URL |
| NPROXY_TOKEN | nproxy-cli | nginx-proxy-manager JWT |

### Common Flags
| Flag | Short | Description |
|------|-------|-------------|
| --url | | Server URL |
| --token | | API token |
| --insecure | -k | Skip TLS verify |

## portainer-cli

### Stacks
[examples]

### Endpoints
[examples]

## nproxy-cli

### Login
[examples]

### Hosts
[examples]

### Certificates
[examples]

## Output Format
[shared - YAML format, error structure, exit codes]

## Shell Completions
[bash/zsh instructions]

## Uninstall
[rm command]
```

### 3. Delete Separate Tool Docs

- Remove `docs/portainer.md`
- Remove `docs/nproxy.md`
- Keep `docs/plans/` for design documents

## Files Changed

- `install.sh` - add completion instructions
- `README.md` - consolidate all tool docs
- `docs/portainer.md` - delete
- `docs/nproxy.md` - delete
