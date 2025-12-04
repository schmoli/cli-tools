package nproxy

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(url, token string, insecure bool) *Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &Client{
		baseURL:    strings.TrimSuffix(url, "/"),
		token:      token,
		httpClient: client,
	}
}

func (c *Client) get(path string, result interface{}) error {
	url := c.baseURL + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NetworkError(err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NetworkError(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return APIError(fmt.Sprintf("failed to parse response from %s: %s", path, err))
		}
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return AuthError(fmt.Sprintf("invalid or expired token for %s", path))
	case http.StatusNotFound:
		return NotFoundError(fmt.Sprintf("resource not found: %s", path))
	default:
		return APIError(fmt.Sprintf("unexpected status %d from %s", resp.StatusCode, path))
	}
}

// Login authenticates and returns a token
func Login(url, email, password string, insecure bool) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	payload := map[string]string{
		"identity": email,
		"secret":   password,
	}
	body, _ := json.Marshal(payload)

	reqURL := strings.TrimSuffix(url, "/") + "/api/tokens"
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return "", NetworkError(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", NetworkError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", AuthError("invalid credentials")
	}
	if resp.StatusCode != http.StatusOK {
		return "", APIError(fmt.Sprintf("login failed with status %d", resp.StatusCode))
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", APIError("failed to parse login response")
	}

	return result.Token, nil
}

func (c *Client) ListProxyHosts() ([]APIProxyHost, error) {
	var hosts []APIProxyHost
	if err := c.get("/api/nginx/proxy-hosts", &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

func (c *Client) GetProxyHost(id int64) (*APIProxyHost, error) {
	var host APIProxyHost
	path := fmt.Sprintf("/api/nginx/proxy-hosts/%d", id)
	if err := c.get(path, &host); err != nil {
		return nil, err
	}
	return &host, nil
}

func (c *Client) ListCertificates() ([]APICertificate, error) {
	var certs []APICertificate
	if err := c.get("/api/nginx/certificates", &certs); err != nil {
		return nil, err
	}
	return certs, nil
}

func (c *Client) GetCertificate(id int64) (*APICertificate, error) {
	var cert APICertificate
	path := fmt.Sprintf("/api/nginx/certificates/%d", id)
	if err := c.get(path, &cert); err != nil {
		return nil, err
	}
	return &cert, nil
}
