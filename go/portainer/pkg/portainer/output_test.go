package portainer

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPrintYAMLFormat(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := StackListItem{
		ID:         1,
		Name:       "test",
		Type:       "compose",
		Status:     "active",
		EndpointID: 2,
	}
	err := PrintYAML(data)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("PrintYAML returned error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "name: test") {
		t.Errorf("output missing 'name: test', got: %s", output)
	}
	if !strings.Contains(output, "type: compose") {
		t.Errorf("output missing 'type: compose', got: %s", output)
	}
}

func TestPrintErrorFormat(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := ConfigError("test error message")
	PrintError(err)

	w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "code: CONFIG_ERROR") {
		t.Errorf("output missing 'code: CONFIG_ERROR', got: %s", output)
	}
	if !strings.Contains(output, "message: test error message") {
		t.Errorf("output missing message, got: %s", output)
	}
}
