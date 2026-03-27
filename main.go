package main

import (
	"log"
	"net/http"
	"smply/config"
	"smply/handlers"
)

func main() {
	config.LoadEnv()

	err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	files := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", files))

	// Page routes
	mux.HandleFunc("GET /", handlers.Home)
	mux.HandleFunc("GET /shorten", handlers.ShortenPage)
	mux.HandleFunc("GET /api", handlers.ApiPage)
	mux.HandleFunc("GET /stats/{code}", handlers.Stats)
	mux.HandleFunc("GET /{code}", handlers.Redirect)

	// Private API routes
	mux.HandleFunc("POST /api/key/request", handlers.RequestApiKey)
	mux.HandleFunc("GET /key/activate", handlers.CreateApiKey)

	//Public API routes
	mux.HandleFunc("POST /api/v1/shorten", handlers.Shorten)
	mux.HandleFunc("GET /api/v1/stats/{code}", handlers.StatsApi)
	mux.HandleFunc("GET /api/v1/redirect/{code}", handlers.RedirectApi)

	port := config.Env.AppPort

	if port == "" && config.Env.AppEnv == "development" {
		port = "8000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s in %s mode", port, config.Env.AppEnv)
	log.Fatal(server.ListenAndServe())
}
