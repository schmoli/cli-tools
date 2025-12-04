# portainer-cli Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build `portainer-cli` in both Rust and Go as parallel implementations for comparison.

**Architecture:** Two independent implementations (`rust/` and `go/`) sharing the same requirements. Each has a shared library for HTTP/auth/output and a CLI binary. Monorepo structure for future tools.

**Tech Stack:**
- Rust: clap (CLI), reqwest (HTTP), serde + serde_yaml (serialization)
- Go: cobra (CLI), net/http (HTTP), gopkg.in/yaml.v3 (serialization)

**API Endpoints:**
- `GET /api/stacks` - list stacks
- `GET /api/stacks/{id}?endpointId={eid}` - get stack (needs endpointId param)
- `GET /api/stacks/{id}/file` - get stack file content
- `GET /api/endpoints` - list endpoints
- `GET /api/endpoints/{id}` - get endpoint

**Auth:** `X-API-Key: <token>` header

---

## Part 1: Rust Implementation

### Task R1: Initialize Rust Workspace

**Files:**
- Create: `rust/Cargo.toml` (workspace)
- Create: `rust/portainer-lib/Cargo.toml`
- Create: `rust/portainer-lib/src/lib.rs`
- Create: `rust/portainer-cli/Cargo.toml`
- Create: `rust/portainer-cli/src/main.rs`

**Step 1: Create workspace Cargo.toml**

```toml
[workspace]
members = ["portainer-lib", "portainer-cli"]
resolver = "2"
```

**Step 2: Create library crate**

`rust/portainer-lib/Cargo.toml`:
```toml
[package]
name = "portainer-lib"
version = "0.1.0"
edition = "2021"

[dependencies]
reqwest = { version = "0.11", features = ["json", "blocking"] }
serde = { version = "1.0", features = ["derive"] }
serde_yaml = "0.9"
thiserror = "1.0"
```

`rust/portainer-lib/src/lib.rs`:
```rust
pub mod client;
pub mod models;
pub mod error;
pub mod output;
```

**Step 3: Create CLI crate**

`rust/portainer-cli/Cargo.toml`:
```toml
[package]
name = "portainer-cli"
version = "0.1.0"
edition = "2021"

[[bin]]
name = "portainer-cli"
path = "src/main.rs"

[dependencies]
portainer-lib = { path = "../portainer-lib" }
clap = { version = "4.4", features = ["derive"] }
```

`rust/portainer-cli/src/main.rs`:
```rust
fn main() {
    println!("portainer-cli");
}
```

**Step 4: Verify build**

Run: `cd rust && cargo build`
Expected: Compiles successfully

**Step 5: Commit**

```bash
git add rust/
git commit -m "feat(rust): init workspace with lib and cli crates"
```

---

### Task R2: Define Error Types

**Files:**
- Create: `rust/portainer-lib/src/error.rs`

**Step 1: Write error enum**

```rust
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
```

**Step 2: Verify build**

Run: `cd rust && cargo build`
Expected: Compiles

**Step 3: Commit**

```bash
git add rust/portainer-lib/src/error.rs
git commit -m "feat(rust): add error types with exit codes"
```

---

### Task R3: Define Data Models

**Files:**
- Create: `rust/portainer-lib/src/models.rs`

**Step 1: Write models**

```rust
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
```

**Step 2: Verify build**

Run: `cd rust && cargo build`
Expected: Compiles

**Step 3: Commit**

```bash
git add rust/portainer-lib/src/models.rs
git commit -m "feat(rust): add API and output models with type mappings"
```

---

### Task R4: Implement HTTP Client

**Files:**
- Create: `rust/portainer-lib/src/client.rs`

**Step 1: Write client**

```rust
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
    pub fn new(url: &str, token: &str) -> Result<Self, PortainerError> {
        let base_url = url.trim_end_matches('/').to_string();

        let client = Client::builder()
            .timeout(Duration::from_secs(10))
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
```

**Step 2: Verify build**

Run: `cd rust && cargo build`
Expected: Compiles

**Step 3: Commit**

```bash
git add rust/portainer-lib/src/client.rs
git commit -m "feat(rust): add HTTP client with auth and error handling"
```

---

### Task R5: Implement Output Formatting

**Files:**
- Create: `rust/portainer-lib/src/output.rs`

**Step 1: Write output module**

```rust
use crate::error::PortainerError;
use serde::Serialize;

#[derive(Serialize)]
pub struct ErrorOutput {
    pub error: ErrorDetail,
}

#[derive(Serialize)]
pub struct ErrorDetail {
    pub code: String,
    pub message: String,
}

pub fn print_yaml<T: Serialize>(data: &T) -> Result<(), PortainerError> {
    let yaml = serde_yaml::to_string(data)
        .map_err(|e| PortainerError::ApiError(format!("Failed to serialize: {}", e)))?;
    print!("{}", yaml);
    Ok(())
}

pub fn print_error(err: &PortainerError) {
    let output = ErrorOutput {
        error: ErrorDetail {
            code: err.code().to_string(),
            message: err.to_string(),
        },
    };
    if let Ok(yaml) = serde_yaml::to_string(&output) {
        eprint!("{}", yaml);
    } else {
        eprintln!("error: {}", err);
    }
}
```

**Step 2: Verify build**

Run: `cd rust && cargo build`
Expected: Compiles

**Step 3: Commit**

```bash
git add rust/portainer-lib/src/output.rs
git commit -m "feat(rust): add YAML output formatting"
```

---

### Task R6: Implement CLI with Clap

**Files:**
- Modify: `rust/portainer-cli/src/main.rs`

**Step 1: Write CLI**

```rust
use clap::{Parser, Subcommand};
use portainer_lib::client::PortainerClient;
use portainer_lib::error::PortainerError;
use portainer_lib::models::{EndpointList, StackList};
use portainer_lib::output::{print_error, print_yaml};
use std::process::ExitCode;

#[derive(Parser)]
#[command(name = "portainer-cli")]
#[command(version, about = "CLI for Portainer API")]
struct Cli {
    /// Portainer URL (or set PORTAINER_URL)
    #[arg(long, global = true)]
    url: Option<String>,

    /// API token (or set PORTAINER_TOKEN)
    #[arg(long, global = true)]
    token: Option<String>,

    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Manage stacks
    Stacks {
        #[command(subcommand)]
        action: StacksAction,
    },
    /// Manage endpoints
    Endpoints {
        #[command(subcommand)]
        action: EndpointsAction,
    },
}

#[derive(Subcommand)]
enum StacksAction {
    /// List all stacks
    List,
    /// Show a stack by ID
    Show {
        /// Stack ID
        id: i64,
    },
}

#[derive(Subcommand)]
enum EndpointsAction {
    /// List all endpoints
    List,
    /// Show an endpoint by ID
    Show {
        /// Endpoint ID
        id: i64,
    },
}

fn get_config(cli: &Cli) -> Result<(String, String), PortainerError> {
    let url = cli.url.clone()
        .or_else(|| std::env::var("PORTAINER_URL").ok())
        .ok_or_else(|| PortainerError::ConfigError(
            "Missing URL. Use --url or set PORTAINER_URL".to_string()
        ))?;

    let token = cli.token.clone()
        .or_else(|| std::env::var("PORTAINER_TOKEN").ok())
        .ok_or_else(|| PortainerError::ConfigError(
            "Missing token. Use --token or set PORTAINER_TOKEN".to_string()
        ))?;

    Ok((url, token))
}

fn run() -> Result<(), PortainerError> {
    let cli = Cli::parse();
    let (url, token) = get_config(&cli)?;
    let client = PortainerClient::new(&url, &token)?;

    match cli.command {
        Commands::Stacks { action } => match action {
            StacksAction::List => {
                let stacks = client.list_stacks()?;
                let output = StackList {
                    stacks: stacks.iter().map(|s| s.to_list_item()).collect(),
                };
                print_yaml(&output)?;
            }
            StacksAction::Show { id } => {
                // First get stack list to find endpoint_id
                let stacks = client.list_stacks()?;
                let api_stack = stacks.iter()
                    .find(|s| s.id == id)
                    .ok_or_else(|| PortainerError::NotFound(format!("Stack with ID {}", id)))?;

                // Get stack file content
                let file = client.get_stack_file(id)?;
                let stack = api_stack.to_stack(Some(file.stack_file_content));
                print_yaml(&stack)?;
            }
        },
        Commands::Endpoints { action } => match action {
            EndpointsAction::List => {
                let endpoints = client.list_endpoints()?;
                let output = EndpointList {
                    endpoints: endpoints.iter().map(|e| e.to_endpoint()).collect(),
                };
                print_yaml(&output)?;
            }
            EndpointsAction::Show { id } => {
                let endpoint = client.get_endpoint(id)?;
                print_yaml(&endpoint.to_endpoint())?;
            }
        },
    }

    Ok(())
}

fn main() -> ExitCode {
    match run() {
        Ok(()) => ExitCode::SUCCESS,
        Err(e) => {
            print_error(&e);
            ExitCode::from(e.exit_code() as u8)
        }
    }
}
```

**Step 2: Verify build**

Run: `cd rust && cargo build`
Expected: Compiles

**Step 3: Test help output**

Run: `cd rust && cargo run -- --help`
Expected: Shows CLI help with stacks and endpoints commands

**Step 4: Test missing config error**

Run: `cd rust && cargo run -- stacks list 2>&1`
Expected: YAML error output about missing URL

**Step 5: Commit**

```bash
git add rust/portainer-cli/src/main.rs
git commit -m "feat(rust): implement CLI with stacks and endpoints commands"
```

---

### Task R7: Add .gitignore for Rust

**Files:**
- Create: `rust/.gitignore`

**Step 1: Write gitignore**

```
/target/
Cargo.lock
```

**Step 2: Commit**

```bash
git add rust/.gitignore
git commit -m "chore(rust): add gitignore"
```

---

## Part 2: Go Implementation

### Task G1: Initialize Go Module

**Files:**
- Create: `go/go.mod`
- Create: `go/cmd/portainer-cli/main.go`
- Create: `go/pkg/portainer/client.go`

**Step 1: Create go.mod**

```
module github.com/toli/portainer-cli

go 1.21

require (
	github.com/spf13/cobra v1.8.0
	gopkg.in/yaml.v3 v3.0.1
)
```

**Step 2: Create placeholder files**

`go/cmd/portainer-cli/main.go`:
```go
package main

import "fmt"

func main() {
	fmt.Println("portainer-cli")
}
```

`go/pkg/portainer/client.go`:
```go
package portainer

// Client placeholder
type Client struct{}
```

**Step 3: Initialize dependencies**

Run: `cd go && go mod tidy`
Expected: Downloads dependencies

**Step 4: Verify build**

Run: `cd go && go build ./cmd/portainer-cli`
Expected: Compiles

**Step 5: Commit**

```bash
git add go/
git commit -m "feat(go): init module with cmd and pkg structure"
```

---

### Task G2: Define Error Types

**Files:**
- Create: `go/pkg/portainer/errors.go`

**Step 1: Write errors**

```go
package portainer

import "fmt"

type ErrorCode string

const (
	ErrConfig  ErrorCode = "CONFIG_ERROR"
	ErrAuth    ErrorCode = "AUTH_FAILED"
	ErrNotFound ErrorCode = "NOT_FOUND"
	ErrNetwork ErrorCode = "NETWORK_ERROR"
	ErrAPI     ErrorCode = "API_ERROR"
)

type PortainerError struct {
	Code    ErrorCode
	Message string
}

func (e *PortainerError) Error() string {
	return e.Message
}

func (e *PortainerError) ExitCode() int {
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

func ConfigError(msg string) *PortainerError {
	return &PortainerError{Code: ErrConfig, Message: msg}
}

func AuthError(msg string) *PortainerError {
	return &PortainerError{Code: ErrAuth, Message: msg}
}

func NotFoundError(msg string) *PortainerError {
	return &PortainerError{Code: ErrNotFound, Message: msg}
}

func NetworkError(msg string) *PortainerError {
	return &PortainerError{Code: ErrNetwork, Message: msg}
}

func APIError(msg string) *PortainerError {
	return &PortainerError{Code: ErrAPI, Message: fmt.Sprintf("API error: %s", msg)}
}
```

**Step 2: Verify build**

Run: `cd go && go build ./...`
Expected: Compiles

**Step 3: Commit**

```bash
git add go/pkg/portainer/errors.go
git commit -m "feat(go): add error types with exit codes"
```

---

### Task G3: Define Data Models

**Files:**
- Create: `go/pkg/portainer/models.go`

**Step 1: Write models**

```go
package portainer

// API response types (match Portainer JSON)
type APIStack struct {
	ID         int64       `json:"Id"`
	Name       string      `json:"Name"`
	Type       int         `json:"Type"`
	Status     int         `json:"Status"`
	EndpointID int64       `json:"EndpointId"`
	Env        []APIEnvVar `json:"Env"`
}

type APIEnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type APIStackFile struct {
	StackFileContent string `json:"StackFileContent"`
}

type APIEndpoint struct {
	ID     int64  `json:"Id"`
	Name   string `json:"Name"`
	Type   int    `json:"Type"`
	Status int    `json:"Status"`
	URL    string `json:"URL"`
}

// Output types (curated, YAML output)
type Stack struct {
	ID         int64       `yaml:"id"`
	Name       string      `yaml:"name"`
	Type       string      `yaml:"type"`
	Status     string      `yaml:"status"`
	EndpointID int64       `yaml:"endpointId"`
	Env        []APIEnvVar `yaml:"env,omitempty"`
	StackFile  string      `yaml:"stackFile,omitempty"`
}

type StackList struct {
	Stacks []StackListItem `yaml:"stacks"`
}

type StackListItem struct {
	ID         int64  `yaml:"id"`
	Name       string `yaml:"name"`
	Type       string `yaml:"type"`
	Status     string `yaml:"status"`
	EndpointID int64  `yaml:"endpointId"`
}

type Endpoint struct {
	ID     int64  `yaml:"id"`
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Status string `yaml:"status"`
	URL    string `yaml:"url"`
}

type EndpointList struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Mapping functions
func (s *APIStack) TypeLabel() string {
	switch s.Type {
	case 1:
		return "swarm"
	case 2:
		return "compose"
	case 3:
		return "kubernetes"
	default:
		return "unknown"
	}
}

func (s *APIStack) StatusLabel() string {
	switch s.Status {
	case 1:
		return "active"
	case 2:
		return "inactive"
	default:
		return "unknown"
	}
}

func (s *APIStack) ToListItem() StackListItem {
	return StackListItem{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.TypeLabel(),
		Status:     s.StatusLabel(),
		EndpointID: s.EndpointID,
	}
}

func (s *APIStack) ToStack(stackFile string) Stack {
	return Stack{
		ID:         s.ID,
		Name:       s.Name,
		Type:       s.TypeLabel(),
		Status:     s.StatusLabel(),
		EndpointID: s.EndpointID,
		Env:        s.Env,
		StackFile:  stackFile,
	}
}

func (e *APIEndpoint) TypeLabel() string {
	switch e.Type {
	case 1:
		return "docker"
	case 2:
		return "agent"
	case 3:
		return "azure"
	case 4:
		return "edge-agent"
	case 5:
		return "kubernetes"
	default:
		return "unknown"
	}
}

func (e *APIEndpoint) StatusLabel() string {
	switch e.Status {
	case 1:
		return "up"
	case 2:
		return "down"
	default:
		return "unknown"
	}
}

func (e *APIEndpoint) ToEndpoint() Endpoint {
	return Endpoint{
		ID:     e.ID,
		Name:   e.Name,
		Type:   e.TypeLabel(),
		Status: e.StatusLabel(),
		URL:    e.URL,
	}
}
```

**Step 2: Verify build**

Run: `cd go && go build ./...`
Expected: Compiles

**Step 3: Commit**

```bash
git add go/pkg/portainer/models.go
git commit -m "feat(go): add API and output models with type mappings"
```

---

### Task G4: Implement HTTP Client

**Files:**
- Modify: `go/pkg/portainer/client.go`

**Step 1: Write client**

```go
package portainer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(url, token string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(url, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) get(path string, result interface{}) error {
	url := c.baseURL + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NetworkError(err.Error())
	}

	req.Header.Set("X-API-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NetworkError(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return NetworkError(err.Error())
		}
		if err := json.Unmarshal(body, result); err != nil {
			return APIError(fmt.Sprintf("failed to parse response: %s", err))
		}
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return AuthError("invalid or expired token")
	case http.StatusNotFound:
		return NotFoundError(fmt.Sprintf("resource not found: %s", path))
	default:
		return APIError(fmt.Sprintf("unexpected status: %d", resp.StatusCode))
	}
}

func (c *Client) ListStacks() ([]APIStack, error) {
	var stacks []APIStack
	if err := c.get("/api/stacks", &stacks); err != nil {
		return nil, err
	}
	return stacks, nil
}

func (c *Client) GetStack(id, endpointID int64) (*APIStack, error) {
	var stack APIStack
	path := fmt.Sprintf("/api/stacks/%d?endpointId=%d", id, endpointID)
	if err := c.get(path, &stack); err != nil {
		return nil, err
	}
	return &stack, nil
}

func (c *Client) GetStackFile(id int64) (*APIStackFile, error) {
	var file APIStackFile
	path := fmt.Sprintf("/api/stacks/%d/file", id)
	if err := c.get(path, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (c *Client) ListEndpoints() ([]APIEndpoint, error) {
	var endpoints []APIEndpoint
	if err := c.get("/api/endpoints", &endpoints); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (c *Client) GetEndpoint(id int64) (*APIEndpoint, error) {
	var endpoint APIEndpoint
	path := fmt.Sprintf("/api/endpoints/%d", id)
	if err := c.get(path, &endpoint); err != nil {
		return nil, err
	}
	return &endpoint, nil
}
```

**Step 2: Verify build**

Run: `cd go && go build ./...`
Expected: Compiles

**Step 3: Commit**

```bash
git add go/pkg/portainer/client.go
git commit -m "feat(go): add HTTP client with auth and error handling"
```

---

### Task G5: Implement Output Formatting

**Files:**
- Create: `go/pkg/portainer/output.go`

**Step 1: Write output module**

```go
package portainer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ErrorOutput struct {
	Error ErrorDetail `yaml:"error"`
}

type ErrorDetail struct {
	Code    string `yaml:"code"`
	Message string `yaml:"message"`
}

func PrintYAML(data interface{}) error {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return APIError(fmt.Sprintf("failed to serialize: %s", err))
	}
	return nil
}

func PrintError(err error) {
	pe, ok := err.(*PortainerError)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}

	output := ErrorOutput{
		Error: ErrorDetail{
			Code:    string(pe.Code),
			Message: pe.Message,
		},
	}

	enc := yaml.NewEncoder(os.Stderr)
	enc.SetIndent(2)
	enc.Encode(output)
}
```

**Step 2: Verify build**

Run: `cd go && go build ./...`
Expected: Compiles

**Step 3: Commit**

```bash
git add go/pkg/portainer/output.go
git commit -m "feat(go): add YAML output formatting"
```

---

### Task G6: Implement CLI with Cobra

**Files:**
- Modify: `go/cmd/portainer-cli/main.go`
- Create: `go/cmd/portainer-cli/stacks.go`
- Create: `go/cmd/portainer-cli/endpoints.go`

**Step 1: Write main.go**

```go
package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/toli/portainer-cli/pkg/portainer"
)

var (
	flagURL   string
	flagToken string
)

var rootCmd = &cobra.Command{
	Use:     "portainer-cli",
	Short:   "CLI for Portainer API",
	Version: "0.1.0",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Portainer URL (or set PORTAINER_URL)")
	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "API token (or set PORTAINER_TOKEN)")

	rootCmd.AddCommand(stacksCmd)
	rootCmd.AddCommand(endpointsCmd)
}

func getConfig() (string, string, error) {
	url := flagURL
	if url == "" {
		url = os.Getenv("PORTAINER_URL")
	}
	if url == "" {
		return "", "", portainer.ConfigError("missing URL. Use --url or set PORTAINER_URL")
	}

	token := flagToken
	if token == "" {
		token = os.Getenv("PORTAINER_TOKEN")
	}
	if token == "" {
		return "", "", portainer.ConfigError("missing token. Use --token or set PORTAINER_TOKEN")
	}

	return url, token, nil
}

func getClient() (*portainer.Client, error) {
	url, token, err := getConfig()
	if err != nil {
		return nil, err
	}
	return portainer.NewClient(url, token), nil
}

func handleError(err error) {
	portainer.PrintError(err)
	if pe, ok := err.(*portainer.PortainerError); ok {
		os.Exit(pe.ExitCode())
	}
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

**Step 2: Write stacks.go**

```go
package main

import (
	"github.com/spf13/cobra"
	"github.com/toli/portainer-cli/pkg/portainer"
)

var stacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "Manage stacks",
}

var stacksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stacks",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		stacks, err := client.ListStacks()
		if err != nil {
			handleError(err)
			return
		}

		output := portainer.StackList{
			Stacks: make([]portainer.StackListItem, len(stacks)),
		}
		for i, s := range stacks {
			output.Stacks[i] = s.ToListItem()
		}

		if err := portainer.PrintYAML(output); err != nil {
			handleError(err)
		}
	},
}

var stacksShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a stack by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		var id int64
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			handleError(portainer.ConfigError("invalid stack ID"))
			return
		}

		// Get stack list to find endpoint_id
		stacks, err := client.ListStacks()
		if err != nil {
			handleError(err)
			return
		}

		var apiStack *portainer.APIStack
		for _, s := range stacks {
			if s.ID == id {
				apiStack = &s
				break
			}
		}
		if apiStack == nil {
			handleError(portainer.NotFoundError(fmt.Sprintf("stack with ID %d", id)))
			return
		}

		// Get stack file content
		file, err := client.GetStackFile(id)
		if err != nil {
			handleError(err)
			return
		}

		stack := apiStack.ToStack(file.StackFileContent)
		if err := portainer.PrintYAML(stack); err != nil {
			handleError(err)
		}
	},
}

func init() {
	stacksCmd.AddCommand(stacksListCmd)
	stacksCmd.AddCommand(stacksShowCmd)
}
```

**Step 3: Add missing import to stacks.go**

Add `"fmt"` to imports in stacks.go.

**Step 4: Write endpoints.go**

```go
package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/toli/portainer-cli/pkg/portainer"
)

var endpointsCmd = &cobra.Command{
	Use:   "endpoints",
	Short: "Manage endpoints",
}

var endpointsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		endpoints, err := client.ListEndpoints()
		if err != nil {
			handleError(err)
			return
		}

		output := portainer.EndpointList{
			Endpoints: make([]portainer.Endpoint, len(endpoints)),
		}
		for i, e := range endpoints {
			output.Endpoints[i] = e.ToEndpoint()
		}

		if err := portainer.PrintYAML(output); err != nil {
			handleError(err)
		}
	},
}

var endpointsShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show an endpoint by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		var id int64
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			handleError(portainer.ConfigError("invalid endpoint ID"))
			return
		}

		endpoint, err := client.GetEndpoint(id)
		if err != nil {
			handleError(err)
			return
		}

		if err := portainer.PrintYAML(endpoint.ToEndpoint()); err != nil {
			handleError(err)
		}
	},
}

func init() {
	endpointsCmd.AddCommand(endpointsListCmd)
	endpointsCmd.AddCommand(endpointsShowCmd)
}
```

**Step 5: Run go mod tidy**

Run: `cd go && go mod tidy`
Expected: Dependencies resolved

**Step 6: Verify build**

Run: `cd go && go build ./cmd/portainer-cli`
Expected: Compiles

**Step 7: Test help output**

Run: `cd go && ./portainer-cli --help`
Expected: Shows CLI help with stacks and endpoints commands

**Step 8: Test missing config error**

Run: `cd go && ./portainer-cli stacks list 2>&1`
Expected: YAML error output about missing URL

**Step 9: Commit**

```bash
git add go/cmd/portainer-cli/
git commit -m "feat(go): implement CLI with stacks and endpoints commands"
```

---

### Task G7: Add .gitignore for Go

**Files:**
- Create: `go/.gitignore`

**Step 1: Write gitignore**

```
portainer-cli
*.exe
```

**Step 2: Commit**

```bash
git add go/.gitignore
git commit -m "chore(go): add gitignore"
```

---

## Part 3: Integration Verification

### Task V1: Test Both Implementations

**Step 1: Build both**

Run: `cd rust && cargo build --release`
Run: `cd go && go build -o portainer-cli ./cmd/portainer-cli`

**Step 2: Compare help output**

Run: `rust/target/release/portainer-cli --help`
Run: `go/portainer-cli --help`
Expected: Both show same command structure

**Step 3: Compare error output**

Run: `rust/target/release/portainer-cli stacks list 2>&1`
Run: `go/portainer-cli stacks list 2>&1`
Expected: Both show YAML error about missing URL

**Step 4: Document any differences**

If outputs differ, note in README which to prefer and why.

---

### Task V2: Add Root README

**Files:**
- Create: `README.md`

**Step 1: Write README**

```markdown
# portainer-cli

CLI tool for Portainer API - backup and viewing operations.

## Implementations

Two parallel implementations for comparison:

- `rust/` - Rust implementation using clap + reqwest
- `go/` - Go implementation using cobra + net/http

## Usage

```bash
# Set credentials
export PORTAINER_URL=https://portainer.example.com
export PORTAINER_TOKEN=ptr_xxxxxxxxxxxx

# List stacks
portainer-cli stacks list

# Show stack details
portainer-cli stacks show 1

# List endpoints
portainer-cli endpoints list

# Show endpoint details
portainer-cli endpoints show 1
```

## Build

### Rust
```bash
cd rust && cargo build --release
# Binary: rust/target/release/portainer-cli
```

### Go
```bash
cd go && go build -o portainer-cli ./cmd/portainer-cli
# Binary: go/portainer-cli
```

## Output Format

All output is YAML. Errors go to stderr.

See `docs/plans/2025-12-03-portainer-cli-requirements.md` for full specification.
```

**Step 2: Commit**

```bash
git add README.md
git commit -m "docs: add root README"
```

---

## Summary

| Task | Description |
|------|-------------|
| R1 | Init Rust workspace |
| R2 | Rust error types |
| R3 | Rust data models |
| R4 | Rust HTTP client |
| R5 | Rust output formatting |
| R6 | Rust CLI with clap |
| R7 | Rust gitignore |
| G1 | Init Go module |
| G2 | Go error types |
| G3 | Go data models |
| G4 | Go HTTP client |
| G5 | Go output formatting |
| G6 | Go CLI with cobra |
| G7 | Go gitignore |
| V1 | Verify both implementations |
| V2 | Root README |

Total: 16 tasks, ~2-5 min each.
