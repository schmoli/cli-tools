# portainer-cli Requirements

## Overview

`portainer-cli` - A read-only CLI tool for interacting with the Portainer API, focused on backup/viewing use cases.

## Scope

### In Scope (v1)
- List and show stacks
- List and show endpoints
- YAML output
- Auth via flags and env vars

### Out of Scope (v1)
- Create/update/delete operations
- Table/human-readable output formats
- Config file support
- Field selection flags
- Retry logic

### Future Considerations
- Part of a monorepo with shared library for future CLI tools
- Additional commands (create, update, delete)
- Configurable field selection via env vars

---

## CLI Interface

### Commands

```
portainer-cli stacks list              # List all stacks
portainer-cli stacks show <id>         # Show single stack by ID
portainer-cli endpoints list           # List all endpoints
portainer-cli endpoints show <id>      # Show single endpoint by ID
```

### Global Flags

```
--url <url>       Portainer instance URL (fallback: PORTAINER_URL)
--token <token>   API token (fallback: PORTAINER_TOKEN)
--help            Show help
--version         Show version
```

---

## Output Format

All output is YAML.

### Success Output (stdout)

**stacks list:**
```yaml
stacks:
  - id: 1
    name: monitoring
    type: compose
    status: active
    endpointId: 1
  - id: 2
    name: webapp
    type: compose
    status: active
    endpointId: 1
```

**stacks show \<id\>:**
```yaml
id: 1
name: monitoring
type: compose
status: active
endpointId: 1
env:
  - name: GRAFANA_PORT
    value: "3000"
stackFile: |
  version: "3"
  services:
    grafana:
      image: grafana/grafana
      ports:
        - "${GRAFANA_PORT}:3000"
```

**endpoints list:**
```yaml
endpoints:
  - id: 1
    name: local
    type: docker
    status: active
    url: unix:///var/run/docker.sock
  - id: 2
    name: production
    type: docker
    status: active
    url: tcp://prod.example.com:2375
```

**endpoints show \<id\>:**
```yaml
id: 1
name: local
type: docker
status: active
url: unix:///var/run/docker.sock
```

### Error Output (stderr)

```yaml
error:
  code: NOT_FOUND
  message: Stack with ID 99 not found
```

---

## Data Fields

### Stack Fields (curated)

| Field | Description |
|-------|-------------|
| id | Stack ID |
| name | Stack name |
| type | `swarm`, `compose`, or `kubernetes` |
| status | `active` or `inactive` |
| endpointId | Environment/endpoint ID |
| env | List of environment variables (name/value pairs) |
| stackFile | Docker compose file content |

**Type mapping:** 1=swarm, 2=compose, 3=kubernetes
**Status mapping:** 1=active, 2=inactive

### Endpoint Fields (curated)

| Field | Description |
|-------|-------------|
| id | Endpoint ID |
| name | Endpoint name |
| type | Endpoint type (docker, kubernetes, etc.) |
| status | Endpoint status |
| url | Endpoint URL |

---

## Authentication

### Methods
- API token-based (Portainer access tokens)
- Stateless requests (no session/login flow)

### Precedence (highest to lowest)
1. `--url` and `--token` flags
2. `PORTAINER_URL` and `PORTAINER_TOKEN` env vars

### Validation
- Both URL and token required - fail fast with clear error if missing
- URL validated as valid HTTP(S) URL before first request

### Example Usage

```bash
# Via env vars (typical for scripts/cron)
export PORTAINER_URL=https://portainer.example.com
export PORTAINER_TOKEN=ptr_xxxxxxxxxxxx
portainer-cli stacks list

# Via flags (one-off use)
portainer-cli --url https://portainer.example.com --token ptr_xxx stacks list
```

---

## Error Handling

### Exit Codes

| Exit Code | Category | Examples |
|-----------|----------|----------|
| 0 | Success | Request completed |
| 1 | Config error | Missing URL/token, invalid URL |
| 2 | Auth error | Invalid token, expired token |
| 3 | Not found | Stack/endpoint ID doesn't exist |
| 4 | Network error | Connection refused, timeout |
| 5 | API error | Unexpected Portainer response |

### Behavior
- Fail fast on first error (no partial output)
- No retries in v1
- Stderr for errors, stdout stays clean for piping

---

## Timeouts

- Default request timeout: 10 seconds
- No configurable timeout flag in v1

---

## Binary

- Name: `portainer-cli`
- Single executable
- Language: TBD (Rust or Go, decided in implementation planning)
