package service

import (
	"context"
	"smply/internal/storage"
	"smply/model"
	"smply/utils"
)

func StoreUrl(ctx context.Context, url string, alias string) (model.Url, error) {
	saved, _ := GetByOriginal(ctx, url)

	if (saved != model.Url{}) {
		return saved, nil
	}

	tx, err := storage.DB.Begin(ctx)
	if err != nil {
		return model.Url{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var id int64
	err = tx.QueryRow(
		ctx,
		`INSERT INTO urls (original, alias) VALUES ($1, $2) RETURNING id`,
		url,
		alias,
	).Scan(&id)
	if err != nil {
		return model.Url{}, err
	}

	if alias == "" {
		alias = utils.EncodeWithPadding(id, 2)
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE urls SET alias = $1 WHERE id = $2`,
		alias,
		id,
	)
	if err != nil {
		return model.Url{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return model.Url{}, err
	}

	result := model.Url{
		Id:       id,
		Original: url,
		Alias:    alias,
	}

	result.BuildUrls()
	return result, nil
}

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
