package trans

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// RPC request/response types
type RPCRequest struct {
	Method    string      `json:"method"`
	Arguments interface{} `json:"arguments,omitempty"`
}

type RPCResponse struct {
	Result    string          `json:"result"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

type TorrentGetArgs struct {
	Fields []string `json:"fields"`
	IDs    []int64  `json:"ids,omitempty"`
}

type TorrentGetResponse struct {
	Torrents []APITorrent `json:"torrents"`
}

type TorrentActionArgs struct {
	IDs []int64 `json:"ids"`
}

type TorrentAddArgs struct {
	Filename string   `json:"filename,omitempty"` // magnet URI
	Metainfo string   `json:"metainfo,omitempty"` // base64 torrent file
	Labels   []string `json:"labels,omitempty"`
}

type TorrentAddResponse struct {
	TorrentAdded     *TorrentAddedInfo `json:"torrent-added,omitempty"`
	TorrentDuplicate *TorrentAddedInfo `json:"torrent-duplicate,omitempty"`
}

type TorrentAddedInfo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Hash string `json:"hashString"`
}

// API types from Transmission RPC
type APITorrent struct {
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	Status         int          `json:"status"`
	PercentDone    float64      `json:"percentDone"`
	TotalSize      int64        `json:"totalSize"`
	SizeWhenDone   int64        `json:"sizeWhenDone"`
	UploadRatio    float64      `json:"uploadRatio"`
	RateDownload   int64        `json:"rateDownload"`
	RateUpload     int64        `json:"rateUpload"`
	ETA            int64        `json:"eta"`
	PeersConnected int          `json:"peersConnected"`
	Trackers       []APITracker `json:"trackers"`
	DownloadedEver int64        `json:"downloadedEver"`
	UploadedEver   int64        `json:"uploadedEver"`
	AddedDate      int64        `json:"addedDate"`
	DoneDate       int64        `json:"doneDate"`
	DownloadDir    string       `json:"downloadDir"`
}

type APITracker struct {
	Announce string `json:"announce"`
}

// Torrent status constants
const (
	StatusStopped      = 0
	StatusCheckWait    = 1
	StatusCheck        = 2
	StatusDownloadWait = 3
	StatusDownload     = 4
	StatusSeedWait     = 5
	StatusSeed         = 6
)

// Output types for YAML
type TorrentList struct {
	Torrents []TorrentListItem `yaml:"torrents"`
}

type TorrentListItem struct {
	ID             int64   `yaml:"id"`
	Name           string  `yaml:"name"`
	Status         string  `yaml:"status"`
	PercentDone    string  `yaml:"percentDone"`
	TotalSize      string  `yaml:"totalSize"`
	UploadRatio    string  `yaml:"uploadRatio"`
	RateDownload   string  `yaml:"rateDownload"`
	RateUpload     string  `yaml:"rateUpload"`
	ETA            string  `yaml:"eta"`
	Tracker        string  `yaml:"tracker"`
	PeersConnected int     `yaml:"peersConnected"`
}

type TorrentDetail struct {
	Torrent TorrentDetailItem `yaml:"torrent"`
}

type TorrentDetailItem struct {
	ID             int64  `yaml:"id"`
	Name           string `yaml:"name"`
	Status         string `yaml:"status"`
	PercentDone    string `yaml:"percentDone"`
	TotalSize      string `yaml:"totalSize"`
	DownloadedEver string `yaml:"downloadedEver"`
	UploadedEver   string `yaml:"uploadedEver"`
	UploadRatio    string `yaml:"uploadRatio"`
	RateDownload   string `yaml:"rateDownload"`
	RateUpload     string `yaml:"rateUpload"`
	ETA            string `yaml:"eta"`
	Tracker        string `yaml:"tracker"`
	PeersConnected int    `yaml:"peersConnected"`
	AddedDate      string `yaml:"addedDate"`
	DoneDate       string `yaml:"doneDate,omitempty"`
	DownloadDir    string `yaml:"downloadDir"`
}

// Status helpers
func (t *APITorrent) StatusLabel() string {
	switch t.Status {
	case StatusStopped:
		return "stopped"
	case StatusCheckWait, StatusCheck:
		return "verifying"
	case StatusDownloadWait, StatusDownload:
		return "downloading"
	case StatusSeedWait, StatusSeed:
		return "seeding"
	default:
		return "unknown"
	}
}

func (t *APITorrent) IsDownloading() bool {
	return t.Status == StatusDownloadWait || t.Status == StatusDownload
}

func (t *APITorrent) IsSeeding() bool {
	return t.Status == StatusSeedWait || t.Status == StatusSeed
}

func (t *APITorrent) IsStopped() bool {
	return t.Status == StatusStopped
}

// Formatting helpers
func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatSpeed(bps int64) string {
	return formatBytes(bps) + "/s"
}

func formatETA(seconds int64) string {
	if seconds < 0 {
		return "-"
	}
	if seconds == 0 {
		return "done"
	}
	d := time.Duration(seconds) * time.Second
	if d >= 24*time.Hour {
		days := d / (24 * time.Hour)
		hours := (d % (24 * time.Hour)) / time.Hour
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if d >= time.Hour {
		return fmt.Sprintf("%dh %dm", d/time.Hour, (d%time.Hour)/time.Minute)
	}
	if d >= time.Minute {
		return fmt.Sprintf("%dm %ds", d/time.Minute, (d%time.Minute)/time.Second)
	}
	return fmt.Sprintf("%ds", d/time.Second)
}

func formatTime(unix int64) string {
	if unix == 0 {
		return "-"
	}
	return time.Unix(unix, 0).Format("2006-01-02 15:04:05")
}

func formatPercent(p float64) string {
	return fmt.Sprintf("%.0f%%", p*100)
}

func formatRatio(r float64) string {
	if r < 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f", r)
}

func extractTracker(trackers []APITracker) string {
	if len(trackers) == 0 {
		return "-"
	}
	u, err := url.Parse(trackers[0].Announce)
	if err != nil {
		return "-"
	}
	return u.Host
}

// Conversion methods
func (t *APITorrent) ToListItem() TorrentListItem {
	return TorrentListItem{
		ID:             t.ID,
		Name:           t.Name,
		Status:         t.StatusLabel(),
		PercentDone:    formatPercent(t.PercentDone),
		TotalSize:      formatBytes(t.TotalSize),
		UploadRatio:    formatRatio(t.UploadRatio),
		RateDownload:   formatSpeed(t.RateDownload),
		RateUpload:     formatSpeed(t.RateUpload),
		ETA:            formatETA(t.ETA),
		Tracker:        extractTracker(t.Trackers),
		PeersConnected: t.PeersConnected,
	}
}

func (t *APITorrent) ToDetail() TorrentDetailItem {
	item := TorrentDetailItem{
		ID:             t.ID,
		Name:           t.Name,
		Status:         t.StatusLabel(),
		PercentDone:    formatPercent(t.PercentDone),
		TotalSize:      formatBytes(t.TotalSize),
		DownloadedEver: formatBytes(t.DownloadedEver),
		UploadedEver:   formatBytes(t.UploadedEver),
		UploadRatio:    formatRatio(t.UploadRatio),
		RateDownload:   formatSpeed(t.RateDownload),
		RateUpload:     formatSpeed(t.RateUpload),
		ETA:            formatETA(t.ETA),
		Tracker:        extractTracker(t.Trackers),
		PeersConnected: t.PeersConnected,
		AddedDate:      formatTime(t.AddedDate),
		DownloadDir:    t.DownloadDir,
	}
	if t.DoneDate > 0 {
		item.DoneDate = formatTime(t.DoneDate)
	}
	return item
}
