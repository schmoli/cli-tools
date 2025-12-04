package nproxy

import "testing"

func TestNewClientTrimsTrailingSlash(t *testing.T) {
	client := NewClient("https://example.com/", "token", false)
	if client.baseURL != "https://example.com" {
		t.Errorf("baseURL = %s, want https://example.com", client.baseURL)
	}
}

func TestNewClientWithInsecure(t *testing.T) {
	client := NewClient("https://example.com", "token", true)
	if client.httpClient.Transport == nil {
		t.Error("Transport should not be nil when insecure=true")
	}
}

func TestNewClientCreation(t *testing.T) {
	client := NewClient("https://example.com", "token", false)
	if client == nil {
		t.Error("Client should not be nil")
	}
	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
}

func TestClientStoresToken(t *testing.T) {
	client := NewClient("https://example.com", "mytoken", false)
	if client.token != "mytoken" {
		t.Errorf("token = %s, want mytoken", client.token)
	}
}
