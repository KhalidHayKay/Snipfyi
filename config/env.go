package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Environment string
	Port        string
	Url         string
}

type MailerConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

type RedisConfig struct {
	Url      string
	Password string
}

type EnvType struct {
	App AppConfig

	InternalApiKey string

	DbUrl string

	Mailer MailerConfig

	Redis RedisConfig
}

var Env *EnvType

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	mailerPort, err := strconv.Atoi(os.Getenv("MAILER_PORT"))
	if err != nil {
		// handle error
	}

	Env = &EnvType{
		App: AppConfig{
			Environment: os.Getenv("APP_ENV"),
			Port:        os.Getenv("APP_PORT"),
			Url:         os.Getenv("APP_URL"),
		},

		InternalApiKey: os.Getenv("INTERNAL_API_KEY"),

		DbUrl: os.Getenv("DB_URL"),

		Mailer: MailerConfig{
			Host: os.Getenv("MAILER_HOST"),
			Port: mailerPort,
			User: os.Getenv("MAILER_USER"),
			Pass: os.Getenv("MAILER_PASS"),
			From: os.Getenv("MAILER_FROM"),
		},

		Redis: RedisConfig{
			Url:      os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}
	if Env.App.Url == "" || Env.DbUrl == "" {
		log.Println(errors.New("APP_URL or DB_URL not set"))
	}
}
