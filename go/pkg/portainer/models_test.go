package portainer

import "testing"

func TestStackTypeLabel(t *testing.T) {
	tests := []struct {
		name     string
		typeCode int
		expected string
	}{
		{"swarm", 1, "swarm"},
		{"compose", 2, "compose"},
		{"kubernetes", 3, "kubernetes"},
		{"unknown", 99, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &APIStack{Type: tt.typeCode}
			got := stack.TypeLabel()
			if got != tt.expected {
				t.Errorf("TypeLabel() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestStackStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected string
	}{
		{"active", 1, "active"},
		{"inactive", 2, "inactive"},
		{"unknown", 99, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := &APIStack{Status: tt.status}
			got := stack.StatusLabel()
			if got != tt.expected {
				t.Errorf("StatusLabel() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestEndpointTypeLabel(t *testing.T) {
	tests := []struct {
		name     string
		typeCode int
		expected string
	}{
		{"docker", 1, "docker"},
		{"agent", 2, "agent"},
		{"azure", 3, "azure"},
		{"edge-agent", 4, "edge-agent"},
		{"kubernetes", 5, "kubernetes"},
		{"unknown", 99, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &APIEndpoint{Type: tt.typeCode}
			got := ep.TypeLabel()
			if got != tt.expected {
				t.Errorf("TypeLabel() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestEndpointStatusLabel(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		expected string
	}{
		{"active", 1, "active"},
		{"inactive", 2, "inactive"},
		{"unknown", 99, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &APIEndpoint{Status: tt.status}
			got := ep.StatusLabel()
			if got != tt.expected {
				t.Errorf("StatusLabel() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestStackToListItem(t *testing.T) {
	stack := &APIStack{
		ID:         42,
		Name:       "mystack",
		Type:       2,
		Status:     1,
		EndpointID: 5,
	}

	item := stack.ToListItem()

	if item.ID != 42 {
		t.Errorf("ID = %d, want 42", item.ID)
	}
	if item.Name != "mystack" {
		t.Errorf("Name = %q, want %q", item.Name, "mystack")
	}
	if item.Type != "compose" {
		t.Errorf("Type = %q, want %q", item.Type, "compose")
	}
	if item.Status != "active" {
		t.Errorf("Status = %q, want %q", item.Status, "active")
	}
	if item.EndpointID != 5 {
		t.Errorf("EndpointID = %d, want 5", item.EndpointID)
	}
}

func TestEndpointToEndpoint(t *testing.T) {
	apiEp := &APIEndpoint{
		ID:     10,
		Name:   "prod",
		Type:   1,
		Status: 2,
		URL:    "tcp://docker:2375",
	}

	ep := apiEp.ToEndpoint()

	if ep.ID != 10 {
		t.Errorf("ID = %d, want 10", ep.ID)
	}
	if ep.Name != "prod" {
		t.Errorf("Name = %q, want %q", ep.Name, "prod")
	}
	if ep.Type != "docker" {
		t.Errorf("Type = %q, want %q", ep.Type, "docker")
	}
	if ep.Status != "inactive" {
		t.Errorf("Status = %q, want %q", ep.Status, "inactive")
	}
	if ep.URL != "tcp://docker:2375" {
		t.Errorf("URL = %q, want %q", ep.URL, "tcp://docker:2375")
	}
}
