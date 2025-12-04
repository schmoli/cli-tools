package nproxy

import "testing"

func TestProxyHostToListItem(t *testing.T) {
	host := APIProxyHost{
		ID:            1,
		DomainNames:   []string{"example.com"},
		ForwardHost:   "backend",
		ForwardPort:   8080,
		SSLForced:     true,
		Enabled:       true,
	}

	item := host.ToListItem()

	if item.ID != 1 {
		t.Errorf("ID = %d, want 1", item.ID)
	}
	if len(item.DomainNames) != 1 || item.DomainNames[0] != "example.com" {
		t.Errorf("DomainNames = %v, want [example.com]", item.DomainNames)
	}
	if item.ForwardHost != "backend" {
		t.Errorf("ForwardHost = %s, want backend", item.ForwardHost)
	}
	if item.ForwardPort != 8080 {
		t.Errorf("ForwardPort = %d, want 8080", item.ForwardPort)
	}
	if !item.SSLForced {
		t.Error("SSLForced = false, want true")
	}
	if !item.Enabled {
		t.Error("Enabled = false, want true")
	}
}

func TestProxyHostToProxyHost(t *testing.T) {
	certID := int64(5)
	host := APIProxyHost{
		ID:             1,
		DomainNames:    []string{"example.com"},
		ForwardScheme:  "https",
		ForwardHost:    "backend",
		ForwardPort:    8080,
		CertificateID:  &certID,
		SSLForced:      true,
		BlockExploits:  true,
		CachingEnabled: false,
		AllowWebsocket: true,
		Enabled:        true,
		AdvancedConfig: "proxy_pass http://test;",
	}

	result := host.ToProxyHost()

	if result.ID != 1 {
		t.Errorf("ID = %d, want 1", result.ID)
	}
	if result.ForwardScheme != "https" {
		t.Errorf("ForwardScheme = %s, want https", result.ForwardScheme)
	}
	if result.CertificateID == nil || *result.CertificateID != 5 {
		t.Errorf("CertificateID = %v, want 5", result.CertificateID)
	}
	if !result.Websocket {
		t.Error("Websocket = false, want true")
	}
	if result.AdvancedConfig != "proxy_pass http://test;" {
		t.Errorf("AdvancedConfig = %s, want 'proxy_pass http://test;'", result.AdvancedConfig)
	}
}

func TestCertificateToListItem(t *testing.T) {
	cert := APICertificate{
		ID:        1,
		NiceName:  "My Cert",
		Provider:  "letsencrypt",
		ExpiresOn: "2024-12-31",
	}

	item := cert.ToListItem()

	if item.ID != 1 {
		t.Errorf("ID = %d, want 1", item.ID)
	}
	if item.NiceName != "My Cert" {
		t.Errorf("NiceName = %s, want 'My Cert'", item.NiceName)
	}
	if item.Provider != "letsencrypt" {
		t.Errorf("Provider = %s, want letsencrypt", item.Provider)
	}
	if item.ExpiresOn != "2024-12-31" {
		t.Errorf("ExpiresOn = %s, want 2024-12-31", item.ExpiresOn)
	}
}

func TestCertificateToCertificate(t *testing.T) {
	cert := APICertificate{
		ID:          1,
		Provider:    "letsencrypt",
		NiceName:    "My Cert",
		DomainNames: []string{"example.com", "www.example.com"},
		ExpiresOn:   "2024-12-31",
	}

	result := cert.ToCertificate()

	if result.ID != 1 {
		t.Errorf("ID = %d, want 1", result.ID)
	}
	if len(result.DomainNames) != 2 {
		t.Errorf("DomainNames length = %d, want 2", len(result.DomainNames))
	}
}
