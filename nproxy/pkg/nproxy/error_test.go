package nproxy

import "testing"

func TestExitCodes(t *testing.T) {
	tests := []struct {
		err      *NproxyError
		expected int
	}{
		{ConfigError("test"), 1},
		{AuthError("test"), 2},
		{NotFoundError("test"), 3},
		{NetworkError("test"), 4},
		{APIError("test"), 5},
	}

	for _, tt := range tests {
		if got := tt.err.ExitCode(); got != tt.expected {
			t.Errorf("ExitCode() = %d, want %d for %s", got, tt.expected, tt.err.Code)
		}
	}
}

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		err      *NproxyError
		expected ErrorCode
	}{
		{ConfigError("test"), ErrConfig},
		{AuthError("test"), ErrAuth},
		{NotFoundError("test"), ErrNotFound},
		{NetworkError("test"), ErrNetwork},
		{APIError("test"), ErrAPI},
	}

	for _, tt := range tests {
		if tt.err.Code != tt.expected {
			t.Errorf("Code = %s, want %s", tt.err.Code, tt.expected)
		}
	}
}

func TestErrorMessage(t *testing.T) {
	err := ConfigError("missing URL")
	if err.Error() != "missing URL" {
		t.Errorf("Error() = %s, want 'missing URL'", err.Error())
	}
}
