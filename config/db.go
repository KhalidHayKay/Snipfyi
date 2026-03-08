package config

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() error {

	db, err := sql.Open("sqlite", Env.DbUrl)
	if err != nil {
		return err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return pingErr
	}

	db.SetMaxOpenConns(1)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			original TEXT UNIQUE,
			short TEXT UNIQUE,
			visited INTEGER DEFAULT 0,
			created DATETIME
		);
	`)
	if err != nil {
		return err
	}

	DB = db

	return nil
}
