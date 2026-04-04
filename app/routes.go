package app

import (
	"net/http"
	"smply/handler"
	"smply/internal/limiter"
	"smply/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	router.Use(chiMiddleware.Logger, chiMiddleware.Recoverer)

	// Static files
	files := http.FileServer(http.Dir("./public"))
	router.Handle("/static/*", http.StripPrefix("/static/", files))

	// Page routes
	router.Get("/", handler.Home)
	router.Get("/shorten", handler.ShortenPage)
	router.Get("/api", handler.ApiPage)
	router.Get("/key/activate", handler.CreateApiKey)
	router.Get("/stats/{code}", handler.StatsPage)
	router.Get("/{code}", handler.ResolveRedirect)
	// Note: Specific stats routes should be defined before the catch-all redirect route to avoid conflicts

	// Private API routes (TODO: protect these routes)
	router.Post("/api/internal/shorten", handler.Shorten)
	router.Route("/api/internal/key/request", func(r chi.Router) {
		rateLimiter := limiter.NewRateLimiter(1/60.0, 1)
		r.Use(rateLimiter.Middleware)

		r.Post("/", handler.RequestApiKey)
	})

	// Public API routes with rate limiting and key validation
	router.Route("/api/v1", func(r chi.Router) {
		rateLimiter := limiter.NewRateLimiter(10/60.0, 10)
		r.Use(rateLimiter.Middleware)
		r.Use(middleware.RequireKey)

		r.Post("/shorten", handler.Shorten)
		r.Get("/stats/{code}", handler.Stats)
		r.Get("/redirect/{code}", handler.Redirect)
	})

	return router
}
