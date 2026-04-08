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
	router.Use(middleware.CORSMiddleware, chiMiddleware.Logger, chiMiddleware.Recoverer)

	// Rate limiters
	keyRequestRateLimiter := limiter.NewRateLimiter(2/60.0, 2)
	publicAPIRateLimiter := limiter.NewRateLimiter(10/60.0, 10)

	// Page routes
	router.Get("/", handler.Home)
	router.Post("/", handler.HomeShorten)

	router.Get("/shorten", handler.ShortenPage)
	router.Post("/shorten", handler.ShortenPageShorten)

	router.Get("/stats/{code}", handler.StatsPage)

	router.Get("/api", handler.ApiPage)
	router.With(keyRequestRateLimiter.Middleware).Post("/api", handler.RequestApiKey)

	router.Get("/key/activate", handler.CreateApiKey)

	router.Get("/{code}", handler.ResolveRedirect)
	// Note: Specific stats routes should be defined before the catch-all redirect route to avoid conflicts

	// Public API routes with rate limiting and key validation
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(publicAPIRateLimiter.Middleware)
		r.Use(middleware.RequireKey)

		r.Post("/shorten", handler.Shorten)
		r.Get("/stats/{code}", handler.Stats)
		r.Get("/redirect/{code}", handler.Redirect)
	})

	// Static files
	publicFiles := http.FileServer(http.Dir("./public"))
	staticFiles := http.FileServer(http.Dir("./static"))

	router.Handle("/robots.txt", publicFiles)
	router.Handle("/sitemap.xml", publicFiles)
	router.Handle("/favicon.ico", publicFiles)
	router.Handle("/static/*", http.StripPrefix("/static/", staticFiles))

	return router
}
