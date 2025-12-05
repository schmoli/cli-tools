package trans

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const sessionHeader = "X-Transmission-Session-Id"

// Fields requested for list view
var listFields = []string{
	"id", "name", "status",
	"percentDone", "totalSize", "sizeWhenDone",
	"uploadRatio", "rateDownload", "rateUpload",
	"eta", "peersConnected", "trackers",
}

// Additional fields for detail view
var detailFields = append(listFields,
	"downloadedEver", "uploadedEver",
	"addedDate", "doneDate", "downloadDir",
)

type Client struct {
	rpcURL     string
	user       string
	pass       string
	sessionID  string
	httpClient *http.Client
}

func NewClient(url, user, pass string, insecure bool) *Client {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	rpcURL := strings.TrimSuffix(url, "/") + "/transmission/rpc"

	return &Client{
		rpcURL:     rpcURL,
		user:       user,
		pass:       pass,
		httpClient: client,
	}
}

func (c *Client) rpc(req *RPCRequest, result interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return APIError(fmt.Sprintf("failed to encode request: %s", err))
	}

	// Try request, retry once on 409 (session ID expired)
	for attempt := 0; attempt < 2; attempt++ {
		httpReq, err := http.NewRequest("POST", c.rpcURL, bytes.NewReader(body))
		if err != nil {
			return NetworkError(err.Error())
		}

		httpReq.Header.Set("Content-Type", "application/json")
		if c.sessionID != "" {
			httpReq.Header.Set(sessionHeader, c.sessionID)
		}
		if c.user != "" {
			httpReq.SetBasicAuth(c.user, c.pass)
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			return NetworkError(err.Error())
		}
		defer resp.Body.Close()

		// Handle 409 - need new session ID
		if resp.StatusCode == http.StatusConflict {
			c.sessionID = resp.Header.Get(sessionHeader)
			if c.sessionID == "" {
				return APIError("received 409 but no session ID in response")
			}
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return AuthError("invalid credentials")
		}

		if resp.StatusCode != http.StatusOK {
			return APIError(fmt.Sprintf("unexpected status %d", resp.StatusCode))
		}

		// Parse response
		var rpcResp RPCResponse
		if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
			return APIError(fmt.Sprintf("failed to parse response: %s", err))
		}

		if rpcResp.Result != "success" {
			return APIError(rpcResp.Result)
		}

		if result != nil && rpcResp.Arguments != nil {
			if err := json.Unmarshal(rpcResp.Arguments, result); err != nil {
				return APIError(fmt.Sprintf("failed to parse arguments: %s", err))
			}
		}

		return nil
	}

	return APIError("failed to get valid session after retry")
}

func (c *Client) ListTorrents() ([]APITorrent, error) {
	req := &RPCRequest{
		Method: "torrent-get",
		Arguments: TorrentGetArgs{
			Fields: listFields,
		},
	}

	var resp TorrentGetResponse
	if err := c.rpc(req, &resp); err != nil {
		return nil, err
	}
	return resp.Torrents, nil
}

func (c *Client) GetTorrent(id int64) (*APITorrent, error) {
	req := &RPCRequest{
		Method: "torrent-get",
		Arguments: TorrentGetArgs{
			Fields: detailFields,
			IDs:    []int64{id},
		},
	}

	var resp TorrentGetResponse
	if err := c.rpc(req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Torrents) == 0 {
		return nil, NotFoundError(fmt.Sprintf("torrent %d not found", id))
	}

	return &resp.Torrents[0], nil
}

func (c *Client) StartTorrent(id int64) error {
	req := &RPCRequest{
		Method: "torrent-start",
		Arguments: TorrentActionArgs{
			IDs: []int64{id},
		},
	}
	return c.rpc(req, nil)
}

func (c *Client) StopTorrent(id int64) error {
	req := &RPCRequest{
		Method: "torrent-stop",
		Arguments: TorrentActionArgs{
			IDs: []int64{id},
		},
	}
	return c.rpc(req, nil)
}

func (c *Client) AddTorrentMagnet(magnet string) (*TorrentAddedInfo, error) {
	req := &RPCRequest{
		Method: "torrent-add",
		Arguments: TorrentAddArgs{
			Filename: magnet,
		},
	}

	var resp TorrentAddResponse
	if err := c.rpc(req, &resp); err != nil {
		return nil, err
	}

	if resp.TorrentAdded != nil {
		return resp.TorrentAdded, nil
	}
	if resp.TorrentDuplicate != nil {
		return resp.TorrentDuplicate, nil
	}
	return nil, APIError("no torrent info in response")
}

func (c *Client) AddTorrentFile(path string) (*TorrentAddedInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, ConfigError(fmt.Sprintf("failed to read file: %s", err))
	}

	req := &RPCRequest{
		Method: "torrent-add",
		Arguments: TorrentAddArgs{
			Metainfo: base64.StdEncoding.EncodeToString(data),
		},
	}

	var resp TorrentAddResponse
	if err := c.rpc(req, &resp); err != nil {
		return nil, err
	}

	if resp.TorrentAdded != nil {
		return resp.TorrentAdded, nil
	}
	if resp.TorrentDuplicate != nil {
		return resp.TorrentDuplicate, nil
	}
	return nil, APIError("no torrent info in response")
}
