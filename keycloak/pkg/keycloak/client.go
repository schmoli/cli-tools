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
