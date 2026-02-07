package abs

import "fmt"

type ErrorCode string

const (
	ErrConfig   ErrorCode = "CONFIG_ERROR"
	ErrAuth     ErrorCode = "AUTH_FAILED"
	ErrNotFound ErrorCode = "NOT_FOUND"
	ErrNetwork  ErrorCode = "NETWORK_ERROR"
	ErrAPI      ErrorCode = "API_ERROR"
)

type AbsError struct {
	Code    ErrorCode
	Message string
}

func (e *AbsError) Error() string {
	return e.Message
}

func (e *AbsError) ExitCode() int {
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

func ConfigError(msg string) *AbsError {
	return &AbsError{Code: ErrConfig, Message: msg}
}

func AuthError(msg string) *AbsError {
	return &AbsError{Code: ErrAuth, Message: msg}
}

func NotFoundError(msg string) *AbsError {
	return &AbsError{Code: ErrNotFound, Message: msg}
}

func NetworkError(msg string) *AbsError {
	return &AbsError{Code: ErrNetwork, Message: msg}
}

func APIError(msg string) *AbsError {
	return &AbsError{Code: ErrAPI, Message: fmt.Sprintf("API error: %s", msg)}
}
