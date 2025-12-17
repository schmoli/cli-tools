package keycloak

import "testing"

func TestConfigError(t *testing.T) {
	err := ConfigError("test message")
	if err.Code != ErrConfig {
		t.Errorf("expected ErrConfig, got %s", err.Code)
	}
	if err.Error() != "test message" {
		t.Errorf("expected 'test message', got %s", err.Error())
	}
	if err.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", err.ExitCode())
	}
}

func TestAuthError(t *testing.T) {
	err := AuthError("auth failed")
	if err.Code != ErrAuth {
		t.Errorf("expected ErrAuth, got %s", err.Code)
	}
	if err.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got %d", err.ExitCode())
	}
}
