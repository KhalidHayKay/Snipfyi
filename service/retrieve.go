package service

import (
	"context"
	"smply/config"
	"smply/model"
)

var pgCtx = context.Background()

func GetByShort(short string) (model.Url, error) {
	var url model.Url

	err := config.DB.QueryRow(
		pgCtx,
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

func GetByOriginal(originalUrl string) (model.Url, error) {
	var url model.Url

	err := config.DB.QueryRow(
		pgCtx,
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
