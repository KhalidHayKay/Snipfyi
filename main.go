package main

import (
	"log"
	"net/http"
	"snipfyi/config"
	"snipfyi/handlers"
)

func main() {

	err := config.LoadEnv()
	if err != nil {
		log.Println(err)
	}

	err = config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/shorten", handlers.ShortenPage)

	http.HandleFunc("POST /api/shorten", handlers.Shorten)
	http.HandleFunc("/r/{code}", handlers.Redirect)
	http.HandleFunc("/r/{code}/stats", handlers.Stats)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
