use thiserror::Error;

#[derive(Error, Debug)]
pub enum PortainerError {
    #[error("Missing configuration: {0}")]
    ConfigError(String),

    #[error("Authentication failed: {0}")]
    AuthError(String),

    #[error("Resource not found: {0}")]
    NotFound(String),

    #[error("Network error: {0}")]
    NetworkError(String),

    #[error("API error: {0}")]
    ApiError(String),
}

impl PortainerError {
    pub fn exit_code(&self) -> i32 {
        match self {
            PortainerError::ConfigError(_) => 1,
            PortainerError::AuthError(_) => 2,
            PortainerError::NotFound(_) => 3,
            PortainerError::NetworkError(_) => 4,
            PortainerError::ApiError(_) => 5,
        }
    }

    pub fn code(&self) -> &'static str {
        match self {
            PortainerError::ConfigError(_) => "CONFIG_ERROR",
            PortainerError::AuthError(_) => "AUTH_FAILED",
            PortainerError::NotFound(_) => "NOT_FOUND",
            PortainerError::NetworkError(_) => "NETWORK_ERROR",
            PortainerError::ApiError(_) => "API_ERROR",
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_exit_codes() {
        assert_eq!(PortainerError::ConfigError("".into()).exit_code(), 1);
        assert_eq!(PortainerError::AuthError("".into()).exit_code(), 2);
        assert_eq!(PortainerError::NotFound("".into()).exit_code(), 3);
        assert_eq!(PortainerError::NetworkError("".into()).exit_code(), 4);
        assert_eq!(PortainerError::ApiError("".into()).exit_code(), 5);
    }

    #[test]
    fn test_error_codes() {
        assert_eq!(PortainerError::ConfigError("".into()).code(), "CONFIG_ERROR");
        assert_eq!(PortainerError::AuthError("".into()).code(), "AUTH_FAILED");
        assert_eq!(PortainerError::NotFound("".into()).code(), "NOT_FOUND");
        assert_eq!(PortainerError::NetworkError("".into()).code(), "NETWORK_ERROR");
        assert_eq!(PortainerError::ApiError("".into()).code(), "API_ERROR");
    }
}