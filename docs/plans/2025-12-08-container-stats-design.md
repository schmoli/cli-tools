# Container Stats Feature Design

## Overview

Add container listing to portainer-cli with two commands:
- `containers list` - list all containers across all endpoints
- `stacks containers <id>` - list containers for a specific stack

## Commands

### `containers list [--endpoint <id>]`

Lists containers across all endpoints by default. Optional `--endpoint` flag to filter to one endpoint.

Sorted by: endpoint ID, then container name.

### `stacks containers <stack-id>`

Lists containers belonging to a specific stack. Finds endpoint from stack data, filters by `com.docker.compose.project` label.

## Output Fields

```yaml
- id: abc123def456
  name: myapp_web_1
  image: nginx:latest
  state: running
  stack: myapp
  endpoint: 1
  ports:
    - host: 8080
      container: 80
      protocol: tcp
  created: 2024-01-15T10:30:00Z
  health: healthy
```

| Field | Description |
|-------|-------------|
| id | Short container ID (12 chars) |
| name | Container name (without leading /) |
| image | Image name with tag |
| state | running, exited, paused, restarting, dead |
| stack | From `com.docker.compose.project` label, empty if none |
| endpoint | Endpoint ID |
| ports | Array of port mappings |
| created | ISO 8601 timestamp |
| health | healthy, unhealthy, starting, none |

## API Implementation

**Endpoints used:**
1. `GET /api/endpoints` - list endpoints (existing)
2. `GET /api/endpoints/{id}/docker/containers/json?all=true` - Docker API proxy

**New client method:**
```go
func (c *Client) ListContainers(endpointID int64) ([]APIContainer, error)
```

**`containers list` flow:**
1. Fetch all endpoints (or use specified one)
2. For each endpoint, fetch containers via Docker proxy
3. Aggregate, transform, sort, output

**`stacks containers` flow:**
1. Fetch stack by ID to get endpointId
2. Fetch containers for that endpoint
3. Filter by stack label
4. Transform and output

## File Changes

**New file:**
- `cmd/portainer-cli/containers.go` - container commands

**Modified files:**
- `pkg/portainer/client.go` - add `ListContainers()` method
- `pkg/portainer/models.go` - add container types

## Models

```go
// API response (Docker API proxy)
type APIContainer struct {
    Id      string            `json:"Id"`
    Names   []string          `json:"Names"`
    Image   string            `json:"Image"`
    State   string            `json:"State"`
    Status  string            `json:"Status"`
    Created int64             `json:"Created"`
    Ports   []APIPort         `json:"Ports"`
    Labels  map[string]string `json:"Labels"`
}

type APIPort struct {
    IP          string `json:"IP"`
    PrivatePort int    `json:"PrivatePort"`
    PublicPort  int    `json:"PublicPort"`
    Type        string `json:"Type"`
}

// Output types
type ContainerListItem struct {
    ID       string `yaml:"id"`
    Name     string `yaml:"name"`
    Image    string `yaml:"image"`
    State    string `yaml:"state"`
    Stack    string `yaml:"stack"`
    Endpoint int64  `yaml:"endpoint"`
    Ports    []Port `yaml:"ports"`
    Created  string `yaml:"created"`
    Health   string `yaml:"health"`
}

type Port struct {
    Host      int    `yaml:"host,omitempty"`
    Container int    `yaml:"container"`
    Protocol  string `yaml:"protocol"`
}

type ContainerList []ContainerListItem
```
