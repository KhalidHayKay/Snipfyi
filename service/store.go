package service

import (
	"context"
	"smply/config"
	"smply/model"
	"smply/utils"
)

func StoreUrl(ctx context.Context, url string, short string) (model.Url, error) {
	saved, _ := GetByOriginal(ctx, url)

	if (saved != model.Url{}) {
		return saved, nil
	}

	tx, err := config.DB.Begin(ctx)
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
		`INSERT INTO urls (original, short) VALUES ($1, $2) RETURNING id`,
		url,
		short,
	).Scan(&id)
	if err != nil {
		return model.Url{}, err
	}

	if short == "" {
		short = utils.EncodeWithPadding(id)
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE urls SET short = $1 WHERE id = $2`,
		short,
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
		Short:    short,
	}

	result.BuildUrls()
	return result, nil
}
