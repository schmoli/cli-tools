package pve

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL     string
	tokenID     string
	tokenSecret string
	httpClient  *http.Client
	node        string // cached node name
}

func NewClient(url, tokenID, tokenSecret string, insecure bool) *Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return &Client{
		baseURL:     strings.TrimSuffix(url, "/"),
		tokenID:     tokenID,
		tokenSecret: tokenSecret,
		httpClient:  client,
	}
}

func (c *Client) authHeader() string {
	return fmt.Sprintf("PVEAPIToken=%s=%s", c.tokenID, c.tokenSecret)
}

func (c *Client) get(path string, result interface{}) error {
	url := c.baseURL + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NetworkError(err.Error())
	}

	req.Header.Set("Authorization", c.authHeader())

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

func (c *Client) post(path string, result interface{}) error {
	url := c.baseURL + path

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return NetworkError(err.Error())
	}

	req.Header.Set("Authorization", c.authHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NetworkError(err.Error())
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if result != nil {
			if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
				return APIError(fmt.Sprintf("failed to parse response from %s: %s", path, err))
			}
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

func (c *Client) GetNode() (string, error) {
	if c.node != "" {
		return c.node, nil
	}

	var resp APINodesResponse
	if err := c.get("/api2/json/nodes", &resp); err != nil {
		return "", err
	}

	if len(resp.Data) == 0 {
		return "", APIError("no nodes found")
	}

	c.node = resp.Data[0].Node
	return c.node, nil
}
