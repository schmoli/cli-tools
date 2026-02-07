package abs

import (
	"fmt"
	"time"
)

// API response types from Audiobookshelf

type APILibrary struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	MediaType    string `json:"mediaType"`
	Folders      []APIFolder `json:"folders"`
	CreatedAt    int64  `json:"createdAt"`
}

type APIFolder struct {
	ID       string `json:"id"`
	FullPath string `json:"fullPath"`
}

type APILibraryItemsResponse struct {
	Results []APILibraryItem `json:"results"`
	Total   int              `json:"total"`
	Limit   int              `json:"limit"`
	Page    int              `json:"page"`
}

type APILibraryItem struct {
	ID            string   `json:"id"`
	LibraryID     string   `json:"libraryId"`
	Media         APIMedia `json:"media"`
	NumFiles      int      `json:"numFiles"`
	Size          int64    `json:"size"`
	AddedAt       int64    `json:"addedAt"`
	UpdatedAt     int64    `json:"updatedAt"`
}

type APIMedia struct {
	Metadata    APIMetadata    `json:"metadata"`
	Duration    float64        `json:"duration"`
	NumTracks   int            `json:"numTracks,omitempty"`
	NumChapters int            `json:"numChapters,omitempty"`
}

type APIMetadata struct {
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle,omitempty"`
	AuthorName  string   `json:"authorName,omitempty"`
	Authors     []APIAuthor `json:"authors,omitempty"`
	Narrators   []string `json:"narrators,omitempty"`
	SeriesName  string   `json:"seriesName,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	PublishedYear string `json:"publishedYear,omitempty"`
	Description string   `json:"description,omitempty"`
}

type APIAuthor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type APIMediaProgress struct {
	ID                string  `json:"id"`
	LibraryItemID     string  `json:"libraryItemId"`
	EpisodeID         string  `json:"episodeId,omitempty"`
	Duration          float64 `json:"duration"`
	Progress          float64 `json:"progress"`
	CurrentTime       float64 `json:"currentTime"`
	IsFinished        bool    `json:"isFinished"`
	LastUpdate        int64   `json:"lastUpdate"`
	StartedAt         int64   `json:"startedAt"`
	FinishedAt        int64   `json:"finishedAt,omitempty"`
}

type APISearchResponse struct {
	Book     []APISearchResult `json:"book,omitempty"`
	Podcast  []APISearchResult `json:"podcast,omitempty"`
	Authors  []APIAuthor       `json:"authors,omitempty"`
	Series   []APISeriesResult `json:"series,omitempty"`
}

type APISearchResult struct {
	LibraryItem APILibraryItem `json:"libraryItem"`
	MatchKey    string         `json:"matchKey,omitempty"`
	MatchText   string         `json:"matchText,omitempty"`
}

type APISeriesResult struct {
	Series APISeriesInfo `json:"series"`
}

type APISeriesInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type APIScanResponse struct {
	Success bool `json:"success"`
}

// Output types for YAML

type LibraryList struct {
	Libraries []LibraryListItem `yaml:"libraries"`
}

type LibraryListItem struct {
	ID        string `yaml:"id"`
	Name      string `yaml:"name"`
	MediaType string `yaml:"mediaType"`
	Folders   int    `yaml:"folders"`
}

type BookList struct {
	Books []BookListItem `yaml:"books"`
	Total int            `yaml:"total"`
}

type BookListItem struct {
	ID       string `yaml:"id"`
	Title    string `yaml:"title"`
	Author   string `yaml:"author"`
	Duration string `yaml:"duration"`
	Size     string `yaml:"size"`
}

type BookDetail struct {
	Book BookDetailItem `yaml:"book"`
}

type BookDetailItem struct {
	ID            string   `yaml:"id"`
	Title         string   `yaml:"title"`
	Subtitle      string   `yaml:"subtitle,omitempty"`
	Author        string   `yaml:"author"`
	Narrators     []string `yaml:"narrators,omitempty"`
	Series        string   `yaml:"series,omitempty"`
	Duration      string   `yaml:"duration"`
	Chapters      int      `yaml:"chapters,omitempty"`
	Size          string   `yaml:"size"`
	PublishedYear string   `yaml:"publishedYear,omitempty"`
	AddedAt       string   `yaml:"addedAt"`
	Description   string   `yaml:"description,omitempty"`
}

type ProgressList struct {
	Progress []ProgressListItem `yaml:"progress"`
}

type ProgressListItem struct {
	LibraryItemID string `yaml:"libraryItemId"`
	Title         string `yaml:"title,omitempty"`
	Progress      string `yaml:"progress"`
	CurrentTime   string `yaml:"currentTime"`
	Duration      string `yaml:"duration"`
	IsFinished    bool   `yaml:"isFinished"`
	LastUpdate    string `yaml:"lastUpdate"`
}

type SearchResults struct {
	Books   []BookListItem `yaml:"books,omitempty"`
	Authors []string       `yaml:"authors,omitempty"`
	Series  []string       `yaml:"series,omitempty"`
}

// Formatting helpers

func formatDuration(seconds float64) string {
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

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

func formatTime(unix int64) string {
	if unix == 0 {
		return "-"
	}
	return time.Unix(unix/1000, 0).Format("2006-01-02 15:04")
}

func formatPercent(p float64) string {
	return fmt.Sprintf("%.0f%%", p*100)
}

func formatTimePosition(seconds float64) string {
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

func getAuthorName(m APIMetadata) string {
	if m.AuthorName != "" {
		return m.AuthorName
	}
	if len(m.Authors) > 0 {
		return m.Authors[0].Name
	}
	return "-"
}

// Conversion methods

func (l *APILibrary) ToListItem() LibraryListItem {
	return LibraryListItem{
		ID:        l.ID,
		Name:      l.Name,
		MediaType: l.MediaType,
		Folders:   len(l.Folders),
	}
}

func (i *APILibraryItem) ToListItem() BookListItem {
	return BookListItem{
		ID:       i.ID,
		Title:    i.Media.Metadata.Title,
		Author:   getAuthorName(i.Media.Metadata),
		Duration: formatDuration(i.Media.Duration),
		Size:     formatBytes(i.Size),
	}
}

func (i *APILibraryItem) ToDetail() BookDetailItem {
	item := BookDetailItem{
		ID:            i.ID,
		Title:         i.Media.Metadata.Title,
		Author:        getAuthorName(i.Media.Metadata),
		Duration:      formatDuration(i.Media.Duration),
		Chapters:      i.Media.NumChapters,
		Size:          formatBytes(i.Size),
		AddedAt:       formatTime(i.AddedAt),
	}
	if i.Media.Metadata.Subtitle != "" {
		item.Subtitle = i.Media.Metadata.Subtitle
	}
	if len(i.Media.Metadata.Narrators) > 0 {
		item.Narrators = i.Media.Metadata.Narrators
	}
	if i.Media.Metadata.SeriesName != "" {
		item.Series = i.Media.Metadata.SeriesName
	}
	if i.Media.Metadata.PublishedYear != "" {
		item.PublishedYear = i.Media.Metadata.PublishedYear
	}
	if i.Media.Metadata.Description != "" {
		// Truncate long descriptions
		desc := i.Media.Metadata.Description
		if len(desc) > 500 {
			desc = desc[:497] + "..."
		}
		item.Description = desc
	}
	return item
}

func (p *APIMediaProgress) ToListItem() ProgressListItem {
	return ProgressListItem{
		LibraryItemID: p.LibraryItemID,
		Progress:      formatPercent(p.Progress),
		CurrentTime:   formatTimePosition(p.CurrentTime),
		Duration:      formatDuration(p.Duration),
		IsFinished:    p.IsFinished,
		LastUpdate:    formatTime(p.LastUpdate),
	}
}
