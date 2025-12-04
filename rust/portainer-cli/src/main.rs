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
