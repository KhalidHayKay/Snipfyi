package config

import (
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

	return nil
}
