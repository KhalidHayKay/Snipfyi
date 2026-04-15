package service

import (
	"context"
	"smply/internal/storage"
	"smply/model"
)

func GetByShort(ctx context.Context, short string) (model.Url, error) {
	var url model.Url

	err := storage.DB.QueryRow(
		ctx,
		`SELECT id, original, short FROM urls WHERE short = $1`,
		short).Scan(
		&url.Id,
		&url.Original,
		&url.Short,
	)

	if err != nil {
		return model.Url{}, err
	}

	url.BuildUrls()
	return url, nil
}

func GetByOriginal(ctx context.Context, originalUrl string) (model.Url, error) {
	var url model.Url

	err := storage.DB.QueryRow(
		ctx,
		`SELECT id, original, short FROM urls WHERE original = $1`,
		originalUrl).Scan(
		&url.Id,
		&url.Original,
		&url.Short,
	)

	if err != nil {
		return model.Url{}, err
	}

	url.BuildUrls()
	return url, nil
}
