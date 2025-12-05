# portainer-cli

CLI tool for Portainer API - backup and viewing operations.

## Configuration

Set credentials via environment variables (recommended):
```bash
export PORTAINER_URL=https://portainer.example.com
export PORTAINER_TOKEN=ptr_xxxxxxxxxxxx
```

Or pass as flags:
```bash
portainer-cli --url https://portainer.example.com --token ptr_xxx stacks list
```

### Options

| Flag | Short | Description |
|------|-------|-------------|
| `--url` | | Portainer server URL |
| `--token` | | API token |
| `--insecure` | `-k` | Skip TLS certificate verification |
| `--help` | `-h` | Show help |
| `--version` | `-v`/`-V` | Show version |

## Usage

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

## Output Format

All output is YAML. Errors go to stderr with structured format:

```yaml
error:
  code: AUTH_FAILED
  message: Invalid or expired token
```

Exit codes: 1=config, 2=auth, 3=not found, 4=network, 5=api error

## Shell Completion (Go only)

The Go implementation supports shell completion via Cobra.

### Bash

```bash
# Add to ~/.bashrc
echo 'source <(portainer-cli completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Zsh

```bash
# Add to ~/.zshrc
echo 'source <(portainer-cli completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

If you get "command not found: compdef", add this before the source line:
```bash
autoload -Uz compinit && compinit
```
