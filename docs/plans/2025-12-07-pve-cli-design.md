# pve-cli Design

CLI for Proxmox VE API - list, start, and stop VMs/LXCs.

## Configuration

**Environment Variables:**
| Variable | Description |
|----------|-------------|
| `PVE_URL` | Proxmox API URL (e.g., https://pve.local:8006) |
| `PVE_TOKEN_ID` | Token ID (format: user@realm!tokenname) |
| `PVE_TOKEN_SECRET` | Token secret/value |

**Flags:**
| Flag | Description |
|------|-------------|
| `--url` | Override PVE_URL |
| `--token-id` | Override PVE_TOKEN_ID |
| `--token` | Override PVE_TOKEN_SECRET |
| `--insecure`, `-k` | Skip TLS verification |

**Auth header:** `Authorization: PVEAPIToken=USER@REALM!TOKENID=SECRET`

**Node discovery:** Auto-discover from `/api2/json/nodes`, cache first node.

## Commands

```bash
pve-cli list              # List all VMs and LXCs
pve-cli start <vmid>      # Start a VM or LXC
pve-cli stop <vmid>       # Stop a VM or LXC
```

## API Mapping

| Command | Proxmox API |
|---------|-------------|
| list | `GET /nodes/{node}/qemu` + `GET /nodes/{node}/lxc` |
| start (VM) | `POST /nodes/{node}/qemu/{vmid}/status/start` |
| start (LXC) | `POST /nodes/{node}/lxc/{vmid}/status/start` |
| stop (VM) | `POST /nodes/{node}/qemu/{vmid}/status/shutdown` |
| stop (LXC) | `POST /nodes/{node}/lxc/{vmid}/status/shutdown` |

**Type auto-detection:** Try QEMU first, then LXC if not found.

**IP fetching:**
- VMs: `/nodes/{node}/qemu/{vmid}/agent/network-get-interfaces`
- LXCs: `/nodes/{node}/lxc/{vmid}/interfaces`
- Fallback to "N/A" if unavailable

## Output

**List:**
```yaml
- vmid: 100
  name: ubuntu-server
  type: qemu
  status: running
  cpu: 4
  memory: 8192
  uptime: 3d 4h 12m
  ip: 192.168.1.100
```

**Start/stop:**
```yaml
vmid: 100
name: ubuntu-server
action: started
```

## Project Structure

```
pve/
├── cmd/pve-cli/
│   └── main.go
├── pkg/pve/
│   ├── client.go
│   ├── models.go
│   ├── errors.go
│   └── output.go
└── go.mod
```

## Error Handling

Same exit codes as other tools: 1=config, 2=auth, 3=not found, 4=network, 5=api error
