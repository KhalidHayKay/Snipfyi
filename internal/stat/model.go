package stat

import (
	"fmt"
	"smply/config"
	"time"
)

type ClickEvent struct {
	Id        int64     `json:"id"`
	UrlId     int64     `json:"link_id"`
	Referer   string    `json:"referer"`
	UserAgent string    `json:"user_agent"`
	IpAddress string    `json:"ip_address"`
	Timestamp time.Time `json:"timestamp"`
}

type Stats struct {
	Visited     int64      `json:"visited"`
	LastVisited *time.Time `json:"last_visited"`

	ClickEvents []ClickEvent `json:"click_events"`

	// derived fields
	Original string     `json:"original"`
	Alias    string     `json:"alias"`
	Created  *time.Time `json:"created"`
	ShortUrl string     `json:"short_url"`
}

func (s *Stats) BuildShortUrl() {
	s.ShortUrl = fmt.Sprintf("%s/%s", config.Env.App.Url, s.Alias)
}

// Admin Stats

type TopLink struct {
	Alias  string
	Clicks int64
}

type DailyStat struct {
	Date   string
	Clicks int64
}

type TopReferer struct {
	Referer string
	Clicks  int64
}

type TopDevice struct {
	Device string
	Clicks int64
}

type AdminStats struct {
	// URLs
	TotalUrls int64
	UrlsToday int64

	// Redirects (sum of urls.visited)
	TotalRedirects int64

	// Click events
	TotalClicks        int64
	ClicksToday        int64
	ClicksThisWeek     int64
	UniqueLinksClicked int64

	// Peak activity
	PeakDay        string
	PeakHour       int
	PeakHourClicks int64

	// Click trend — last 7 days
	DailyTrend []DailyStat

	// Top links
	TopLinks []TopLink

	// Top referers
	TopReferers []TopReferer

	// Device breakdown
	TopDevices []TopDevice

	// API
	TotalApiKeys  int64
	ActiveApiKeys int64

	// Magic tokens
	TotalTokens int64
	TokensToday int64
}
