# cli-tools

CLI tools for various APIs.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/schmoli/cli-tools/main/install.sh | bash
```

Installs `portainer-cli` and `nproxy-cli` to `~/.local/bin`.

## Tools

### portainer-cli

CLI for Portainer API. See [docs/portainer.md](docs/portainer.md).

### nproxy-cli

CLI for nginx-proxy-manager API. See [docs/nproxy.md](docs/nproxy.md).

## Development

```bash
./build.sh
```

Binaries: `bin/portainer-cli`, `bin/nproxy-cli`
