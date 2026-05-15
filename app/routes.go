package app

import (
	"net/http"
	"smply/app/middleware"
	"smply/internal/home"
	"smply/internal/limiter"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func setupRouter(handlers Handlers, middleware *middleware.Middleware) *chi.Mux {
	router := chi.NewRouter()

	// Global middleware
	router.Use(middleware.CORS, chiMiddleware.Logger, chiMiddleware.Recoverer)

	// Rate limiters
	keyRequestRateLimiter := limiter.NewRateLimiter(2/60.0, 2)
	publicAPIRateLimiter := limiter.NewRateLimiter(10/60.0, 10)

	// Page routes
	router.Get("/", home.Page)
	router.Get("/shorten", handlers.URL.ShortenPage)

	// TODO: rate limit
	router.Post("/shorten", handlers.URL.HandleShortenForm)

	router.Get("/api", handlers.APIKey.Page)
	router.With(keyRequestRateLimiter.Middleware).Post("/api", handlers.APIKey.Request)

	router.Get("/stats/{alias}", handlers.Stat.Page)

	router.Get("/key/activate", handlers.APIKey.Create)

	router.Route("/admin", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminGuest)
			r.Get("/login", handlers.Admin.Page)
		})
		r.Post("/login", handlers.Admin.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthenticateAdmin)
			r.Get("/stats", handlers.Stat.GetForAdmin)
		})

		r.Get("/auth/redirect", handlers.Admin.Auth)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/admin/stats", http.StatusFound)
		})
	})

	router.Get("/{alias}", handlers.URL.Redirect)
	// Note: Specific stats routes should be defined before the catch-all redirect route to avoid conflicts

	// Public API routes with rate limiting and key validation
	router.Route("/api/v1", func(r chi.Router) {
		r.Use(publicAPIRateLimiter.Middleware)
		r.Use(middleware.RequireKey)

		r.Post("/shorten", handlers.URL.Shorten)
		r.Get("/stats/{alias}", handlers.Stat.Get)
		r.Get("/show/{alias}", handlers.URL.Get)
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
