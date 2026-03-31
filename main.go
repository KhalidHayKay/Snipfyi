package main

import (
	"log"
	"net/http"
	"smply/config"
	"smply/handler"
	"smply/middleware"
)

func apply(h http.HandlerFunc, m ...func(http.Handler) http.Handler) http.Handler {
	var res http.Handler = h
	for _, middleware := range m {
		res = middleware(res)
	}
	return res
}

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
	mux.Handle("GET /", http.HandlerFunc(handler.Home))
	mux.Handle("GET /shorten", http.HandlerFunc(handler.ShortenPage))
	mux.Handle("GET /api", http.HandlerFunc(handler.ApiPage))
	mux.Handle("GET /stats/{code}", http.HandlerFunc(handler.Stats))
	mux.Handle("GET /{code}", http.HandlerFunc(handler.Redirect))
	mux.Handle("GET /key/activate", http.HandlerFunc(handler.CreateApiKey))

	// Private API routes
	//#TODO - Need to protect these routes
	mux.Handle("POST /api/internal/shorten", http.HandlerFunc(handler.Shorten))
	mux.Handle("POST /api/internal/key/request", http.HandlerFunc(handler.RequestApiKey))

	//Public API routes
	mux.Handle("POST /api/v1/shorten", middleware.Apply(handler.Shorten, middleware.RequireKey))
	mux.Handle("GET /api/v1/stats/{code}", middleware.Apply(handler.StatsApi, middleware.RequireKey))
	mux.Handle("GET /api/v1/redirect/{code}", middleware.Apply(handler.RedirectApi, middleware.RequireKey))

	port := config.Env.App.Port

	if port == "" && config.Env.App.Environment == "development" {
		port = "8000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s in %s mode", port, config.Env.App.Environment)
	log.Fatal(server.ListenAndServe())
}
