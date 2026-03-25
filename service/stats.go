package service

import (
	"context"
	"smply/config"
	"smply/model"
	"time"
)

func RunStats(ctx context.Context, id int64) error {
	_, err := config.DB.Exec(ctx, `
		UPDATE urls
		SET visited = visited + 1, last_visited = $1
		WHERE id = $2
	`, time.Now(), id)

	return err
}

func GetStats(ctx context.Context, short string) (model.Url, error) {
	var url model.Url
	err := config.DB.QueryRow(ctx,
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
