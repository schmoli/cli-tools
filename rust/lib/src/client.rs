use crate::error::PortainerError;
use crate::models::{ApiEndpoint, ApiStack, ApiStackFile};
use reqwest::blocking::Client;
use reqwest::StatusCode;
use std::time::Duration;

pub struct PortainerClient {
    base_url: String,
    token: String,
    client: Client,
}

impl PortainerClient {
    pub fn new(url: &str, token: &str, insecure: bool) -> Result<Self, PortainerError> {
        let base_url = url.trim_end_matches('/').to_string();

        let mut builder = Client::builder()
            .timeout(Duration::from_secs(10));

        if insecure {
            builder = builder.danger_accept_invalid_certs(true);
        }

        let client = builder
            .build()
            .map_err(|e| PortainerError::NetworkError(e.to_string()))?;

        Ok(Self {
            base_url,
            token: token.to_string(),
            client,
        })
    }

    fn get<T: serde::de::DeserializeOwned>(&self, path: &str) -> Result<T, PortainerError> {
        let url = format!("{}{}", self.base_url, path);

        let response = self.client
            .get(&url)
            .header("X-API-Key", &self.token)
            .send()
            .map_err(|e| PortainerError::NetworkError(e.to_string()))?;

        match response.status() {
            StatusCode::OK => {
                response.json::<T>()
                    .map_err(|e| PortainerError::ApiError(format!("Failed to parse response: {}", e)))
            }
            StatusCode::UNAUTHORIZED | StatusCode::FORBIDDEN => {
                Err(PortainerError::AuthError("Invalid or expired token".to_string()))
            }
            StatusCode::NOT_FOUND => {
                Err(PortainerError::NotFound(format!("Resource not found: {}", path)))
            }
            status => {
                Err(PortainerError::ApiError(format!("Unexpected status: {}", status)))
            }
        }
    }

    pub fn list_stacks(&self) -> Result<Vec<ApiStack>, PortainerError> {
        self.get("/api/stacks")
    }

    pub fn get_stack(&self, id: i64, endpoint_id: i64) -> Result<ApiStack, PortainerError> {
        self.get(&format!("/api/stacks/{}?endpointId={}", id, endpoint_id))
    }

    pub fn get_stack_file(&self, id: i64) -> Result<ApiStackFile, PortainerError> {
        self.get(&format!("/api/stacks/{}/file", id))
    }

    pub fn list_endpoints(&self) -> Result<Vec<ApiEndpoint>, PortainerError> {
        self.get("/api/endpoints")
    }

    pub fn get_endpoint(&self, id: i64) -> Result<ApiEndpoint, PortainerError> {
        self.get(&format!("/api/endpoints/{}", id))
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_url_trailing_slash_trimmed() {
        let client = PortainerClient::new("http://example.com/", "token", false).unwrap();
        assert_eq!(client.base_url, "http://example.com");
    }

    #[test]
    fn test_client_creation_succeeds() {
        let result = PortainerClient::new("http://example.com", "token", false);
        assert!(result.is_ok());
    }

    #[test]
    fn test_client_creation_with_insecure() {
        let result = PortainerClient::new("http://example.com", "token", true);
        assert!(result.is_ok());
    }
}
