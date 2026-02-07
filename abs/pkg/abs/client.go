package abs

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewClient(baseURL, token string, insecure bool) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		token:      token,
		httpClient: client,
	}
}

func (c *Client) request(method, path string, result interface{}) error {
	reqURL := c.baseURL + path
	
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return NetworkError(err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NetworkError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return AuthError("invalid or expired token")
	}

	if resp.StatusCode == http.StatusNotFound {
		return NotFoundError("resource not found")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return APIError(fmt.Sprintf("unexpected status %d", resp.StatusCode))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return APIError(fmt.Sprintf("failed to parse response: %s", err))
		}
	}

	return nil
}

func (c *Client) ListLibraries() ([]APILibrary, error) {
	var resp struct {
		Libraries []APILibrary `json:"libraries"`
	}
	if err := c.request("GET", "/api/libraries", &resp); err != nil {
		return nil, err
	}
	return resp.Libraries, nil
}

func (c *Client) ListLibraryItems(libraryID string, limit int) ([]APILibraryItem, int, error) {
	path := fmt.Sprintf("/api/libraries/%s/items?limit=%d&sort=media.metadata.title", libraryID, limit)
	
	var resp APILibraryItemsResponse
	if err := c.request("GET", path, &resp); err != nil {
		return nil, 0, err
	}
	return resp.Results, resp.Total, nil
}

func (c *Client) GetItem(itemID string) (*APILibraryItem, error) {
	var item APILibraryItem
	if err := c.request("GET", "/api/items/"+itemID, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (c *Client) GetProgress() ([]APIMediaProgress, error) {
	var resp struct {
		ID            string             `json:"id"`
		MediaProgress []APIMediaProgress `json:"mediaProgress"`
	}
	if err := c.request("GET", "/api/me", &resp); err != nil {
		return nil, err
	}
	return resp.MediaProgress, nil
}

func (c *Client) Search(libraryID, query string) (*APISearchResponse, error) {
	path := fmt.Sprintf("/api/libraries/%s/search?q=%s", libraryID, url.QueryEscape(query))
	
	var resp APISearchResponse
	if err := c.request("GET", path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) ScanLibrary(libraryID string) error {
	return c.request("POST", "/api/libraries/"+libraryID+"/scan", nil)
}
