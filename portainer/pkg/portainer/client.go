package portainer

import (
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

	req.Header.Set("X-API-Key", c.token)

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

func (c *Client) ListStacks() ([]APIStack, error) {
	var stacks []APIStack
	if err := c.get("/api/stacks", &stacks); err != nil {
		return nil, err
	}
	return stacks, nil
}

func (c *Client) GetStack(id, endpointID int64) (*APIStack, error) {
	var stack APIStack
	path := fmt.Sprintf("/api/stacks/%d?endpointId=%d", id, endpointID)
	if err := c.get(path, &stack); err != nil {
		return nil, err
	}
	return &stack, nil
}

func (c *Client) GetStackFile(id int64) (*APIStackFile, error) {
	var file APIStackFile
	path := fmt.Sprintf("/api/stacks/%d/file", id)
	if err := c.get(path, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

func (c *Client) ListEndpoints() ([]APIEndpoint, error) {
	var endpoints []APIEndpoint
	if err := c.get("/api/endpoints", &endpoints); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (c *Client) GetEndpoint(id int64) (*APIEndpoint, error) {
	var endpoint APIEndpoint
	path := fmt.Sprintf("/api/endpoints/%d", id)
	if err := c.get(path, &endpoint); err != nil {
		return nil, err
	}
	return &endpoint, nil
}
