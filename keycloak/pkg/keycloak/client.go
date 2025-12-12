// keycloak/pkg/keycloak/client.go
package keycloak

import (
	"context"
	"crypto/tls"

	"github.com/Nerzal/gocloak/v13"
)

type Client struct {
	gocloak *gocloak.GoCloak
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

func derefInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

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
			IPAddress:  deref(s.IPAddress),
			Started:    derefInt64(s.Start),
			LastAccess: derefInt64(s.LastAccess),
		}
	}
	return list, nil
}

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
			IPAddress:  deref(s.IPAddress),
			Started:    derefInt64(s.Start),
			LastAccess: derefInt64(s.LastAccess),
		}
	}
	return list, nil
}

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
