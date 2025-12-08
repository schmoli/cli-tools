package pve

import (
	"bytes"
	"testing"
)

func TestPrintYAML(t *testing.T) {
	var buf bytes.Buffer
	data := map[string]string{"name": "test"}

	err := PrintYAMLTo(&buf, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "name: test\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}
