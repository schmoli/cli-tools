use serde::{Deserialize, Serialize};

// Helper function for status formatting
fn format_status(status: i32) -> String {
    match status {
        1 => "active".to_string(),
        2 => "inactive".to_string(),
        _ => format!("unknown({})", status),
    }
}

// Raw API response types
#[derive(Debug, Deserialize)]
pub struct ApiStack {
    #[serde(rename = "Id")]
    pub id: i64,
    #[serde(rename = "Name")]
    pub name: String,
    #[serde(rename = "Type")]
    pub stack_type: i32,
    #[serde(rename = "Status")]
    pub status: i32,
    #[serde(rename = "EndpointId")]
    pub endpoint_id: i64,
    #[serde(rename = "Env", default)]
    pub env: Vec<ApiEnvVar>,
}

#[derive(Debug, Deserialize, Serialize, Clone)]
pub struct ApiEnvVar {
    pub name: String,
    pub value: String,
}

#[derive(Debug, Deserialize)]
pub struct ApiStackFile {
    #[serde(rename = "StackFileContent")]
    pub stack_file_content: String,
}

#[derive(Debug, Deserialize)]
pub struct ApiEndpoint {
    #[serde(rename = "Id")]
    pub id: i64,
    #[serde(rename = "Name")]
    pub name: String,
    #[serde(rename = "Type")]
    pub endpoint_type: i32,
    #[serde(rename = "Status")]
    pub status: i32,
    #[serde(rename = "URL")]
    pub url: String,
}

// Output types (curated, with string labels)
#[derive(Debug, Serialize)]
pub struct Stack {
    pub id: i64,
    pub name: String,
    #[serde(rename = "type")]
    pub stack_type: String,
    pub status: String,
    #[serde(rename = "endpointId")]
    pub endpoint_id: i64,
    #[serde(skip_serializing_if = "Vec::is_empty")]
    pub env: Vec<ApiEnvVar>,
    #[serde(rename = "stackFile", skip_serializing_if = "Option::is_none")]
    pub stack_file: Option<String>,
}

#[derive(Debug, Serialize)]
pub struct StackList {
    pub stacks: Vec<StackListItem>,
}

#[derive(Debug, Serialize)]
pub struct StackListItem {
    pub id: i64,
    pub name: String,
    #[serde(rename = "type")]
    pub stack_type: String,
    pub status: String,
    #[serde(rename = "endpointId")]
    pub endpoint_id: i64,
}

#[derive(Debug, Serialize)]
pub struct Endpoint {
    pub id: i64,
    pub name: String,
    #[serde(rename = "type")]
    pub endpoint_type: String,
    pub status: String,
    pub url: String,
}

#[derive(Debug, Serialize)]
pub struct EndpointList {
    pub endpoints: Vec<Endpoint>,
}

// Mapping functions
impl ApiStack {
    pub fn stack_type_label(&self) -> String {
        match self.stack_type {
            1 => "swarm".to_string(),
            2 => "compose".to_string(),
            3 => "kubernetes".to_string(),
            _ => format!("unknown({})", self.stack_type),
        }
    }

    pub fn status_label(&self) -> String {
        format_status(self.status)
    }

    pub fn to_list_item(&self) -> StackListItem {
        StackListItem {
            id: self.id,
            name: self.name.clone(),
            stack_type: self.stack_type_label(),
            status: self.status_label(),
            endpoint_id: self.endpoint_id,
        }
    }

    pub fn to_stack(&self, stack_file: Option<String>) -> Stack {
        Stack {
            id: self.id,
            name: self.name.clone(),
            stack_type: self.stack_type_label(),
            status: self.status_label(),
            endpoint_id: self.endpoint_id,
            env: self.env.clone(),
            stack_file,
        }
    }
}

impl ApiEndpoint {
    pub fn endpoint_type_label(&self) -> String {
        match self.endpoint_type {
            1 => "docker".to_string(),
            2 => "agent".to_string(),
            3 => "azure".to_string(),
            4 => "edge-agent".to_string(),
            5 => "kubernetes".to_string(),
            _ => format!("unknown({})", self.endpoint_type),
        }
    }

    pub fn status_label(&self) -> String {
        format_status(self.status)
    }

    pub fn to_endpoint(&self) -> Endpoint {
        Endpoint {
            id: self.id,
            name: self.name.clone(),
            endpoint_type: self.endpoint_type_label(),
            status: self.status_label(),
            url: self.url.clone(),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_stack_type_label() {
        let cases = [(1, "swarm"), (2, "compose"), (3, "kubernetes")];
        for (type_code, expected) in cases {
            let stack = ApiStack {
                id: 1, name: "test".into(), stack_type: type_code,
                status: 1, endpoint_id: 1, env: vec![],
            };
            assert_eq!(stack.stack_type_label(), expected);
        }
        // Unknown type
        let stack = ApiStack {
            id: 1, name: "test".into(), stack_type: 99,
            status: 1, endpoint_id: 1, env: vec![],
        };
        assert_eq!(stack.stack_type_label(), "unknown(99)");
    }

    #[test]
    fn test_stack_status_label() {
        let cases = [(1, "active"), (2, "inactive")];
        for (status, expected) in cases {
            let stack = ApiStack {
                id: 1, name: "test".into(), stack_type: 1,
                status, endpoint_id: 1, env: vec![],
            };
            assert_eq!(stack.status_label(), expected);
        }
        let stack = ApiStack {
            id: 1, name: "test".into(), stack_type: 1,
            status: 99, endpoint_id: 1, env: vec![],
        };
        assert_eq!(stack.status_label(), "unknown(99)");
    }

    #[test]
    fn test_endpoint_type_label() {
        let cases = [
            (1, "docker"), (2, "agent"), (3, "azure"),
            (4, "edge-agent"), (5, "kubernetes"),
        ];
        for (type_code, expected) in cases {
            let ep = ApiEndpoint {
                id: 1, name: "test".into(), endpoint_type: type_code,
                status: 1, url: "http://test".into(),
            };
            assert_eq!(ep.endpoint_type_label(), expected);
        }
        let ep = ApiEndpoint {
            id: 1, name: "test".into(), endpoint_type: 99,
            status: 1, url: "http://test".into(),
        };
        assert_eq!(ep.endpoint_type_label(), "unknown(99)");
    }

    #[test]
    fn test_endpoint_status_label() {
        let cases = [(1, "active"), (2, "inactive")];
        for (status, expected) in cases {
            let ep = ApiEndpoint {
                id: 1, name: "test".into(), endpoint_type: 1,
                status, url: "http://test".into(),
            };
            assert_eq!(ep.status_label(), expected);
        }
    }

    #[test]
    fn test_to_list_item() {
        let stack = ApiStack {
            id: 42, name: "mystack".into(), stack_type: 2,
            status: 1, endpoint_id: 5, env: vec![],
        };
        let item = stack.to_list_item();
        assert_eq!(item.id, 42);
        assert_eq!(item.name, "mystack");
        assert_eq!(item.stack_type, "compose");
        assert_eq!(item.status, "active");
        assert_eq!(item.endpoint_id, 5);
    }

    #[test]
    fn test_to_endpoint() {
        let ep = ApiEndpoint {
            id: 10, name: "prod".into(), endpoint_type: 1,
            status: 2, url: "tcp://docker:2375".into(),
        };
        let result = ep.to_endpoint();
        assert_eq!(result.id, 10);
        assert_eq!(result.name, "prod");
        assert_eq!(result.endpoint_type, "docker");
        assert_eq!(result.status, "inactive");
        assert_eq!(result.url, "tcp://docker:2375");
    }
}
