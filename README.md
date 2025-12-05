# cli-tools

[![Release](https://img.shields.io/github/v/release/schmoli/cli-tools)](https://github.com/schmoli/cli-tools/releases)
[![License](https://img.shields.io/github/license/schmoli/cli-tools)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/schmoli/cli-tools/release-please.yml?branch=main)](https://github.com/schmoli/cli-tools/actions)

> CLI tools for Portainer and nginx-proxy-manager APIs

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/schmoli/cli-tools/main/install.sh | bash
```

Installs to `~/.local/bin`.

## Tools

| Tool | Description |
|------|-------------|
| [portainer-cli](docs/portainer.md) | Portainer API - stacks, endpoints |
| [nproxy-cli](docs/nproxy.md) | nginx-proxy-manager API - hosts, certs |

## Uninstall

```bash
rm ~/.local/bin/{portainer-cli,nproxy-cli}
```

## Development

```bash
./build.sh
```
