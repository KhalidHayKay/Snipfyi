package app

import (
	"log"
	"net/http"
	"smply/config"
)

func Start() {
	router := setupRouter()

	port := config.Env.App.Port
	if port == "" && config.Env.App.Environment == "development" {
		port = "8000"
	}

	log.Printf("Starting server on port %s in %s mode", port, config.Env.App.Environment)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
