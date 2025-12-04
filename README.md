# cli-tools

Monorepo for CLI tools that fetch/transform data from various sources.

## Structure

```
cli-tools/
├── go/
│   ├── common/         # Shared Go code
│   └── portainer/      # portainer-cli
└── rust/
    ├── common/         # Shared Rust code
    └── portainer/      # portainer-cli
```

## Tools

### portainer-cli

CLI for Portainer API - backup and viewing operations.

| Feature | Go | Rust |
|---------|-----|------|
| Binary size | ~9 MB | ~4 MB |
| Shell completions | bash, zsh, fish, powershell | - |

See [portainer documentation](docs/portainer.md) for usage details.

## Build

```bash
./build.sh
```

Builds all tools. Binaries are created in each tool's directory:
- `go/portainer/portainer-cli`
- `rust/portainer/portainer-cli`
