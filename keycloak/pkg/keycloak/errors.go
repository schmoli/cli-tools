package keycloak

type ErrorCode string

const (
	ErrConfig   ErrorCode = "CONFIG_ERROR"
	ErrAuth     ErrorCode = "AUTH_FAILED"
	ErrNotFound ErrorCode = "NOT_FOUND"
	ErrNetwork  ErrorCode = "NETWORK_ERROR"
	ErrAPI      ErrorCode = "API_ERROR"
)

type KeycloakError struct {
	Code    ErrorCode
	Message string
}

func (e *KeycloakError) Error() string {
	return e.Message
}

func (e *KeycloakError) ExitCode() int {
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

func ConfigError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrConfig, Message: msg}
}

func AuthError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrAuth, Message: msg}
}

func NotFoundError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrNotFound, Message: msg}
}

func NetworkError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrNetwork, Message: msg}
}

func APIError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrAPI, Message: msg}
}
