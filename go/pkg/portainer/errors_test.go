package portainer

import "testing"

func TestErrorExitCodes(t *testing.T) {
	tests := []struct {
		name     string
		err      *PortainerError
		expected int
	}{
		{"config", ConfigError("test"), 1},
		{"auth", AuthError("test"), 2},
		{"not_found", NotFoundError("test"), 3},
		{"network", NetworkError("test"), 4},
		{"api", APIError("test"), 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.ExitCode()
			if got != tt.expected {
				t.Errorf("ExitCode() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		err      *PortainerError
		expected ErrorCode
	}{
		{"config", ConfigError("test"), ErrConfig},
		{"auth", AuthError("test"), ErrAuth},
		{"not_found", NotFoundError("test"), ErrNotFound},
		{"network", NetworkError("test"), ErrNetwork},
		{"api", APIError("test"), ErrAPI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.expected {
				t.Errorf("Code = %q, want %q", tt.err.Code, tt.expected)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	err := ConfigError("missing url")
	if err.Error() != "missing url" {
		t.Errorf("Error() = %q, want %q", err.Error(), "missing url")
	}
}
