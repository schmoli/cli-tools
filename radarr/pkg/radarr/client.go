package radarr

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
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string, insecure bool) *Client {
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
		apiKey:     apiKey,
		httpClient: client,
	}
}

func (c *Client) request(method, path string, result interface{}) error {
	reqURL := c.baseURL + "/api/v3" + path

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return NetworkError(err.Error())
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return NetworkError(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return AuthError("invalid API key")
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

func (c *Client) ListMovies() ([]APIMovie, error) {
	var movies []APIMovie
	if err := c.request("GET", "/movie", &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (c *Client) GetMovie(id int) (*APIMovie, error) {
	var movie APIMovie
	if err := c.request("GET", fmt.Sprintf("/movie/%d", id), &movie); err != nil {
		return nil, err
	}
	return &movie, nil
}

func (c *Client) GetCalendar(days int) ([]APICalendarEntry, error) {
	start := time.Now()
	end := start.AddDate(0, 0, days)

	path := fmt.Sprintf("/calendar?start=%s&end=%s",
		url.QueryEscape(start.Format(time.RFC3339)),
		url.QueryEscape(end.Format(time.RFC3339)))

	var entries []APICalendarEntry
	if err := c.request("GET", path, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func (c *Client) GetQueue() (*APIQueueResponse, error) {
	var queue APIQueueResponse
	if err := c.request("GET", "/queue?pageSize=100", &queue); err != nil {
		return nil, err
	}
	return &queue, nil
}

func (c *Client) GetWanted(limit int) (*APIWantedResponse, error) {
	var wanted APIWantedResponse
	path := fmt.Sprintf("/wanted/missing?pageSize=%d&sortKey=digitalRelease&sortDirection=descending", limit)
	if err := c.request("GET", path, &wanted); err != nil {
		return nil, err
	}
	return &wanted, nil
}

func (c *Client) Search(term string) ([]APISearchResult, error) {
	path := fmt.Sprintf("/movie/lookup?term=%s", url.QueryEscape(term))
	var results []APISearchResult
	if err := c.request("GET", path, &results); err != nil {
		return nil, err
	}
	return results, nil
}
