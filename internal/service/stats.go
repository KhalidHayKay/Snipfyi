package service

import (
	"context"
	"log"
	"smply/internal/storage"
	"smply/model"
	"time"
)

func RunStats(ctx context.Context, alias, referer, userAgent, ipAddress string, timestamp time.Time) error {
	tx, err := storage.DB.Begin(ctx)
	if err != nil {
		log.Printf("Failed to begin transaction for stats update: %v", err)
		return err
	}
	defer tx.Rollback(ctx)

	var id int
	err = tx.QueryRow(ctx, `
		UPDATE urls
		SET visited = visited + 1, last_visited = $1
		WHERE alias = $2
		RETURNING id
	`, timestamp, alias).Scan(&id)
	if err != nil {
		log.Printf("Failed to update stats for alias %s: %v", alias, err)
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO click_events (url_id, referer, user_agent, ip_address, timestamp)
		VALUES ($1, $2, $3, $4, $5)
	`, id, referer, userAgent, ipAddress, timestamp)
	if err != nil {
		log.Printf("Failed to insert click event for URL ID %d: %v", id, err)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction for stats update: %v", err)
		return err
	}

	log.Printf("Stats update for URL ID: %d successful", id)
	return nil
}

func GetStats(ctx context.Context, alias string) (model.Url, error) {
	var url model.Url
	err := storage.DB.QueryRow(ctx,
		`SELECT id, original, alias, visited, created, last_visited FROM urls WHERE alias = $1`,
		alias).Scan(
		&url.Id,
		&url.Original,
		&url.Alias,
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
