# portainer-cli

CLI tool for Portainer API - backup and viewing operations.

## Implementations

Two parallel implementations for comparison:

- `rust/` - Rust implementation using clap + reqwest
- `go/` - Go implementation using cobra + net/http

## Usage

```bash
# Set credentials
export PORTAINER_URL=https://portainer.example.com
export PORTAINER_TOKEN=ptr_xxxxxxxxxxxx

# List stacks
portainer-cli stacks list

# Show stack details
portainer-cli stacks show 1

# List endpoints
portainer-cli endpoints list

# Show endpoint details
portainer-cli endpoints show 1
```

## Build

### Rust
```bash
cd rust && cargo build --release
# Binary: rust/target/release/portainer-cli
```

### Go
```bash
cd go && go build -o portainer-cli ./cmd/portainer-cli
# Binary: go/portainer-cli
```

## Output Format

All output is YAML. Errors go to stderr.

See `docs/plans/2025-12-03-portainer-cli-requirements.md` for full specification.
