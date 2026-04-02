package handler

import (
	"context"
	"log"
	"net/http"
	"smply/internal/render"
	"smply/model"
	"smply/service"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "home.html", render.ViewData{
		Title: "Home",
		Page:  "home",
	})
}

func ShortenPage(w http.ResponseWriter, r *http.Request) {
	data := render.ViewData{
		Title: "Shorten URL",
		Page:  "shorten",
	}

	render.Page(w, "shorten.html", data)
}

func ApiPage(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "api.html", render.ViewData{
		Title: "API",
		Page:  "api",
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

	render.Page(w, "stats.html", render.ViewData{
		Title: "URL Stats",
		Page:  "stats",
		Data: map[string]model.Url{
			"Stats": stat,
		},
	})
}
