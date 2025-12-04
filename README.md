# cli-tools

CLI tools for various APIs.

## Structure

```
cli-tools/
├── common/         # Shared Go code
├── portainer/      # portainer-cli
└── nproxy/         # nproxy-cli
```

## Tools

### portainer-cli

CLI for Portainer API. See [docs/portainer.md](docs/portainer.md).

### nproxy-cli

CLI for nginx-proxy-manager API. See [docs/nproxy.md](docs/nproxy.md).

## Build

```bash
./build.sh
```

Binaries: `bin/portainer-cli`, `bin/nproxy-cli`
