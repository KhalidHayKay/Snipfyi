package service

import (
	"context"
	"smply/internal/storage"
	"smply/model"
)

func GetByAlias(ctx context.Context, alias string) (model.Url, error) {
	var url model.Url

	err := storage.DB.QueryRow(
		ctx,
		`SELECT id, original, alias FROM urls WHERE alias = $1`,
		alias).Scan(
		&url.Id,
		&url.Original,
		&url.Alias,
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
		`SELECT id, original, alias FROM urls WHERE original = $1`,
		originalUrl).Scan(
		&url.Id,
		&url.Original,
		&url.Alias,
	)

	if err != nil {
		return model.Url{}, err
	}

	url.BuildUrls()
	return url, nil
}
