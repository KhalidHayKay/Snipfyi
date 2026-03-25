package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"smply/service"
	"smply/views"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title": "Home",
		"Page":  "home",
	}

	views.Render(w, "home.html", data)
}

func ShortenPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title": "Shorten URL",
		"Page":  "shorten",
	}

	views.Render(w, "shorten.html", data)
}

func ApiPage(w http.ResponseWriter, r *http.Request) {
	views.Render(w, "api.html", map[string]any{
		"Title": "API",
		"Page":  "api",
	})
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	url, err := service.GetByShort(r.Context(), code)

	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "Not found")
			return
		}

		Error(w, http.StatusInternalServerError, "Internal server error")
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

func Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	stat, err := service.GetStats(r.Context(), code)

	if err != nil {
		log.Println(err)
		Error(w, http.StatusNotFound, "Not found")
		return
	}

	data := map[string]any{
		"Title": "URL Stats",
		"Page":  "stats",
		"Stats": stat,
	}

	views.Render(w, "stats.html", data)
}
