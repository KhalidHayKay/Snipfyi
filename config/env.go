package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type EnvType struct {
	AppUrl string
	DbUrl  string
}

var Env *EnvType

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	Env = &EnvType{
		AppUrl: os.Getenv("APP_URL"),
		DbUrl:  os.Getenv("DB_URL"),
	}

	if Env.AppUrl == "" || Env.DbUrl == "" {
		return errors.New("APP_URL or DB_URL not set")

	}

	return nil
}
