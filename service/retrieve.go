package service

import (
	"smply/config"
	"smply/model"
)

func GetByShort(short string) (model.Url, error) {
	var url model.Url

	err := config.DB.QueryRow(
		`SELECT id, original, short FROM urls WHERE short = ?`,
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
		`SELECT id, original, short FROM urls WHERE original = ?`,
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
