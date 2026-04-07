package app

import (
	"log"
	"net/http"
	"smply/config"
)

func Start() {
	config.LoadEnv()

	if err := config.InitDB(); err != nil {
		log.Fatal(err)
	}

	router := setupRouter()

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":"+config.Env.App.Port, router))
}
