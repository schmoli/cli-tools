package pve

import "testing"

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		seconds int64
		want    string
	}{
		{0, "0s"},
		{59, "59s"},
		{60, "1m"},
		{3600, "1h"},
		{3661, "1h 1m"},
		{86400, "1d"},
		{90061, "1d 1h 1m"},
	}

	for _, tt := range tests {
		got := FormatUptime(tt.seconds)
		if got != tt.want {
			t.Errorf("FormatUptime(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestAPIVMStatusLabel(t *testing.T) {
	vm := &APIVM{Status: "running"}
	if vm.StatusLabel() != "running" {
		t.Errorf("got %q, want running", vm.StatusLabel())
	}

	vm.Status = "stopped"
	if vm.StatusLabel() != "stopped" {
		t.Errorf("got %q, want stopped", vm.StatusLabel())
	}
}
