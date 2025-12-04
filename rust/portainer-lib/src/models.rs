use serde::{Deserialize, Serialize};

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
        match self.status {
            1 => "active".to_string(),
            2 => "inactive".to_string(),
            _ => format!("unknown({})", self.status),
        }
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
        match self.status {
            1 => "up".to_string(),
            2 => "down".to_string(),
            _ => format!("unknown({})", self.status),
        }
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
