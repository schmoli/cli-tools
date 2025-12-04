use clap::{Parser, Subcommand};
use nproxy_lib::client::{login, NproxyClient};
use nproxy_lib::error::NproxyError;
use nproxy_lib::models::{CertificateList, ProxyHostList};
use nproxy_lib::output::{print_error, print_yaml};
use std::io::{self, Write};
use std::process::ExitCode;

#[derive(Parser)]
#[command(name = "nproxy-cli")]
#[command(version, about = "CLI for nginx-proxy-manager API")]
struct Cli {
    /// nginx-proxy-manager URL (or set NPROXY_URL)
    #[arg(long, global = true)]
    url: Option<String>,

    /// API token (or set NPROXY_TOKEN)
    #[arg(long, global = true)]
    token: Option<String>,

    /// Skip TLS certificate verification
    #[arg(short = 'k', long, global = true)]
    insecure: bool,

    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Authenticate and get a token
    Login,
    /// Manage proxy hosts
    Hosts {
        #[command(subcommand)]
        action: HostsAction,
    },
    /// Manage certificates
    #[command(alias = "certs")]
    Certificates {
        #[command(subcommand)]
        action: CertificatesAction,
    },
}

#[derive(Subcommand)]
enum HostsAction {
    /// List all proxy hosts
    List,
    /// Show a proxy host by ID
    Show {
        /// Proxy host ID
        id: i64,
    },
}

#[derive(Subcommand)]
enum CertificatesAction {
    /// List all certificates
    List,
    /// Show a certificate by ID
    Show {
        /// Certificate ID
        id: i64,
    },
}

fn get_url(cli: &Cli) -> Result<String, NproxyError> {
    cli.url.clone()
        .or_else(|| std::env::var("NPROXY_URL").ok())
        .ok_or_else(|| NproxyError::ConfigError(
            "Missing URL. Use --url or set NPROXY_URL".to_string()
        ))
}

fn get_config(cli: &Cli) -> Result<(String, String), NproxyError> {
    let url = get_url(cli)?;

    let token = cli.token.clone()
        .or_else(|| std::env::var("NPROXY_TOKEN").ok())
        .ok_or_else(|| NproxyError::ConfigError(
            "Missing token. Use --token or set NPROXY_TOKEN".to_string()
        ))?;

    Ok((url, token))
}

fn read_line(prompt: &str) -> Result<String, NproxyError> {
    print!("{}", prompt);
    io::stdout().flush().ok();
    let mut input = String::new();
    io::stdin().read_line(&mut input)
        .map_err(|_| NproxyError::ConfigError("Failed to read input".to_string()))?;
    Ok(input.trim().to_string())
}

fn read_password(prompt: &str) -> Result<String, NproxyError> {
    print!("{}", prompt);
    io::stdout().flush().ok();
    let password = rpassword::read_password()
        .map_err(|_| NproxyError::ConfigError("Failed to read password".to_string()))?;
    Ok(password)
}

fn run() -> Result<(), NproxyError> {
    let cli = Cli::parse();

    match cli.command {
        Commands::Login => {
            let url = get_url(&cli)?;
            let email = read_line("Email: ")?;
            let password = read_password("Password: ")?;

            let token = login(&url, &email, &password, cli.insecure)?;
            println!("{}", token);
        }
        Commands::Hosts { ref action } => {
            let (url, token) = get_config(&cli)?;
            let client = NproxyClient::new(&url, &token, cli.insecure)?;

            match action {
                HostsAction::List => {
                    let hosts = client.list_proxy_hosts()?;
                    let output = ProxyHostList {
                        hosts: hosts.iter().map(|h| h.to_list_item()).collect(),
                    };
                    print_yaml(&output)?;
                }
                HostsAction::Show { id } => {
                    if *id <= 0 {
                        return Err(NproxyError::ConfigError("ID must be positive".to_string()));
                    }
                    let host = client.get_proxy_host(*id)?;
                    print_yaml(&host.to_proxy_host())?;
                }
            }
        }
        Commands::Certificates { ref action } => {
            let (url, token) = get_config(&cli)?;
            let client = NproxyClient::new(&url, &token, cli.insecure)?;

            match action {
                CertificatesAction::List => {
                    let certs = client.list_certificates()?;
                    let output = CertificateList {
                        certificates: certs.iter().map(|c| c.to_list_item()).collect(),
                    };
                    print_yaml(&output)?;
                }
                CertificatesAction::Show { id } => {
                    if *id <= 0 {
                        return Err(NproxyError::ConfigError("ID must be positive".to_string()));
                    }
                    let cert = client.get_certificate(*id)?;
                    print_yaml(&cert.to_certificate())?;
                }
            }
        }
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
