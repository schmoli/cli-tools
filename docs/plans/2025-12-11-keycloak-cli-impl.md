# keycloak-cli Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Read-only CLI for Keycloak Admin API using gocloak library.

**Architecture:** Cobra CLI wrapping gocloak client. Client credentials auth. YAML output matching portainer-cli pattern.

**Tech Stack:** Go, gocloak/v13, cobra, yaml.v3

---

## Task 1: Initialize Go Module

**Files:**
- Create: `keycloak/go.mod`

**Step 1: Create directory and init module**

```bash
cd /Users/toli/code/claude/cli-tools/.worktrees/keycloak-cli
mkdir -p keycloak/cmd/keycloak-cli keycloak/pkg/keycloak
cd keycloak && go mod init github.com/schmoli/cli-tools/keycloak
```

**Step 2: Add to go.work**

Edit `go.work` to add `./keycloak` to the use block.

**Step 3: Add dependencies**

```bash
cd keycloak && go get github.com/Nerzal/gocloak/v13 github.com/spf13/cobra gopkg.in/yaml.v3
```

**Step 4: Commit**

```bash
git add keycloak/go.mod keycloak/go.sum go.work go.work.sum
git commit -m "feat(keycloak): init module with deps"
```

---

## Task 2: Create Error Types

**Files:**
- Create: `keycloak/pkg/keycloak/errors.go`
- Create: `keycloak/pkg/keycloak/errors_test.go`

**Step 1: Write failing test**

```go
// keycloak/pkg/keycloak/errors_test.go
package keycloak

import "testing"

func TestConfigError(t *testing.T) {
	err := ConfigError("test message")
	if err.Code != ErrConfig {
		t.Errorf("expected ErrConfig, got %s", err.Code)
	}
	if err.Error() != "test message" {
		t.Errorf("expected 'test message', got %s", err.Error())
	}
	if err.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", err.ExitCode())
	}
}

func TestAuthError(t *testing.T) {
	err := AuthError("auth failed")
	if err.Code != ErrAuth {
		t.Errorf("expected ErrAuth, got %s", err.Code)
	}
	if err.ExitCode() != 2 {
		t.Errorf("expected exit code 2, got %d", err.ExitCode())
	}
}
```

**Step 2: Run test to verify failure**

```bash
cd keycloak && go test ./pkg/keycloak/... -v
```

Expected: FAIL

**Step 3: Implement errors.go**

```go
// keycloak/pkg/keycloak/errors.go
package keycloak

type ErrorCode string

const (
	ErrConfig   ErrorCode = "CONFIG_ERROR"
	ErrAuth     ErrorCode = "AUTH_FAILED"
	ErrNotFound ErrorCode = "NOT_FOUND"
	ErrNetwork  ErrorCode = "NETWORK_ERROR"
	ErrAPI      ErrorCode = "API_ERROR"
)

type KeycloakError struct {
	Code    ErrorCode
	Message string
}

func (e *KeycloakError) Error() string {
	return e.Message
}

func (e *KeycloakError) ExitCode() int {
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

func ConfigError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrConfig, Message: msg}
}

func AuthError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrAuth, Message: msg}
}

func NotFoundError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrNotFound, Message: msg}
}

func NetworkError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrNetwork, Message: msg}
}

func APIError(msg string) *KeycloakError {
	return &KeycloakError{Code: ErrAPI, Message: msg}
}
```

**Step 4: Run test to verify pass**

```bash
cd keycloak && go test ./pkg/keycloak/... -v
```

Expected: PASS

**Step 5: Commit**

```bash
git add keycloak/pkg/keycloak/errors.go keycloak/pkg/keycloak/errors_test.go
git commit -m "feat(keycloak): add error types"
```

---

## Task 3: Create Output Helpers

**Files:**
- Create: `keycloak/pkg/keycloak/output.go`

**Step 1: Implement output.go**

```go
// keycloak/pkg/keycloak/output.go
package keycloak

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
	ke, ok := err.(*KeycloakError)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return
	}

	output := ErrorOutput{
		Error: ErrorDetail{
			Code:    string(ke.Code),
			Message: ke.Message,
		},
	}

	enc := yaml.NewEncoder(os.Stderr)
	enc.SetIndent(2)
	if err := enc.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", ke.Message)
	}
}
```

**Step 2: Verify compiles**

```bash
cd keycloak && go build ./pkg/keycloak/...
```

**Step 3: Commit**

```bash
git add keycloak/pkg/keycloak/output.go
git commit -m "feat(keycloak): add YAML output helpers"
```

---

## Task 4: Create Client Wrapper

**Files:**
- Create: `keycloak/pkg/keycloak/client.go`

**Step 1: Implement client.go**

```go
// keycloak/pkg/keycloak/client.go
package keycloak

import (
	"context"
	"crypto/tls"

	"github.com/Nerzal/gocloak/v13"
)

type Client struct {
	gocloak gocloak.GoCloak
	token   string
	ctx     context.Context
}

type Config struct {
	URL          string
	Realm        string
	ClientID     string
	ClientSecret string
	Insecure     bool
}

func NewClient(cfg Config) (*Client, error) {
	if cfg.URL == "" {
		return nil, ConfigError("missing URL")
	}
	if cfg.Realm == "" {
		return nil, ConfigError("missing realm")
	}
	if cfg.ClientID == "" {
		return nil, ConfigError("missing client ID")
	}
	if cfg.ClientSecret == "" {
		return nil, ConfigError("missing client secret")
	}

	gc := gocloak.NewClient(cfg.URL)
	if cfg.Insecure {
		restyClient := gc.RestyClient()
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	ctx := context.Background()
	token, err := gc.LoginClient(ctx, cfg.ClientID, cfg.ClientSecret, cfg.Realm)
	if err != nil {
		return nil, AuthError(err.Error())
	}

	return &Client{
		gocloak: gc,
		token:   token.AccessToken,
		ctx:     ctx,
	}, nil
}
```

**Step 2: Verify compiles**

```bash
cd keycloak && go build ./pkg/keycloak/...
```

**Step 3: Commit**

```bash
git add keycloak/pkg/keycloak/client.go
git commit -m "feat(keycloak): add gocloak client wrapper"
```

---

## Task 5: Add Realm Methods

**Files:**
- Modify: `keycloak/pkg/keycloak/client.go`

**Step 1: Add realm methods to client.go**

```go
// Add to client.go

type RealmInfo struct {
	ID          string `yaml:"id"`
	Realm       string `yaml:"realm"`
	DisplayName string `yaml:"display_name,omitempty"`
	Enabled     bool   `yaml:"enabled"`
}

type RealmList struct {
	Realms []RealmInfo `yaml:"realms"`
}

func (c *Client) ListRealms() (*RealmList, error) {
	realms, err := c.gocloak.GetRealms(c.ctx, c.token)
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &RealmList{Realms: make([]RealmInfo, len(realms))}
	for i, r := range realms {
		list.Realms[i] = RealmInfo{
			ID:          deref(r.ID),
			Realm:       deref(r.Realm),
			DisplayName: deref(r.DisplayName),
			Enabled:     derefBool(r.Enabled),
		}
	}
	return list, nil
}

func (c *Client) GetRealm(name string) (*RealmInfo, error) {
	r, err := c.gocloak.GetRealm(c.ctx, c.token, name)
	if err != nil {
		return nil, NotFoundError(err.Error())
	}
	return &RealmInfo{
		ID:          deref(r.ID),
		Realm:       deref(r.Realm),
		DisplayName: deref(r.DisplayName),
		Enabled:     derefBool(r.Enabled),
	}, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}
```

**Step 2: Verify compiles**

```bash
cd keycloak && go build ./pkg/keycloak/...
```

**Step 3: Commit**

```bash
git add keycloak/pkg/keycloak/client.go
git commit -m "feat(keycloak): add realm list/get methods"
```

---

## Task 6: Create Main CLI Entry Point

**Files:**
- Create: `keycloak/cmd/keycloak-cli/main.go`

**Step 1: Implement main.go**

```go
// keycloak/cmd/keycloak-cli/main.go
package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var version = "dev"

var (
	flagURL          string
	flagRealm        string
	flagClientID     string
	flagClientSecret string
	flagInsecure     bool
	flagTargetRealm  string
)

var rootCmd = &cobra.Command{
	Use:     "keycloak-cli",
	Short:   "CLI for Keycloak Admin API",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "Keycloak URL (or KEYCLOAK_URL)")
	rootCmd.PersistentFlags().StringVar(&flagRealm, "realm", "", "Auth realm (or KEYCLOAK_REALM)")
	rootCmd.PersistentFlags().StringVar(&flagClientID, "client-id", "", "Client ID (or KEYCLOAK_CLIENT_ID)")
	rootCmd.PersistentFlags().StringVar(&flagClientSecret, "client-secret", "", "Client secret (or KEYCLOAK_CLIENT_SECRET)")
	rootCmd.PersistentFlags().BoolVarP(&flagInsecure, "insecure", "k", false, "Skip TLS verification")
	rootCmd.PersistentFlags().StringVar(&flagTargetRealm, "target-realm", "", "Target realm for queries (or KEYCLOAK_TARGET_REALM)")

	rootCmd.AddCommand(realmsCmd)
}

func envOrFlag(flag, env string) string {
	if flag != "" {
		return flag
	}
	return os.Getenv(env)
}

func getClient() (*keycloak.Client, error) {
	cfg := keycloak.Config{
		URL:          envOrFlag(flagURL, "KEYCLOAK_URL"),
		Realm:        envOrFlag(flagRealm, "KEYCLOAK_REALM"),
		ClientID:     envOrFlag(flagClientID, "KEYCLOAK_CLIENT_ID"),
		ClientSecret: envOrFlag(flagClientSecret, "KEYCLOAK_CLIENT_SECRET"),
		Insecure:     flagInsecure,
	}
	return keycloak.NewClient(cfg)
}

func getTargetRealm() string {
	return envOrFlag(flagTargetRealm, "KEYCLOAK_TARGET_REALM")
}

func handleError(err error) {
	keycloak.PrintError(err)
	if ke, ok := err.(*keycloak.KeycloakError); ok {
		os.Exit(ke.ExitCode())
	}
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

**Step 2: Verify compiles (will fail, needs realmsCmd)**

Create stub realms.go first.

**Step 3: Commit with realms.go in next task**

---

## Task 7: Add Realms Commands

**Files:**
- Create: `keycloak/cmd/keycloak-cli/realms.go`

**Step 1: Implement realms.go**

```go
// keycloak/cmd/keycloak-cli/realms.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var realmsCmd = &cobra.Command{
	Use:   "realms",
	Short: "Manage realms",
}

var realmsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all realms",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		realms, err := client.ListRealms()
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(realms); err != nil {
			handleError(err)
		}
	},
}

var realmsGetCmd = &cobra.Command{
	Use:   "get <realm-name>",
	Short: "Get realm details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		realm, err := client.GetRealm(args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(realm); err != nil {
			handleError(err)
		}
	},
}

func init() {
	realmsCmd.AddCommand(realmsListCmd)
	realmsCmd.AddCommand(realmsGetCmd)
}
```

**Step 2: Verify compiles**

```bash
cd keycloak && go build ./cmd/keycloak-cli/...
```

**Step 3: Commit**

```bash
git add keycloak/cmd/keycloak-cli/main.go keycloak/cmd/keycloak-cli/realms.go
git commit -m "feat(keycloak): add realms list/get commands"
```

---

## Task 8: Add Users Commands

**Files:**
- Modify: `keycloak/pkg/keycloak/client.go` (add user methods)
- Create: `keycloak/cmd/keycloak-cli/users.go`

**Step 1: Add user methods to client.go**

```go
// Add to client.go

type UserInfo struct {
	ID        string `yaml:"id"`
	Username  string `yaml:"username"`
	Email     string `yaml:"email,omitempty"`
	FirstName string `yaml:"first_name,omitempty"`
	LastName  string `yaml:"last_name,omitempty"`
	Enabled   bool   `yaml:"enabled"`
}

type UserList struct {
	Users []UserInfo `yaml:"users"`
}

type SessionInfo struct {
	ID         string `yaml:"id"`
	Username   string `yaml:"username"`
	ClientID   string `yaml:"client_id,omitempty"`
	IPAddress  string `yaml:"ip_address,omitempty"`
	Started    int64  `yaml:"started,omitempty"`
	LastAccess int64  `yaml:"last_access,omitempty"`
}

type SessionList struct {
	Sessions []SessionInfo `yaml:"sessions"`
}

func (c *Client) ListUsers(realm string) (*UserList, error) {
	users, err := c.gocloak.GetUsers(c.ctx, c.token, realm, gocloak.GetUsersParams{})
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &UserList{Users: make([]UserInfo, len(users))}
	for i, u := range users {
		list.Users[i] = UserInfo{
			ID:        deref(u.ID),
			Username:  deref(u.Username),
			Email:     deref(u.Email),
			FirstName: deref(u.FirstName),
			LastName:  deref(u.LastName),
			Enabled:   derefBool(u.Enabled),
		}
	}
	return list, nil
}

func (c *Client) GetUser(realm, userID string) (*UserInfo, error) {
	u, err := c.gocloak.GetUserByID(c.ctx, c.token, realm, userID)
	if err != nil {
		return nil, NotFoundError(err.Error())
	}
	return &UserInfo{
		ID:        deref(u.ID),
		Username:  deref(u.Username),
		Email:     deref(u.Email),
		FirstName: deref(u.FirstName),
		LastName:  deref(u.LastName),
		Enabled:   derefBool(u.Enabled),
	}, nil
}

func (c *Client) GetUserSessions(realm, userID string) (*SessionList, error) {
	sessions, err := c.gocloak.GetUserSessions(c.ctx, c.token, realm, userID)
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &SessionList{Sessions: make([]SessionInfo, len(sessions))}
	for i, s := range sessions {
		list.Sessions[i] = SessionInfo{
			ID:         deref(s.ID),
			Username:   deref(s.Username),
			ClientID:   deref(s.ClientID),
			IPAddress:  deref(s.IPAddress),
			Started:    derefInt64(s.Start),
			LastAccess: derefInt64(s.LastAccess),
		}
	}
	return list, nil
}

func derefInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
```

**Step 2: Implement users.go**

```go
// keycloak/cmd/keycloak-cli/users.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users in realm",
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		users, err := client.ListUsers(realm)
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(users); err != nil {
			handleError(err)
		}
	},
}

var usersGetCmd = &cobra.Command{
	Use:   "get <user-id>",
	Short: "Get user details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		user, err := client.GetUser(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(user); err != nil {
			handleError(err)
		}
	},
}

var usersSessionsCmd = &cobra.Command{
	Use:   "sessions <user-id>",
	Short: "List user sessions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		sessions, err := client.GetUserSessions(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(sessions); err != nil {
			handleError(err)
		}
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)
	usersCmd.AddCommand(usersGetCmd)
	usersCmd.AddCommand(usersSessionsCmd)
}
```

**Step 3: Add usersCmd to main.go init()**

```go
rootCmd.AddCommand(usersCmd)
```

**Step 4: Verify compiles**

```bash
cd keycloak && go build ./cmd/keycloak-cli/...
```

**Step 5: Commit**

```bash
git add keycloak/pkg/keycloak/client.go keycloak/cmd/keycloak-cli/users.go keycloak/cmd/keycloak-cli/main.go
git commit -m "feat(keycloak): add users list/get/sessions commands"
```

---

## Task 9: Add Clients Commands

**Files:**
- Modify: `keycloak/pkg/keycloak/client.go`
- Create: `keycloak/cmd/keycloak-cli/clients.go`

**Step 1: Add client methods to client.go**

```go
// Add to client.go

type ClientInfo struct {
	ID          string `yaml:"id"`
	ClientID    string `yaml:"client_id"`
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
	Enabled     bool   `yaml:"enabled"`
	Protocol    string `yaml:"protocol,omitempty"`
}

type ClientList struct {
	Clients []ClientInfo `yaml:"clients"`
}

func (c *Client) ListClients(realm string) (*ClientList, error) {
	clients, err := c.gocloak.GetClients(c.ctx, c.token, realm, gocloak.GetClientsParams{})
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &ClientList{Clients: make([]ClientInfo, len(clients))}
	for i, cl := range clients {
		list.Clients[i] = ClientInfo{
			ID:          deref(cl.ID),
			ClientID:    deref(cl.ClientID),
			Name:        deref(cl.Name),
			Description: deref(cl.Description),
			Enabled:     derefBool(cl.Enabled),
			Protocol:    deref(cl.Protocol),
		}
	}
	return list, nil
}

func (c *Client) GetClient(realm, clientID string) (*ClientInfo, error) {
	clients, err := c.gocloak.GetClients(c.ctx, c.token, realm, gocloak.GetClientsParams{ClientID: &clientID})
	if err != nil {
		return nil, APIError(err.Error())
	}
	if len(clients) == 0 {
		return nil, NotFoundError("client not found: " + clientID)
	}
	cl := clients[0]
	return &ClientInfo{
		ID:          deref(cl.ID),
		ClientID:    deref(cl.ClientID),
		Name:        deref(cl.Name),
		Description: deref(cl.Description),
		Enabled:     derefBool(cl.Enabled),
		Protocol:    deref(cl.Protocol),
	}, nil
}

func (c *Client) GetClientSessions(realm, clientUUID string) (*SessionList, error) {
	sessions, err := c.gocloak.GetClientUserSessions(c.ctx, c.token, realm, clientUUID, gocloak.GetClientUserSessionsParams{})
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &SessionList{Sessions: make([]SessionInfo, len(sessions))}
	for i, s := range sessions {
		list.Sessions[i] = SessionInfo{
			ID:         deref(s.ID),
			Username:   deref(s.Username),
			ClientID:   deref(s.ClientID),
			IPAddress:  deref(s.IPAddress),
			Started:    derefInt64(s.Start),
			LastAccess: derefInt64(s.LastAccess),
		}
	}
	return list, nil
}
```

**Step 2: Implement clients.go**

```go
// keycloak/cmd/keycloak-cli/clients.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Manage clients",
}

var clientsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List clients in realm",
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		clients, err := client.ListClients(realm)
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(clients); err != nil {
			handleError(err)
		}
	},
}

var clientsGetCmd = &cobra.Command{
	Use:   "get <client-id>",
	Short: "Get client details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		cl, err := client.GetClient(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(cl); err != nil {
			handleError(err)
		}
	},
}

var clientsSessionsCmd = &cobra.Command{
	Use:   "sessions <client-uuid>",
	Short: "List client sessions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		sessions, err := client.GetClientSessions(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(sessions); err != nil {
			handleError(err)
		}
	},
}

func init() {
	clientsCmd.AddCommand(clientsListCmd)
	clientsCmd.AddCommand(clientsGetCmd)
	clientsCmd.AddCommand(clientsSessionsCmd)
}
```

**Step 3: Add clientsCmd to main.go init()**

```go
rootCmd.AddCommand(clientsCmd)
```

**Step 4: Verify compiles**

```bash
cd keycloak && go build ./cmd/keycloak-cli/...
```

**Step 5: Commit**

```bash
git add keycloak/pkg/keycloak/client.go keycloak/cmd/keycloak-cli/clients.go keycloak/cmd/keycloak-cli/main.go
git commit -m "feat(keycloak): add clients list/get/sessions commands"
```

---

## Task 10: Add Roles Commands

**Files:**
- Modify: `keycloak/pkg/keycloak/client.go`
- Create: `keycloak/cmd/keycloak-cli/roles.go`

**Step 1: Add role methods to client.go**

```go
// Add to client.go

type RoleInfo struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Composite   bool   `yaml:"composite"`
	ClientRole  bool   `yaml:"client_role"`
}

type RoleList struct {
	Roles []RoleInfo `yaml:"roles"`
}

func (c *Client) ListRealmRoles(realm string) (*RoleList, error) {
	roles, err := c.gocloak.GetRealmRoles(c.ctx, c.token, realm, gocloak.GetRoleParams{})
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &RoleList{Roles: make([]RoleInfo, len(roles))}
	for i, r := range roles {
		list.Roles[i] = RoleInfo{
			ID:          deref(r.ID),
			Name:        deref(r.Name),
			Description: deref(r.Description),
			Composite:   derefBool(r.Composite),
			ClientRole:  derefBool(r.ClientRole),
		}
	}
	return list, nil
}

func (c *Client) GetRealmRole(realm, roleName string) (*RoleInfo, error) {
	r, err := c.gocloak.GetRealmRole(c.ctx, c.token, realm, roleName)
	if err != nil {
		return nil, NotFoundError(err.Error())
	}
	return &RoleInfo{
		ID:          deref(r.ID),
		Name:        deref(r.Name),
		Description: deref(r.Description),
		Composite:   derefBool(r.Composite),
		ClientRole:  derefBool(r.ClientRole),
	}, nil
}

func (c *Client) ListClientRoles(realm, clientUUID string) (*RoleList, error) {
	roles, err := c.gocloak.GetClientRoles(c.ctx, c.token, realm, clientUUID, gocloak.GetRoleParams{})
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &RoleList{Roles: make([]RoleInfo, len(roles))}
	for i, r := range roles {
		list.Roles[i] = RoleInfo{
			ID:          deref(r.ID),
			Name:        deref(r.Name),
			Description: deref(r.Description),
			Composite:   derefBool(r.Composite),
			ClientRole:  derefBool(r.ClientRole),
		}
	}
	return list, nil
}

func (c *Client) GetClientRole(realm, clientUUID, roleName string) (*RoleInfo, error) {
	r, err := c.gocloak.GetClientRole(c.ctx, c.token, realm, clientUUID, roleName)
	if err != nil {
		return nil, NotFoundError(err.Error())
	}
	return &RoleInfo{
		ID:          deref(r.ID),
		Name:        deref(r.Name),
		Description: deref(r.Description),
		Composite:   derefBool(r.Composite),
		ClientRole:  derefBool(r.ClientRole),
	}, nil
}
```

**Step 2: Implement roles.go**

```go
// keycloak/cmd/keycloak-cli/roles.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var (
	flagClientUUID string
)

var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Manage roles",
}

var rolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List roles (realm or client)",
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		var roles *keycloak.RoleList
		if flagClientUUID != "" {
			roles, err = client.ListClientRoles(realm, flagClientUUID)
		} else {
			roles, err = client.ListRealmRoles(realm)
		}
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(roles); err != nil {
			handleError(err)
		}
	},
}

var rolesGetCmd = &cobra.Command{
	Use:   "get <role-name>",
	Short: "Get role details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		var role *keycloak.RoleInfo
		if flagClientUUID != "" {
			role, err = client.GetClientRole(realm, flagClientUUID, args[0])
		} else {
			role, err = client.GetRealmRole(realm, args[0])
		}
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(role); err != nil {
			handleError(err)
		}
	},
}

func init() {
	rolesCmd.PersistentFlags().StringVar(&flagClientUUID, "client", "", "Client UUID for client roles")
	rolesCmd.AddCommand(rolesListCmd)
	rolesCmd.AddCommand(rolesGetCmd)
}
```

**Step 3: Add rolesCmd to main.go init()**

```go
rootCmd.AddCommand(rolesCmd)
```

**Step 4: Verify compiles**

```bash
cd keycloak && go build ./cmd/keycloak-cli/...
```

**Step 5: Commit**

```bash
git add keycloak/pkg/keycloak/client.go keycloak/cmd/keycloak-cli/roles.go keycloak/cmd/keycloak-cli/main.go
git commit -m "feat(keycloak): add roles list/get commands"
```

---

## Task 11: Add Groups Commands

**Files:**
- Modify: `keycloak/pkg/keycloak/client.go`
- Create: `keycloak/cmd/keycloak-cli/groups.go`

**Step 1: Add group methods to client.go**

```go
// Add to client.go

type GroupInfo struct {
	ID        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Path      string   `yaml:"path"`
	SubGroups []string `yaml:"subgroups,omitempty"`
}

type GroupList struct {
	Groups []GroupInfo `yaml:"groups"`
}

type MemberList struct {
	Members []UserInfo `yaml:"members"`
}

func (c *Client) ListGroups(realm string) (*GroupList, error) {
	groups, err := c.gocloak.GetGroups(c.ctx, c.token, realm, gocloak.GetGroupsParams{})
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &GroupList{Groups: make([]GroupInfo, len(groups))}
	for i, g := range groups {
		var subs []string
		if g.SubGroups != nil {
			for _, sg := range *g.SubGroups {
				subs = append(subs, deref(sg.Name))
			}
		}
		list.Groups[i] = GroupInfo{
			ID:        deref(g.ID),
			Name:      deref(g.Name),
			Path:      deref(g.Path),
			SubGroups: subs,
		}
	}
	return list, nil
}

func (c *Client) GetGroup(realm, groupID string) (*GroupInfo, error) {
	g, err := c.gocloak.GetGroup(c.ctx, c.token, realm, groupID)
	if err != nil {
		return nil, NotFoundError(err.Error())
	}
	var subs []string
	if g.SubGroups != nil {
		for _, sg := range *g.SubGroups {
			subs = append(subs, deref(sg.Name))
		}
	}
	return &GroupInfo{
		ID:        deref(g.ID),
		Name:      deref(g.Name),
		Path:      deref(g.Path),
		SubGroups: subs,
	}, nil
}

func (c *Client) GetGroupMembers(realm, groupID string) (*MemberList, error) {
	members, err := c.gocloak.GetGroupMembers(c.ctx, c.token, realm, groupID, gocloak.GetGroupsParams{})
	if err != nil {
		return nil, APIError(err.Error())
	}

	list := &MemberList{Members: make([]UserInfo, len(members))}
	for i, u := range members {
		list.Members[i] = UserInfo{
			ID:        deref(u.ID),
			Username:  deref(u.Username),
			Email:     deref(u.Email),
			FirstName: deref(u.FirstName),
			LastName:  deref(u.LastName),
			Enabled:   derefBool(u.Enabled),
		}
	}
	return list, nil
}
```

**Step 2: Implement groups.go**

```go
// keycloak/cmd/keycloak-cli/groups.go
package main

import (
	"github.com/spf13/cobra"
	"github.com/schmoli/cli-tools/keycloak/pkg/keycloak"
)

var groupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Manage groups",
}

var groupsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List groups in realm",
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		groups, err := client.ListGroups(realm)
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(groups); err != nil {
			handleError(err)
		}
	},
}

var groupsGetCmd = &cobra.Command{
	Use:   "get <group-id>",
	Short: "Get group details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		group, err := client.GetGroup(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(group); err != nil {
			handleError(err)
		}
	},
}

var groupsMembersCmd = &cobra.Command{
	Use:   "members <group-id>",
	Short: "List group members",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		realm := getTargetRealm()
		if realm == "" {
			handleError(keycloak.ConfigError("missing --target-realm or KEYCLOAK_TARGET_REALM"))
			return
		}

		client, err := getClient()
		if err != nil {
			handleError(err)
			return
		}

		members, err := client.GetGroupMembers(realm, args[0])
		if err != nil {
			handleError(err)
			return
		}

		if err := keycloak.PrintYAML(members); err != nil {
			handleError(err)
		}
	},
}

func init() {
	groupsCmd.AddCommand(groupsListCmd)
	groupsCmd.AddCommand(groupsGetCmd)
	groupsCmd.AddCommand(groupsMembersCmd)
}
```

**Step 3: Add groupsCmd to main.go init()**

```go
rootCmd.AddCommand(groupsCmd)
```

**Step 4: Verify compiles**

```bash
cd keycloak && go build ./cmd/keycloak-cli/...
```

**Step 5: Commit**

```bash
git add keycloak/pkg/keycloak/client.go keycloak/cmd/keycloak-cli/groups.go keycloak/cmd/keycloak-cli/main.go
git commit -m "feat(keycloak): add groups list/get/members commands"
```

---

## Task 12: Update Build Scripts

**Files:**
- Modify: `build.sh`
- Modify: `go.work`

**Step 1: Add keycloak to build.sh**

Add `keycloak` to the TOOLS array.

**Step 2: Verify build**

```bash
./build.sh
ls bin/ | grep keycloak
```

**Step 3: Commit**

```bash
git add build.sh go.work
git commit -m "chore: add keycloak-cli to build"
```

---

## Task 13: Final Verification

**Step 1: Run all tests**

```bash
go test ./...
```

**Step 2: Build all**

```bash
./build.sh
```

**Step 3: Verify CLI help**

```bash
./bin/keycloak-cli --help
./bin/keycloak-cli realms --help
./bin/keycloak-cli users --help
```

**Step 4: Final commit if needed**

---

## Summary

11 implementation tasks + 2 verification tasks. Each task is self-contained with TDD where applicable.

Commands implemented:
- `realms list`, `realms get`
- `users list`, `users get`, `users sessions`
- `clients list`, `clients get`, `clients sessions`
- `roles list`, `roles get` (with `--client` for client roles)
- `groups list`, `groups get`, `groups members`
