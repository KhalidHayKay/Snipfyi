package main

import (
	"log"
	"net/http"
	"shortener/config"
	"shortener/handlers"
)

func main() {

	envErr := config.LoadEnv()
	if envErr != nil {
		log.Fatal(envErr)
	}

	dbErr := config.InitDB()
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	http.HandleFunc("POST /shorten", handlers.Shorten)
	http.HandleFunc("/{redirect}", handlers.Redirect)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
