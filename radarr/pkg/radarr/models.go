package radarr

import (
	"fmt"
	"time"
)

// API response types from Radarr

type APIMovie struct {
	ID                int        `json:"id"`
	Title             string     `json:"title"`
	OriginalTitle     string     `json:"originalTitle"`
	SortTitle         string     `json:"sortTitle"`
	Status            string     `json:"status"`
	Overview          string     `json:"overview"`
	Year              int        `json:"year"`
	Runtime           int        `json:"runtime"`
	Path              string     `json:"path"`
	SizeOnDisk        int64      `json:"sizeOnDisk"`
	Genres            []string   `json:"genres"`
	Added             time.Time  `json:"added"`
	HasFile           bool       `json:"hasFile"`
	Monitored         bool       `json:"monitored"`
	IsAvailable       bool       `json:"isAvailable"`
	InCinemas         *time.Time `json:"inCinemas,omitempty"`
	PhysicalRelease   *time.Time `json:"physicalRelease,omitempty"`
	DigitalRelease    *time.Time `json:"digitalRelease,omitempty"`
	Studio            string     `json:"studio"`
	TmdbID            int        `json:"tmdbId"`
	ImdbID            string     `json:"imdbId"`
}

type APICalendarEntry struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	Year            int        `json:"year"`
	Status          string     `json:"status"`
	HasFile         bool       `json:"hasFile"`
	InCinemas       *time.Time `json:"inCinemas,omitempty"`
	PhysicalRelease *time.Time `json:"physicalRelease,omitempty"`
	DigitalRelease  *time.Time `json:"digitalRelease,omitempty"`
}

type APIQueueItem struct {
	ID                      int       `json:"id"`
	MovieID                 int       `json:"movieId"`
	Title                   string    `json:"title"`
	Status                  string    `json:"status"`
	Size                    int64     `json:"size"`
	Sizeleft                int64     `json:"sizeleft"`
	Timeleft                string    `json:"timeleft"`
	EstimatedCompletionTime *time.Time `json:"estimatedCompletionTime,omitempty"`
	Movie                   APIMovie  `json:"movie"`
}

type APIQueueResponse struct {
	Page         int            `json:"page"`
	PageSize     int            `json:"pageSize"`
	TotalRecords int            `json:"totalRecords"`
	Records      []APIQueueItem `json:"records"`
}

type APIWantedResponse struct {
	Page         int        `json:"page"`
	PageSize     int        `json:"pageSize"`
	TotalRecords int        `json:"totalRecords"`
	Records      []APIMovie `json:"records"`
}

type APISearchResult struct {
	Title    string   `json:"title"`
	Year     int      `json:"year"`
	TmdbID   int      `json:"tmdbId"`
	ImdbID   string   `json:"imdbId"`
	Overview string   `json:"overview"`
	Runtime  int      `json:"runtime"`
	Genres   []string `json:"genres"`
	Studio   string   `json:"studio"`
}

// Output types for YAML

type MovieList struct {
	Movies []MovieListItem `yaml:"movies"`
}

type MovieListItem struct {
	ID       int    `yaml:"id"`
	Title    string `yaml:"title"`
	Year     int    `yaml:"year"`
	Status   string `yaml:"status"`
	HasFile  bool   `yaml:"hasFile"`
	Size     string `yaml:"size"`
	Runtime  string `yaml:"runtime"`
}

type MovieDetail struct {
	Movie MovieDetailItem `yaml:"movie"`
}

type MovieDetailItem struct {
	ID          int      `yaml:"id"`
	Title       string   `yaml:"title"`
	Year        int      `yaml:"year"`
	Status      string   `yaml:"status"`
	HasFile     bool     `yaml:"hasFile"`
	Size        string   `yaml:"size"`
	Runtime     string   `yaml:"runtime"`
	Path        string   `yaml:"path"`
	Genres      []string `yaml:"genres,omitempty"`
	Studio      string   `yaml:"studio,omitempty"`
	ImdbID      string   `yaml:"imdbId,omitempty"`
	InCinemas   string   `yaml:"inCinemas,omitempty"`
	PhysicalRelease string `yaml:"physicalRelease,omitempty"`
	Overview    string   `yaml:"overview,omitempty"`
}

type CalendarList struct {
	Movies []CalendarItem `yaml:"movies"`
}

type CalendarItem struct {
	Title       string `yaml:"title"`
	Year        int    `yaml:"year"`
	ReleaseDate string `yaml:"releaseDate"`
	ReleaseType string `yaml:"releaseType"`
	HasFile     bool   `yaml:"hasFile"`
}

type QueueList struct {
	Queue []QueueItem `yaml:"queue"`
	Total int         `yaml:"total"`
}

type QueueItem struct {
	Title     string `yaml:"title"`
	Year      int    `yaml:"year"`
	Status    string `yaml:"status"`
	Size      string `yaml:"size"`
	Remaining string `yaml:"remaining"`
	TimeLeft  string `yaml:"timeLeft,omitempty"`
}

type WantedList struct {
	Movies []WantedItem `yaml:"movies"`
	Total  int          `yaml:"total"`
}

type WantedItem struct {
	ID      int    `yaml:"id"`
	Title   string `yaml:"title"`
	Year    int    `yaml:"year"`
	Status  string `yaml:"status"`
}

type SearchResultList struct {
	Results []SearchResultItem `yaml:"results"`
}

type SearchResultItem struct {
	Title   string `yaml:"title"`
	Year    int    `yaml:"year"`
	TmdbID  int    `yaml:"tmdbId"`
	Runtime string `yaml:"runtime"`
	Studio  string `yaml:"studio,omitempty"`
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

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Local().Format("2006-01-02")
}

func formatRuntime(minutes int) string {
	if minutes == 0 {
		return "-"
	}
	h := minutes / 60
	m := minutes % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// Conversion methods

func (m *APIMovie) ToListItem() MovieListItem {
	return MovieListItem{
		ID:      m.ID,
		Title:   m.Title,
		Year:    m.Year,
		Status:  m.Status,
		HasFile: m.HasFile,
		Size:    formatBytes(m.SizeOnDisk),
		Runtime: formatRuntime(m.Runtime),
	}
}

func (m *APIMovie) ToDetail() MovieDetailItem {
	item := MovieDetailItem{
		ID:      m.ID,
		Title:   m.Title,
		Year:    m.Year,
		Status:  m.Status,
		HasFile: m.HasFile,
		Size:    formatBytes(m.SizeOnDisk),
		Runtime: formatRuntime(m.Runtime),
		Path:    m.Path,
	}
	if len(m.Genres) > 0 {
		item.Genres = m.Genres
	}
	if m.Studio != "" {
		item.Studio = m.Studio
	}
	if m.ImdbID != "" {
		item.ImdbID = m.ImdbID
	}
	if m.InCinemas != nil {
		item.InCinemas = formatTime(m.InCinemas)
	}
	if m.PhysicalRelease != nil {
		item.PhysicalRelease = formatTime(m.PhysicalRelease)
	}
	if m.Overview != "" {
		overview := m.Overview
		if len(overview) > 500 {
			overview = overview[:497] + "..."
		}
		item.Overview = overview
	}
	return item
}

func (c *APICalendarEntry) ToListItem() CalendarItem {
	item := CalendarItem{
		Title:   c.Title,
		Year:    c.Year,
		HasFile: c.HasFile,
	}
	// Prefer digital, then physical, then cinema
	if c.DigitalRelease != nil {
		item.ReleaseDate = formatTime(c.DigitalRelease)
		item.ReleaseType = "digital"
	} else if c.PhysicalRelease != nil {
		item.ReleaseDate = formatTime(c.PhysicalRelease)
		item.ReleaseType = "physical"
	} else if c.InCinemas != nil {
		item.ReleaseDate = formatTime(c.InCinemas)
		item.ReleaseType = "cinema"
	}
	return item
}

func (q *APIQueueItem) ToListItem() QueueItem {
	item := QueueItem{
		Title:     q.Movie.Title,
		Year:      q.Movie.Year,
		Status:    q.Status,
		Size:      formatBytes(q.Size),
		Remaining: formatBytes(q.Sizeleft),
	}
	if q.Timeleft != "" {
		item.TimeLeft = q.Timeleft
	}
	return item
}

func (m *APIMovie) ToWantedItem() WantedItem {
	return WantedItem{
		ID:     m.ID,
		Title:  m.Title,
		Year:   m.Year,
		Status: m.Status,
	}
}

func (r *APISearchResult) ToListItem() SearchResultItem {
	item := SearchResultItem{
		Title:   r.Title,
		Year:    r.Year,
		TmdbID:  r.TmdbID,
		Runtime: formatRuntime(r.Runtime),
	}
	if r.Studio != "" {
		item.Studio = r.Studio
	}
	return item
}
