package service

import (
	"smply/config"
	"smply/model"
	"time"
)

func RunStats(id int64) error {
	_, err := config.DB.Exec(pgCtx, `
		UPDATE urls
		SET visited = visited + 1, last_visited = $1
		WHERE id = $2
	`, time.Now(), id)

	return err
}

func GetStats(short string) (model.Url, error) {
	var url model.Url
	err := config.DB.QueryRow(pgCtx,
		`SELECT id, original, short, visited, created, last_visited FROM urls WHERE short = $1`,
		short).Scan(
		&url.Id,
		&url.Original,
		&url.Short,
		&url.Visited,
		&url.Created,
		&url.LastVisited,
	)

	if err != nil {
		return model.Url{}, err
	}

	url.BuildUrls()
	return url, nil
}
