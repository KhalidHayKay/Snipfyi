package service

import (
	"shortener/config"
)

func Retrieve(short string) (string, error) {
	var url string
	err := config.DB.QueryRow(`
		SELECT original FROM urls WHERE urls.short = ?
	`, short).Scan(&url)

	if err != nil {
		return "", err
	}

	return url, nil
}
