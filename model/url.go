package model

import (
	"fmt"
	"smply/config"
	"time"
)

type Url struct {
	Id          int64      `json:"id"`
	Original    string     `json:"original"`
	Short       string     `json:"short"`
	Stat        string     `json:"stat"`
	Visited     int64      `json:"visited"`
	Created     time.Time  `json:"created"`
	LastVisited *time.Time `json:"last_visited"`
}

func (u *Url) BuildUrls() {
	code := u.Short // save code before overwriting
	u.Short = fmt.Sprintf("%s/%s", config.Env.AppUrl, code)
	u.Stat = fmt.Sprintf("%s/stats/%s", config.Env.AppUrl, code)
}

type ClickEvent struct {
	Id        int64     `json:"id"`
	LinkId    int64     `json:"link_id"`
	Timestamp time.Time `json:"timestamp"`
	Referrer  string    `json:"referrer"`
	UserAgent string    `json:"user_agent"`
}
