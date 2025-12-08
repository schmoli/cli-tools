package pve

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGetNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api2/json/nodes" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") == "" {
			t.Error("missing Authorization header")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"node":"pve","status":"online"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "user@pam!token", "secret", false)
	node, err := client.GetNode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node != "pve" {
		t.Errorf("got node %q, want pve", node)
	}
}

func TestClientAuthHeader(t *testing.T) {
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"node":"pve","status":"online"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "user@pam!mytoken", "abc123", false)
	client.GetNode()

	expected := "PVEAPIToken=user@pam!mytoken=abc123"
	if gotAuth != expected {
		t.Errorf("got auth %q, want %q", gotAuth, expected)
	}
}
