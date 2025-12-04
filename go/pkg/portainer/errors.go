package portainer

import "fmt"

type ErrorCode string

const (
	ErrConfig  ErrorCode = "CONFIG_ERROR"
	ErrAuth    ErrorCode = "AUTH_FAILED"
	ErrNotFound ErrorCode = "NOT_FOUND"
	ErrNetwork ErrorCode = "NETWORK_ERROR"
	ErrAPI     ErrorCode = "API_ERROR"
)

type PortainerError struct {
	Code    ErrorCode
	Message string
}

func (e *PortainerError) Error() string {
	return e.Message
}

func (e *PortainerError) ExitCode() int {
	switch e.Code {
	case ErrConfig:
		return 1
	case ErrAuth:
		return 2
	case ErrNotFound:
		return 3
	case ErrNetwork:
		return 4
	case ErrAPI:
		return 5
	default:
		return 1
	}
}

func ConfigError(msg string) *PortainerError {
	return &PortainerError{Code: ErrConfig, Message: msg}
}

func AuthError(msg string) *PortainerError {
	return &PortainerError{Code: ErrAuth, Message: msg}
}

func NotFoundError(msg string) *PortainerError {
	return &PortainerError{Code: ErrNotFound, Message: msg}
}

func NetworkError(msg string) *PortainerError {
	return &PortainerError{Code: ErrNetwork, Message: msg}
}

func APIError(msg string) *PortainerError {
	return &PortainerError{Code: ErrAPI, Message: fmt.Sprintf("API error: %s", msg)}
}
