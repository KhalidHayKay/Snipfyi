package service

import (
	"fmt"
	"shortener/config"
	"shortener/model"
	"time"
)

func StoreUrl(url string) (model.Url, error) {
	short := "Qr12bwg"
	now := time.Now()

	res, err := config.DB.Exec(
		`INSERT INTO urls (original, short, created) VALUES (?, ?, ?)`,
		url,
		short,
		now,
	)
	if err != nil {
		return model.Url{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Url{}, err
	}

	result := model.Url{
		Id:       id,
		Original: url,
		Short:    fmt.Sprintf("%s/%s", config.Env.AppUrl, short),
		Visited:  0,
		Created:  now,
	}

	return result, nil
}
