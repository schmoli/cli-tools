package nproxy

// API response types (match nginx-proxy-manager JSON)
type APIProxyHost struct {
	ID               int64    `json:"id"`
	DomainNames      []string `json:"domain_names"`
	ForwardScheme    string   `json:"forward_scheme"`
	ForwardHost      string   `json:"forward_host"`
	ForwardPort      int      `json:"forward_port"`
	CertificateID    *int64   `json:"certificate_id"`
	SSLForced        bool     `json:"ssl_forced"`
	HTTPSRedirect    bool     `json:"http2_support"`
	BlockExploits    bool     `json:"block_exploits"`
	CachingEnabled   bool     `json:"caching_enabled"`
	AllowWebsocket   bool     `json:"allow_websocket_upgrade"`
	AccessListID     int64    `json:"access_list_id"`
	AdvancedConfig   string   `json:"advanced_config"`
	Enabled          bool     `json:"enabled"`
	Meta             APIMeta  `json:"meta"`
}

type APICertificate struct {
	ID            int64    `json:"id"`
	Provider      string   `json:"provider"`
	NiceName      string   `json:"nice_name"`
	DomainNames   []string `json:"domain_names"`
	ExpiresOn     string   `json:"expires_on"`
	Meta          APIMeta  `json:"meta"`
}

type APIMeta struct {
	LetsencryptEmail   string `json:"letsencrypt_email,omitempty"`
	LetsencryptAgree   bool   `json:"letsencrypt_agree,omitempty"`
	DNSChallenge       bool   `json:"dns_challenge,omitempty"`
	DNSProvider        string `json:"dns_provider,omitempty"`
}

// Output types (curated, YAML output)
type ProxyHost struct {
	ID             int64    `yaml:"id"`
	DomainNames    []string `yaml:"domainNames"`
	ForwardScheme  string   `yaml:"forwardScheme"`
	ForwardHost    string   `yaml:"forwardHost"`
	ForwardPort    int      `yaml:"forwardPort"`
	CertificateID  *int64   `yaml:"certificateId,omitempty"`
	SSLForced      bool     `yaml:"sslForced"`
	BlockExploits  bool     `yaml:"blockExploits"`
	CachingEnabled bool     `yaml:"cachingEnabled"`
	Websocket      bool     `yaml:"websocket"`
	Enabled        bool     `yaml:"enabled"`
	AdvancedConfig string   `yaml:"advancedConfig,omitempty"`
}

type ProxyHostList struct {
	Hosts []ProxyHostListItem `yaml:"hosts"`
}

type ProxyHostListItem struct {
	ID            int64    `yaml:"id"`
	DomainNames   []string `yaml:"domainNames"`
	ForwardHost   string   `yaml:"forwardHost"`
	ForwardPort   int      `yaml:"forwardPort"`
	SSLForced     bool     `yaml:"sslForced"`
	Enabled       bool     `yaml:"enabled"`
}

type Certificate struct {
	ID          int64    `yaml:"id"`
	Provider    string   `yaml:"provider"`
	NiceName    string   `yaml:"niceName"`
	DomainNames []string `yaml:"domainNames"`
	ExpiresOn   string   `yaml:"expiresOn"`
}

type CertificateList struct {
	Certificates []CertificateListItem `yaml:"certificates"`
}

type CertificateListItem struct {
	ID        int64    `yaml:"id"`
	NiceName  string   `yaml:"niceName"`
	Provider  string   `yaml:"provider"`
	ExpiresOn string   `yaml:"expiresOn"`
}

// Mapping functions
func (h *APIProxyHost) ToListItem() ProxyHostListItem {
	return ProxyHostListItem{
		ID:          h.ID,
		DomainNames: h.DomainNames,
		ForwardHost: h.ForwardHost,
		ForwardPort: h.ForwardPort,
		SSLForced:   h.SSLForced,
		Enabled:     h.Enabled,
	}
}

func (h *APIProxyHost) ToProxyHost() ProxyHost {
	return ProxyHost{
		ID:             h.ID,
		DomainNames:    h.DomainNames,
		ForwardScheme:  h.ForwardScheme,
		ForwardHost:    h.ForwardHost,
		ForwardPort:    h.ForwardPort,
		CertificateID:  h.CertificateID,
		SSLForced:      h.SSLForced,
		BlockExploits:  h.BlockExploits,
		CachingEnabled: h.CachingEnabled,
		Websocket:      h.AllowWebsocket,
		Enabled:        h.Enabled,
		AdvancedConfig: h.AdvancedConfig,
	}
}

func (c *APICertificate) ToListItem() CertificateListItem {
	return CertificateListItem{
		ID:        c.ID,
		NiceName:  c.NiceName,
		Provider:  c.Provider,
		ExpiresOn: c.ExpiresOn,
	}
}

func (c *APICertificate) ToCertificate() Certificate {
	return Certificate{
		ID:          c.ID,
		Provider:    c.Provider,
		NiceName:    c.NiceName,
		DomainNames: c.DomainNames,
		ExpiresOn:   c.ExpiresOn,
	}
}
