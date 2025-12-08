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

func (c *Client) ListGuests() ([]Guest, error) {
	node, err := c.GetNode()
	if err != nil {
		return nil, err
	}

	var guests []Guest

	// Get VMs
	var vmResp APIVMListResponse
	if err := c.get(fmt.Sprintf("/api2/json/nodes/%s/qemu", node), &vmResp); err != nil {
		return nil, err
	}
	for _, vm := range vmResp.Data {
		ip := c.getVMIP(node, vm.VMID)
		guests = append(guests, vm.ToGuest("qemu", ip))
	}

	// Get LXCs
	var lxcResp APIVMListResponse
	if err := c.get(fmt.Sprintf("/api2/json/nodes/%s/lxc", node), &lxcResp); err != nil {
		return nil, err
	}
	for _, lxc := range lxcResp.Data {
		ip := c.getLXCIP(node, lxc.VMID)
		guests = append(guests, lxc.ToGuest("lxc", ip))
	}

	return guests, nil
}

func (c *Client) getVMIP(node string, vmid int64) string {
	var resp APIQemuAgentNetworkResponse
	path := fmt.Sprintf("/api2/json/nodes/%s/qemu/%d/agent/network-get-interfaces", node, vmid)
	if err := c.get(path, &resp); err != nil {
		return "N/A"
	}

	for _, iface := range resp.Data.Result {
		if iface.Name == "lo" {
			continue
		}
		for _, addr := range iface.IPAddresses {
			if addr.IPType == "ipv4" && !strings.HasPrefix(addr.IPAddress, "127.") {
				return addr.IPAddress
			}
		}
	}
	return "N/A"
}

func (c *Client) getLXCIP(node string, vmid int64) string {
	var resp APILXCInterfaceResponse
	path := fmt.Sprintf("/api2/json/nodes/%s/lxc/%d/interfaces", node, vmid)
	if err := c.get(path, &resp); err != nil {
		return "N/A"
	}

	for _, iface := range resp.Data {
		if iface.Name == "lo" {
			continue
		}
		if iface.Inet != "" {
			// Strip CIDR notation
			ip := strings.Split(iface.Inet, "/")[0]
			if !strings.HasPrefix(ip, "127.") {
				return ip
			}
		}
	}
	return "N/A"
}
