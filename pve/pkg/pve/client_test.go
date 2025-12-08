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

func TestClientListGuests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api2/json/nodes":
			w.Write([]byte(`{"data":[{"node":"pve","status":"online"}]}`))
		case "/api2/json/nodes/pve/qemu":
			w.Write([]byte(`{"data":[{"vmid":100,"name":"vm1","status":"running","cpus":4,"maxmem":8589934592,"uptime":3600}]}`))
		case "/api2/json/nodes/pve/lxc":
			w.Write([]byte(`{"data":[{"vmid":101,"name":"ct1","status":"running","cpus":2,"maxmem":536870912,"uptime":7200}]}`))
		case "/api2/json/nodes/pve/qemu/100/agent/network-get-interfaces":
			w.Write([]byte(`{"data":{"result":[{"name":"eth0","ip-addresses":[{"ip-address":"192.168.1.100","ip-address-type":"ipv4"}]}]}}`))
		case "/api2/json/nodes/pve/lxc/101/interfaces":
			w.Write([]byte(`{"data":[{"name":"eth0","inet":"192.168.1.101/24"}]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "user@pam!token", "secret", false)
	guests, err := client.ListGuests()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(guests) != 2 {
		t.Fatalf("got %d guests, want 2", len(guests))
	}

	if guests[0].VMID != 100 || guests[0].Type != "qemu" {
		t.Errorf("guest 0: got vmid=%d type=%s, want vmid=100 type=qemu", guests[0].VMID, guests[0].Type)
	}
	if guests[1].VMID != 101 || guests[1].Type != "lxc" {
		t.Errorf("guest 1: got vmid=%d type=%s, want vmid=101 type=lxc", guests[1].VMID, guests[1].Type)
	}
}
