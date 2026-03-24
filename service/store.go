package service

import (
	"smply/config"
	"smply/model"
	"smply/utils"
	"time"
)

func StoreUrl(url string, short string) (model.Url, error) {
	saved, _ := GetByOriginal(url)

	if (saved != model.Url{}) {
		return saved, nil
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return model.Url{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	res, err := tx.Exec(
		`INSERT INTO urls (original, short) VALUES (?, ?)`,
		url,
		time.Now().Unix(),
	)
	if err != nil {
		return model.Url{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Url{}, err
	}

	if short == "" {
		short = utils.EncodeWithPadding(id)
	}

	_, err = tx.Exec(`
		UPDATE urls
		SET short = ?
		WHERE id = ?
	`, short, id)
	if err != nil {
		return model.Url{}, err
	}

	if err = tx.Commit(); err != nil {
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
