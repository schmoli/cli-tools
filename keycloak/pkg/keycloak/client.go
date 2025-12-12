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
