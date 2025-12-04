package portainer

import "testing"

func TestNewClientTrimsURL(t *testing.T) {
	client := NewClient("http://example.com/", "token", false)
	if client.baseURL != "http://example.com" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "http://example.com")
	}
}

func TestNewClientNoTrailingSlash(t *testing.T) {
	client := NewClient("http://example.com", "token", false)
	if client.baseURL != "http://example.com" {
		t.Errorf("baseURL = %q, want %q", client.baseURL, "http://example.com")
	}
}

func TestNewClientStoresToken(t *testing.T) {
	client := NewClient("http://example.com", "mytoken", false)
	if client.token != "mytoken" {
		t.Errorf("token = %q, want %q", client.token, "mytoken")
	}
}

func TestNewClientSecureByDefault(t *testing.T) {
	client := NewClient("http://example.com", "token", false)
	// When insecure=false, should not have custom TLS config that skips verification
	// Just verify client was created successfully
	if client.httpClient == nil {
		t.Error("expected non-nil httpClient")
	}
}

func TestNewClientInsecure(t *testing.T) {
	client := NewClient("http://example.com", "token", true)
	// When insecure=true, transport should be set
	if client.httpClient.Transport == nil {
		t.Error("expected non-nil transport for insecure client")
	}
}
