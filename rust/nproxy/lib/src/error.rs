use thiserror::Error;

#[derive(Error, Debug)]
pub enum NproxyError {
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

impl NproxyError {
    pub fn exit_code(&self) -> i32 {
        match self {
            NproxyError::ConfigError(_) => 1,
            NproxyError::AuthError(_) => 2,
            NproxyError::NotFound(_) => 3,
            NproxyError::NetworkError(_) => 4,
            NproxyError::ApiError(_) => 5,
        }
    }

    pub fn code(&self) -> &'static str {
        match self {
            NproxyError::ConfigError(_) => "CONFIG_ERROR",
            NproxyError::AuthError(_) => "AUTH_FAILED",
            NproxyError::NotFound(_) => "NOT_FOUND",
            NproxyError::NetworkError(_) => "NETWORK_ERROR",
            NproxyError::ApiError(_) => "API_ERROR",
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_exit_codes() {
        assert_eq!(NproxyError::ConfigError("".into()).exit_code(), 1);
        assert_eq!(NproxyError::AuthError("".into()).exit_code(), 2);
        assert_eq!(NproxyError::NotFound("".into()).exit_code(), 3);
        assert_eq!(NproxyError::NetworkError("".into()).exit_code(), 4);
        assert_eq!(NproxyError::ApiError("".into()).exit_code(), 5);
    }

    #[test]
    fn test_error_codes() {
        assert_eq!(NproxyError::ConfigError("".into()).code(), "CONFIG_ERROR");
        assert_eq!(NproxyError::AuthError("".into()).code(), "AUTH_FAILED");
        assert_eq!(NproxyError::NotFound("".into()).code(), "NOT_FOUND");
        assert_eq!(NproxyError::NetworkError("".into()).code(), "NETWORK_ERROR");
        assert_eq!(NproxyError::ApiError("".into()).code(), "API_ERROR");
    }
}
