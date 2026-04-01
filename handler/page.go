package handler

import (
	"context"
	"log"
	"net/http"
	"smply/internal/render"
	"smply/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Home",
		"Page":  "home",
	}

	render.Page(w, "home.html", data)
}

func ShortenPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Shorten URL",
		"Page":  "shorten",
	}

	render.Page(w, "shorten.html", data)
}

func ApiPage(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "api.html", map[string]any{
		"Title": "API",
		"Page":  "api",
	})
}

func ResolveRedirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := service.GetByShort(r.Context(), code)

	if err != nil {
		if err == pgx.ErrNoRows {
			render.ErrorPage(w, http.StatusNotFound, "Not found")
			return
		}

		render.ErrorPage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := service.RunStats(ctx, url.Id)

		if err != nil {
			log.Println(err)
		}
	}()

	http.Redirect(w, r, url.Original, http.StatusFound)
}

func StatsPage(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	stat, err := service.GetStats(r.Context(), code)
	if err != nil {
		log.Println(err)
		render.ErrorPage(w, http.StatusNotFound, "Not found")
		return
	}

	data := map[string]any{
		"Title": "URL Stats",
		"Page":  "stats",
		"Stats": stat,
	}

	render.Page(w, "stats.html", data)
}
