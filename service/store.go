package service

import (
	"smply/config"
	"smply/model"
	"smply/utils"
)

func StoreUrl(url string, short string) (model.Url, error) {
	saved, _ := GetByOriginal(url)

	if (saved != model.Url{}) {
		return saved, nil
	}

	tx, err := config.DB.Begin(pgCtx)
	if err != nil {
		return model.Url{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(pgCtx)
		}
	}()

	var id int64
	err = tx.QueryRow(
		pgCtx,
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
		pgCtx,
		`UPDATE urls SET short = $1 WHERE id = $2`,
		short,
		id,
	)
	if err != nil {
		return model.Url{}, err
	}

	if err = tx.Commit(pgCtx); err != nil {
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
