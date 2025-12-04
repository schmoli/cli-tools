use serde::{Deserialize, Serialize};

// Raw API response types
#[derive(Debug, Deserialize)]
pub struct ApiProxyHost {
    pub id: i64,
    pub domain_names: Vec<String>,
    pub forward_scheme: String,
    pub forward_host: String,
    pub forward_port: i32,
    pub certificate_id: Option<i64>,
    pub ssl_forced: bool,
    pub block_exploits: bool,
    pub caching_enabled: bool,
    #[serde(default)]
    pub allow_websocket_upgrade: bool,
    pub enabled: bool,
    #[serde(default)]
    pub advanced_config: String,
}

#[derive(Debug, Deserialize)]
pub struct ApiCertificate {
    pub id: i64,
    pub provider: String,
    pub nice_name: String,
    #[serde(default)]
    pub domain_names: Vec<String>,
    pub expires_on: String,
}

#[derive(Debug, Deserialize)]
pub struct ApiLoginResponse {
    pub token: String,
}

// Output types (curated, with camelCase for YAML)
#[derive(Debug, Serialize)]
pub struct ProxyHost {
    pub id: i64,
    #[serde(rename = "domainNames")]
    pub domain_names: Vec<String>,
    #[serde(rename = "forwardScheme")]
    pub forward_scheme: String,
    #[serde(rename = "forwardHost")]
    pub forward_host: String,
    #[serde(rename = "forwardPort")]
    pub forward_port: i32,
    #[serde(rename = "certificateId", skip_serializing_if = "Option::is_none")]
    pub certificate_id: Option<i64>,
    #[serde(rename = "sslForced")]
    pub ssl_forced: bool,
    #[serde(rename = "blockExploits")]
    pub block_exploits: bool,
    #[serde(rename = "cachingEnabled")]
    pub caching_enabled: bool,
    pub websocket: bool,
    pub enabled: bool,
    #[serde(rename = "advancedConfig", skip_serializing_if = "String::is_empty")]
    pub advanced_config: String,
}

#[derive(Debug, Serialize)]
pub struct ProxyHostList {
    pub hosts: Vec<ProxyHostListItem>,
}

#[derive(Debug, Serialize)]
pub struct ProxyHostListItem {
    pub id: i64,
    #[serde(rename = "domainNames")]
    pub domain_names: Vec<String>,
    #[serde(rename = "forwardHost")]
    pub forward_host: String,
    #[serde(rename = "forwardPort")]
    pub forward_port: i32,
    #[serde(rename = "sslForced")]
    pub ssl_forced: bool,
    pub enabled: bool,
}

#[derive(Debug, Serialize)]
pub struct Certificate {
    pub id: i64,
    pub provider: String,
    #[serde(rename = "niceName")]
    pub nice_name: String,
    #[serde(rename = "domainNames")]
    pub domain_names: Vec<String>,
    #[serde(rename = "expiresOn")]
    pub expires_on: String,
}

#[derive(Debug, Serialize)]
pub struct CertificateList {
    pub certificates: Vec<CertificateListItem>,
}

#[derive(Debug, Serialize)]
pub struct CertificateListItem {
    pub id: i64,
    #[serde(rename = "niceName")]
    pub nice_name: String,
    pub provider: String,
    #[serde(rename = "expiresOn")]
    pub expires_on: String,
}

// Mapping functions
impl ApiProxyHost {
    pub fn to_list_item(&self) -> ProxyHostListItem {
        ProxyHostListItem {
            id: self.id,
            domain_names: self.domain_names.clone(),
            forward_host: self.forward_host.clone(),
            forward_port: self.forward_port,
            ssl_forced: self.ssl_forced,
            enabled: self.enabled,
        }
    }

    pub fn to_proxy_host(&self) -> ProxyHost {
        ProxyHost {
            id: self.id,
            domain_names: self.domain_names.clone(),
            forward_scheme: self.forward_scheme.clone(),
            forward_host: self.forward_host.clone(),
            forward_port: self.forward_port,
            certificate_id: self.certificate_id,
            ssl_forced: self.ssl_forced,
            block_exploits: self.block_exploits,
            caching_enabled: self.caching_enabled,
            websocket: self.allow_websocket_upgrade,
            enabled: self.enabled,
            advanced_config: self.advanced_config.clone(),
        }
    }
}

impl ApiCertificate {
    pub fn to_list_item(&self) -> CertificateListItem {
        CertificateListItem {
            id: self.id,
            nice_name: self.nice_name.clone(),
            provider: self.provider.clone(),
            expires_on: self.expires_on.clone(),
        }
    }

    pub fn to_certificate(&self) -> Certificate {
        Certificate {
            id: self.id,
            provider: self.provider.clone(),
            nice_name: self.nice_name.clone(),
            domain_names: self.domain_names.clone(),
            expires_on: self.expires_on.clone(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_proxy_host_to_list_item() {
        let host = ApiProxyHost {
            id: 1,
            domain_names: vec!["example.com".into()],
            forward_scheme: "https".into(),
            forward_host: "backend".into(),
            forward_port: 8080,
            certificate_id: Some(5),
            ssl_forced: true,
            block_exploits: true,
            caching_enabled: false,
            allow_websocket_upgrade: true,
            enabled: true,
            advanced_config: String::new(),
        };

        let item = host.to_list_item();
        assert_eq!(item.id, 1);
        assert_eq!(item.domain_names, vec!["example.com"]);
        assert_eq!(item.forward_host, "backend");
        assert_eq!(item.forward_port, 8080);
        assert!(item.ssl_forced);
        assert!(item.enabled);
    }

    #[test]
    fn test_certificate_to_list_item() {
        let cert = ApiCertificate {
            id: 1,
            provider: "letsencrypt".into(),
            nice_name: "My Cert".into(),
            domain_names: vec!["example.com".into()],
            expires_on: "2024-12-31".into(),
        };

        let item = cert.to_list_item();
        assert_eq!(item.id, 1);
        assert_eq!(item.nice_name, "My Cert");
        assert_eq!(item.provider, "letsencrypt");
        assert_eq!(item.expires_on, "2024-12-31");
    }
}
