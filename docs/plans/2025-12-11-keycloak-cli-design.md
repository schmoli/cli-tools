# keycloak-cli Design

Read-only CLI for Keycloak Admin REST API.

## Scope (Phase 1)

| Resource | Commands |
|----------|----------|
| Realms | `list`, `get` |
| Users | `list`, `get`, `sessions` |
| Clients | `list`, `get`, `sessions` |
| Roles | `list`, `get` (realm + client level) |
| Groups | `list`, `get`, `members` |

## Authentication

Client credentials flow via gocloak.

```
KEYCLOAK_URL=https://keycloak.example.com
KEYCLOAK_REALM=master
KEYCLOAK_CLIENT_ID=keycloak-cli
KEYCLOAK_CLIENT_SECRET=xxx
```

Flags: `--url`, `--realm`, `--client-id`, `--client-secret`, `--insecure/-k`

**Keycloak setup:**
1. Create client `keycloak-cli` in `master` realm
2. Enable "Client authentication" (confidential)
3. Enable "Service accounts roles"
4. Assign realm-management admin roles to service account

## Command Structure

```bash
keycloak-cli realms list
keycloak-cli realms get <realm-name>

keycloak-cli users list --realm <realm>
keycloak-cli users get <user-id> --realm <realm>
keycloak-cli users sessions <user-id> --realm <realm>

keycloak-cli clients list --realm <realm>
keycloak-cli clients get <client-id> --realm <realm>
keycloak-cli clients sessions <client-id> --realm <realm>

keycloak-cli roles list --realm <realm>
keycloak-cli roles get <role-name> --realm <realm>
keycloak-cli roles list --realm <realm> --client <client-id>
keycloak-cli roles get <role-name> --realm <realm> --client <client-id>

keycloak-cli groups list --realm <realm>
keycloak-cli groups get <group-id> --realm <realm>
keycloak-cli groups members <group-id> --realm <realm>
```

## File Structure

```
keycloak/
├── cmd/keycloak-cli/
│   ├── main.go      # cobra root, global flags
│   ├── realms.go
│   ├── users.go
│   ├── clients.go
│   ├── roles.go
│   └── groups.go
├── pkg/keycloak/
│   ├── client.go    # wraps gocloak, auth
│   ├── errors.go    # typed errors + exit codes
│   └── output.go    # YAML printing
├── go.mod
└── go.sum
```

## Dependencies

- `github.com/Nerzal/gocloak/v13` — Keycloak API client
- `github.com/spf13/cobra` — CLI framework
- `gopkg.in/yaml.v3` — YAML output

## Output

YAML only (matching other cli-tools).

## API Reference

Base: `{url}/admin/realms/{realm}/...`

| Resource | Endpoint |
|----------|----------|
| Realms | `GET /admin/realms` |
| Users | `GET /admin/realms/{realm}/users` |
| User sessions | `GET /admin/realms/{realm}/users/{id}/sessions` |
| Clients | `GET /admin/realms/{realm}/clients` |
| Client sessions | `GET /admin/realms/{realm}/clients/{id}/user-sessions` |
| Roles (realm) | `GET /admin/realms/{realm}/roles` |
| Roles (client) | `GET /admin/realms/{realm}/clients/{id}/roles` |
| Groups | `GET /admin/realms/{realm}/groups` |
| Group members | `GET /admin/realms/{realm}/groups/{id}/members` |
