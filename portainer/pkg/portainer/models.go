package portainer

import (
	"strings"
	"time"
)

// API response types (match Portainer JSON)
type APIStack struct {
	ID         int64       `json:"Id"`
	Name       string      `json:"Name"`
	Type       int         `json:"Type"`
	Status     int         `json:"Status"`
	EndpointID int64       `json:"EndpointId"`
	Env        []APIEnvVar `json:"Env"`
}

type APIEnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type APIStackFile struct {
	StackFileContent string `json:"StackFileContent"`
}

type APIEndpoint struct {
	ID     int64  `json:"Id"`
	Name   string `json:"Name"`
	Type   int    `json:"Type"`
	Status int    `json:"Status"`
	URL    string `json:"URL"`
}

// Output types (curated, YAML output)
type Stack struct {
	ID         int64       `yaml:"id"`
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Status     string      `yaml:"status"`
	EndpointID int64       `yaml:"endpointId"`
	Env        []APIEnvVar `yaml:"env,omitempty"`
	StackFile  string      `yaml:"stackFile,omitempty"`
}

type StackList struct {
	Stacks []StackListItem `yaml:"stacks"`
}

type StackListItem struct {
	ID         int64  `yaml:"id"`
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Status     string `yaml:"status"`
	EndpointID int64  `yaml:"endpointId"`
}

type Endpoint struct {
	ID     int64  `yaml:"id"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Status string `yaml:"status"`
	URL    string `yaml:"url"`
}

type EndpointList struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Mapping functions
func (s *APIStack) TypeLabel() string {
	switch s.Type {
	case 1:
		return "swarm"
	case 2:
		return "compose"
	case 3:
		return "kubernetes"
	default:
		return "unknown"
	}
}

func (s *APIStack) StatusLabel() string {
	switch s.Status {
	case 1:
		return "active"
	case 2:
		return "inactive"
	default:
		return "unknown"
	}
}

func (s *APIStack) ToListItem() StackListItem {
	return StackListItem{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.TypeLabel(),
		Status:     s.StatusLabel(),
		EndpointID: s.EndpointID,
	}
}

func (s *APIStack) ToStack(stackFile string) Stack {
	return Stack{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.TypeLabel(),
		Status:     s.StatusLabel(),
		EndpointID: s.EndpointID,
		Env:        s.Env,
		StackFile:  stackFile,
	}
}

func (e *APIEndpoint) TypeLabel() string {
	switch e.Type {
	case 1:
		return "docker"
	case 2:
		return "agent"
	case 3:
		return "azure"
	case 4:
		return "edge-agent"
	case 5:
		return "kubernetes"
	default:
		return "unknown"
	}
}

func (e *APIEndpoint) StatusLabel() string {
	switch e.Status {
	case 1:
		return "active"
	case 2:
		return "inactive"
	default:
		return "unknown"
	}
}

func (e *APIEndpoint) ToEndpoint() Endpoint {
	return Endpoint{
		ID:     e.ID,
		Name:   e.Name,
		Type:   e.TypeLabel(),
		Status: e.StatusLabel(),
		URL:    e.URL,
	}
}

// Container API types (Docker API proxy response)
type APIContainer struct {
	ID      string            `json:"Id"`
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

// Container output types
type ContainerListItem struct {
	ID       string `yaml:"id"`
	Name     string `yaml:"name"`
	Image    string `yaml:"image"`
	State    string `yaml:"state"`
	Stack    string `yaml:"stack"`
	Endpoint int64  `yaml:"endpoint"`
	Ports    []Port `yaml:"ports,omitempty"`
	Created  string `yaml:"created"`
	Health   string `yaml:"health"`
}

type Port struct {
	Host      int    `yaml:"host,omitempty"`
	Container int    `yaml:"container"`
	Protocol  string `yaml:"protocol"`
}

type ContainerList []ContainerListItem

// Container conversion methods
func (c *APIContainer) ToListItem(endpointID int64) ContainerListItem {
	name := ""
	if len(c.Names) > 0 {
		name = c.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}
	}

	// Short ID (12 chars like Docker)
	id := c.ID
	if len(id) > 12 {
		id = id[:12]
	}

	// Stack from compose label
	stack := c.Labels["com.docker.compose.project"]

	// Health from Status field (e.g., "Up 2 hours (healthy)")
	health := "none"
	statusLower := strings.ToLower(c.Status)
	if strings.Contains(statusLower, "(healthy)") {
		health = "healthy"
	} else if strings.Contains(statusLower, "(unhealthy)") {
		health = "unhealthy"
	} else if strings.Contains(statusLower, "(health: starting)") {
		health = "starting"
	}

	// Convert ports
	var ports []Port
	for _, p := range c.Ports {
		ports = append(ports, Port{
			Host:      p.PublicPort,
			Container: p.PrivatePort,
			Protocol:  p.Type,
		})
	}

	// Format created time
	created := formatUnixTime(c.Created)

	return ContainerListItem{
		ID:       id,
		Name:     name,
		Image:    c.Image,
		State:    c.State,
		Stack:    stack,
		Endpoint: endpointID,
		Ports:    ports,
		Created:  created,
		Health:   health,
	}
}

func formatUnixTime(unix int64) string {
	if unix == 0 {
		return ""
	}
	return time.Unix(unix, 0).UTC().Format(time.RFC3339)
}
