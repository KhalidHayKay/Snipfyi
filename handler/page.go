package handler

import (
	"log"
	"net/http"
	"smply/internal/render"
	"smply/internal/service"
	"smply/model"

	"github.com/go-chi/chi/v5"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "home.html", render.ViewData{
		Title: "Home",
		Page:  "home",
	})
}

func ShortenPage(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "shorten.html", render.ViewData{
		Title: "Shorten URL",
		Page:  "shorten",
	})
}

func ApiPage(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "api.html", render.ViewData{
		Title: "API",
		Page:  "api",
	})
}

func StatsPage(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	stat, err := service.GetStats(r.Context(), alias)
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
