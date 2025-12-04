use crate::error::NproxyError;
use crate::models::{ApiCertificate, ApiLoginResponse, ApiProxyHost};
use reqwest::blocking::Client;
use reqwest::StatusCode;
use std::time::Duration;

pub struct NproxyClient {
    base_url: String,
    token: String,
    client: Client,
}

impl NproxyClient {
    pub fn new(url: &str, token: &str, insecure: bool) -> Result<Self, NproxyError> {
        let base_url = url.trim_end_matches('/').to_string();

        if !base_url.starts_with("http://") && !base_url.starts_with("https://") {
            return Err(NproxyError::ConfigError("URL must start with http:// or https://".to_string()));
        }

        let mut builder = Client::builder()
            .timeout(Duration::from_secs(10));

        if insecure {
            builder = builder.danger_accept_invalid_certs(true);
        }

        let client = builder
            .build()
            .map_err(|e| NproxyError::NetworkError(e.to_string()))?;

        Ok(Self {
            base_url,
            token: token.to_string(),
            client,
        })
    }

    fn get<T: serde::de::DeserializeOwned>(&self, path: &str) -> Result<T, NproxyError> {
        let url = format!("{}{}", self.base_url, path);

        let response = self.client
            .get(&url)
            .header("Authorization", format!("Bearer {}", self.token))
            .send()
            .map_err(|e| NproxyError::NetworkError(e.to_string()))?;

        match response.status() {
            StatusCode::OK => {
                response.json::<T>()
                    .map_err(|e| NproxyError::ApiError(format!("Failed to parse response from {}: {}", path, e)))
            }
            StatusCode::UNAUTHORIZED | StatusCode::FORBIDDEN => {
                Err(NproxyError::AuthError("Invalid or expired token".to_string()))
            }
            StatusCode::NOT_FOUND => {
                Err(NproxyError::NotFound(format!("Resource not found: {}", path)))
            }
            status => {
                Err(NproxyError::ApiError(format!("Unexpected status {} from {}", status, path)))
            }
        }
    }

    pub fn list_proxy_hosts(&self) -> Result<Vec<ApiProxyHost>, NproxyError> {
        self.get("/api/nginx/proxy-hosts")
    }

    pub fn get_proxy_host(&self, id: i64) -> Result<ApiProxyHost, NproxyError> {
        self.get(&format!("/api/nginx/proxy-hosts/{}", id))
    }

    pub fn list_certificates(&self) -> Result<Vec<ApiCertificate>, NproxyError> {
        self.get("/api/nginx/certificates")
    }

    pub fn get_certificate(&self, id: i64) -> Result<ApiCertificate, NproxyError> {
        self.get(&format!("/api/nginx/certificates/{}", id))
    }
}

pub fn login(url: &str, email: &str, password: &str, insecure: bool) -> Result<String, NproxyError> {
    let base_url = url.trim_end_matches('/');

    if !base_url.starts_with("http://") && !base_url.starts_with("https://") {
        return Err(NproxyError::ConfigError("URL must start with http:// or https://".to_string()));
    }

    let mut builder = Client::builder()
        .timeout(Duration::from_secs(10));

    if insecure {
        builder = builder.danger_accept_invalid_certs(true);
    }

    let client = builder
        .build()
        .map_err(|e| NproxyError::NetworkError(e.to_string()))?;

    let payload = serde_json::json!({
        "identity": email,
        "secret": password
    });

    let response = client
        .post(format!("{}/api/tokens", base_url))
        .json(&payload)
        .send()
        .map_err(|e| NproxyError::NetworkError(e.to_string()))?;

    match response.status() {
        StatusCode::OK => {
            let result: ApiLoginResponse = response.json()
                .map_err(|e| NproxyError::ApiError(format!("Failed to parse login response: {}", e)))?;
            Ok(result.token)
        }
        StatusCode::UNAUTHORIZED | StatusCode::FORBIDDEN => {
            Err(NproxyError::AuthError("Invalid credentials".to_string()))
        }
        status => {
            Err(NproxyError::ApiError(format!("Login failed with status {}", status)))
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_url_trailing_slash_trimmed() {
        let client = NproxyClient::new("http://example.com/", "token", false).unwrap();
        assert_eq!(client.base_url, "http://example.com");
    }

    #[test]
    fn test_client_creation_succeeds() {
        let result = NproxyClient::new("http://example.com", "token", false);
        assert!(result.is_ok());
    }

    #[test]
    fn test_client_creation_with_insecure() {
        let result = NproxyClient::new("http://example.com", "token", true);
        assert!(result.is_ok());
    }

    #[test]
    fn test_invalid_url_rejected() {
        let result = NproxyClient::new("not-a-url", "token", false);
        assert!(result.is_err());
    }
}
