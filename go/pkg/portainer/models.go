package portainer

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
		return "up"
	case 2:
		return "down"
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
