package model

import (
	"fmt"
	"smply/config"
	"time"
)

type Url struct {
	Id          int64      `json:"id"`
	Original    string     `json:"original"`
	Alias       string     `json:"alias"`
	ShortUrl    string     `json:"short_url"`
	StatUrl     string     `json:"stat_url"`
	Visited     int64      `json:"visited"`
	Created     time.Time  `json:"created"`
	LastVisited *time.Time `json:"last_visited"`
}

func (u *Url) BuildUrls() {
	u.ShortUrl = fmt.Sprintf("%s/%s", config.Env.App.Url, u.Alias)
	u.StatUrl = fmt.Sprintf("%s/stats/%s", config.Env.App.Url, u.Alias)
}

type ClickEvent struct {
	Id        int64     `json:"id"`
	UrlId     int64     `json:"link_id"`
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"user_agent"`
	IpAddress string    `json:"ip_address"`
	Timestamp time.Time `json:"timestamp"`
}
