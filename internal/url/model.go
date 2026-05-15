package url

import (
	"fmt"
	"smply/config"
	"time"
)

type Url struct {
	Id       int64      `json:"id"`
	Original string     `json:"original"`
	Alias    string     `json:"alias"`
	Created  *time.Time `json:"created"`

	// derived fields
	ShortUrl string `json:"short_url"`
	StatUrl  string `json:"stat_url"`
}

func (u *Url) BuildUrls() {
	u.ShortUrl = fmt.Sprintf("%s/%s", config.Env.App.Url, u.Alias)
	u.StatUrl = fmt.Sprintf("%s/stats/%s", config.Env.App.Url, u.Alias)
}
