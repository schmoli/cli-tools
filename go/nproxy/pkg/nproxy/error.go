package nproxy

import "fmt"

type ErrorCode string

const (
	ErrConfig   ErrorCode = "CONFIG_ERROR"
	ErrAuth     ErrorCode = "AUTH_FAILED"
	ErrNotFound ErrorCode = "NOT_FOUND"
	ErrNetwork  ErrorCode = "NETWORK_ERROR"
	ErrAPI      ErrorCode = "API_ERROR"
)

type NproxyError struct {
	Code    ErrorCode
	Message string
}

func (e *NproxyError) Error() string {
	return e.Message
}

func (e *NproxyError) ExitCode() int {
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

func ConfigError(msg string) *NproxyError {
	return &NproxyError{Code: ErrConfig, Message: msg}
}

func AuthError(msg string) *NproxyError {
	return &NproxyError{Code: ErrAuth, Message: msg}
}

func NotFoundError(msg string) *NproxyError {
	return &NproxyError{Code: ErrNotFound, Message: msg}
}

func NetworkError(msg string) *NproxyError {
	return &NproxyError{Code: ErrNetwork, Message: msg}
}

func APIError(msg string) *NproxyError {
	return &NproxyError{Code: ErrAPI, Message: fmt.Sprintf("API error: %s", msg)}
}
