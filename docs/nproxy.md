# nproxy-cli

CLI tool for nginx-proxy-manager API - backup and viewing operations.

## Installation

Copy your preferred binary to your PATH:
```bash
cp go/nproxy/nproxy-cli /usr/local/bin/nproxy-cli
# or
cp rust/nproxy/nproxy-cli /usr/local/bin/nproxy-cli
```

## Configuration

Set credentials via environment variables (recommended):
```bash
export NPROXY_URL=https://nginx.example.com
export NPROXY_TOKEN=eyJhbGciOiJS...
```

Or pass as flags:
```bash
nproxy-cli --url https://nginx.example.com --token eyJ... hosts list
```

### Options

| Flag | Short | Description |
|------|-------|-------------|
| `--url` | | nginx-proxy-manager URL |
| `--token` | | API token (JWT) |
| `--insecure` | `-k` | Skip TLS certificate verification |
| `--help` | `-h` | Show help |
| `--version` | `-v`/`-V` | Show version |

## Usage

### Login

Get a token by authenticating with email/password:
```bash
nproxy-cli login
# Email: admin@example.com
# Password: ****
# eyJhbGciOiJS...
```

Save the output to NPROXY_TOKEN.

### Proxy Hosts

```bash
# List all proxy hosts
nproxy-cli hosts list

# Show proxy host details
nproxy-cli hosts show 1
```

### Certificates

```bash
# List all certificates
nproxy-cli certificates list
nproxy-cli certs list  # alias

# Show certificate details
nproxy-cli certificates show 1
```

## Output Format

All output is YAML. Errors go to stderr with structured format:

```yaml
error:
  code: AUTH_FAILED
  message: Invalid or expired token
```

Exit codes: 1=config, 2=auth, 3=not found, 4=network, 5=api error

## Example Output

```yaml
# hosts list
hosts:
- id: 1
  domainNames:
  - example.com
  forwardHost: 10.0.0.1
  forwardPort: 8080
  sslForced: true
  enabled: true

# certificates list
certificates:
- id: 1
  niceName: '*.example.com'
  provider: letsencrypt
  expiresOn: 2026-02-28 17:40:58
```
