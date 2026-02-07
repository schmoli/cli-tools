package sonarr

import (
	"fmt"
	"time"
)

// API response types from Sonarr

type APISeries struct {
	ID            int          `json:"id"`
	Title         string       `json:"title"`
	SortTitle     string       `json:"sortTitle"`
	Status        string       `json:"status"`
	Overview      string       `json:"overview"`
	Network       string       `json:"network"`
	Year          int          `json:"year"`
	SeasonCount   int          `json:"seasonCount"`
	EpisodeCount  int          `json:"episodeCount"`
	EpisodeFileCount int       `json:"episodeFileCount"`
	Path          string       `json:"path"`
	SizeOnDisk    int64        `json:"sizeOnDisk"`
	Genres        []string     `json:"genres"`
	Added         time.Time    `json:"added"`
	Seasons       []APISeason  `json:"seasons"`
	NextAiring    *time.Time   `json:"nextAiring,omitempty"`
	PreviousAiring *time.Time  `json:"previousAiring,omitempty"`
}

type APISeason struct {
	SeasonNumber int  `json:"seasonNumber"`
	Monitored    bool `json:"monitored"`
}

type APIEpisode struct {
	ID                  int        `json:"id"`
	SeriesID            int        `json:"seriesId"`
	SeasonNumber        int        `json:"seasonNumber"`
	EpisodeNumber       int        `json:"episodeNumber"`
	Title               string     `json:"title"`
	AirDate             string     `json:"airDate"`
	AirDateUtc          *time.Time `json:"airDateUtc,omitempty"`
	HasFile             bool       `json:"hasFile"`
	Monitored           bool       `json:"monitored"`
	Overview            string     `json:"overview"`
}

type APICalendarEntry struct {
	ID            int        `json:"id"`
	SeriesID      int        `json:"seriesId"`
	SeasonNumber  int        `json:"seasonNumber"`
	EpisodeNumber int        `json:"episodeNumber"`
	Title         string     `json:"title"`
	AirDateUtc    time.Time  `json:"airDateUtc"`
	HasFile       bool       `json:"hasFile"`
	SeriesTitle   string     `json:"-"` // Populated by client
}

type APIQueueItem struct {
	ID                  int       `json:"id"`
	SeriesID            int       `json:"seriesId"`
	Title               string    `json:"title"`
	Status              string    `json:"status"`
	Size                int64     `json:"size"`
	Sizeleft            int64     `json:"sizeleft"`
	Timeleft            string    `json:"timeleft"`
	EstimatedCompletionTime *time.Time `json:"estimatedCompletionTime,omitempty"`
	Series              APISeries `json:"series"`
	Episode             APIEpisode `json:"episode"`
}

type APIQueueResponse struct {
	Page          int            `json:"page"`
	PageSize      int            `json:"pageSize"`
	TotalRecords  int            `json:"totalRecords"`
	Records       []APIQueueItem `json:"records"`
}

type APIWantedResponse struct {
	Page         int          `json:"page"`
	PageSize     int          `json:"pageSize"`
	TotalRecords int          `json:"totalRecords"`
	Records      []APIEpisode `json:"records"`
}

type APISearchResult struct {
	Title       string `json:"title"`
	Year        int    `json:"year"`
	TvdbID      int    `json:"tvdbId"`
	Overview    string `json:"overview"`
	Network     string `json:"network"`
	SeasonCount int    `json:"seasonCount"`
}

// Output types for YAML

type SeriesList struct {
	Series []SeriesListItem `yaml:"series"`
}

type SeriesListItem struct {
	ID           int    `yaml:"id"`
	Title        string `yaml:"title"`
	Year         int    `yaml:"year"`
	Status       string `yaml:"status"`
	Network      string `yaml:"network"`
	Seasons      int    `yaml:"seasons"`
	Episodes     string `yaml:"episodes"`
	Size         string `yaml:"size"`
	NextAiring   string `yaml:"nextAiring,omitempty"`
}

type SeriesDetail struct {
	Series SeriesDetailItem `yaml:"series"`
}

type SeriesDetailItem struct {
	ID           int      `yaml:"id"`
	Title        string   `yaml:"title"`
	Year         int      `yaml:"year"`
	Status       string   `yaml:"status"`
	Network      string   `yaml:"network"`
	Seasons      int      `yaml:"seasons"`
	Episodes     string   `yaml:"episodes"`
	Size         string   `yaml:"size"`
	Path         string   `yaml:"path"`
	Genres       []string `yaml:"genres,omitempty"`
	NextAiring   string   `yaml:"nextAiring,omitempty"`
	Overview     string   `yaml:"overview,omitempty"`
}

type CalendarList struct {
	Episodes []CalendarItem `yaml:"episodes"`
}

type CalendarItem struct {
	Series        string `yaml:"series"`
	Episode       string `yaml:"episode"`
	Title         string `yaml:"title"`
	AirDate       string `yaml:"airDate"`
	HasFile       bool   `yaml:"hasFile"`
}

type QueueList struct {
	Queue []QueueItem `yaml:"queue"`
	Total int         `yaml:"total"`
}

type QueueItem struct {
	Series    string `yaml:"series"`
	Episode   string `yaml:"episode"`
	Title     string `yaml:"title"`
	Status    string `yaml:"status"`
	Size      string `yaml:"size"`
	Remaining string `yaml:"remaining"`
	TimeLeft  string `yaml:"timeLeft,omitempty"`
}

type WantedList struct {
	Episodes []WantedItem `yaml:"episodes"`
	Total    int          `yaml:"total"`
}

type WantedItem struct {
	SeriesID int    `yaml:"seriesId"`
	Episode  string `yaml:"episode"`
	Title    string `yaml:"title"`
	AirDate  string `yaml:"airDate"`
}

type SearchResultList struct {
	Results []SearchResultItem `yaml:"results"`
}

type SearchResultItem struct {
	Title    string `yaml:"title"`
	Year     int    `yaml:"year"`
	TvdbID   int    `yaml:"tvdbId"`
	Network  string `yaml:"network"`
	Seasons  int    `yaml:"seasons"`
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
	return t.Local().Format("2006-01-02 15:04")
}

func formatEpisodeCode(season, episode int) string {
	return fmt.Sprintf("S%02dE%02d", season, episode)
}

// Conversion methods

func (s *APISeries) ToListItem() SeriesListItem {
	item := SeriesListItem{
		ID:       s.ID,
		Title:    s.Title,
		Year:     s.Year,
		Status:   s.Status,
		Network:  s.Network,
		Seasons:  s.SeasonCount,
		Episodes: fmt.Sprintf("%d/%d", s.EpisodeFileCount, s.EpisodeCount),
		Size:     formatBytes(s.SizeOnDisk),
	}
	if s.NextAiring != nil {
		item.NextAiring = formatTime(s.NextAiring)
	}
	return item
}

func (s *APISeries) ToDetail() SeriesDetailItem {
	item := SeriesDetailItem{
		ID:       s.ID,
		Title:    s.Title,
		Year:     s.Year,
		Status:   s.Status,
		Network:  s.Network,
		Seasons:  s.SeasonCount,
		Episodes: fmt.Sprintf("%d/%d", s.EpisodeFileCount, s.EpisodeCount),
		Size:     formatBytes(s.SizeOnDisk),
		Path:     s.Path,
	}
	if len(s.Genres) > 0 {
		item.Genres = s.Genres
	}
	if s.NextAiring != nil {
		item.NextAiring = formatTime(s.NextAiring)
	}
	if s.Overview != "" {
		overview := s.Overview
		if len(overview) > 500 {
			overview = overview[:497] + "..."
		}
		item.Overview = overview
	}
	return item
}

func (c *APICalendarEntry) ToListItem() CalendarItem {
	return CalendarItem{
		Series:  c.SeriesTitle,
		Episode: formatEpisodeCode(c.SeasonNumber, c.EpisodeNumber),
		Title:   c.Title,
		AirDate: formatTime(&c.AirDateUtc),
		HasFile: c.HasFile,
	}
}

func (q *APIQueueItem) ToListItem() QueueItem {
	item := QueueItem{
		Series:    q.Series.Title,
		Episode:   formatEpisodeCode(q.Episode.SeasonNumber, q.Episode.EpisodeNumber),
		Title:     q.Episode.Title,
		Status:    q.Status,
		Size:      formatBytes(q.Size),
		Remaining: formatBytes(q.Sizeleft),
	}
	if q.Timeleft != "" {
		item.TimeLeft = q.Timeleft
	}
	return item
}

func (e *APIEpisode) ToWantedItem() WantedItem {
	return WantedItem{
		SeriesID: e.SeriesID,
		Episode:  formatEpisodeCode(e.SeasonNumber, e.EpisodeNumber),
		Title:    e.Title,
		AirDate:  e.AirDate,
	}
}

func (r *APISearchResult) ToListItem() SearchResultItem {
	return SearchResultItem{
		Title:   r.Title,
		Year:    r.Year,
		TvdbID:  r.TvdbID,
		Network: r.Network,
		Seasons: r.SeasonCount,
	}
}
