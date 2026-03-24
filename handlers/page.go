package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"smply/service"
	"smply/views"
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

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	url, err := service.GetByShort(code)

	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "Not found")
			return
		}

		Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	go func() {
		err := service.RunStats(url.Id)

		if err != nil {
			log.Println(err)
		}
	}()

	http.Redirect(w, r, url.Original, http.StatusFound)
}

func Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	stat, err := service.GetStats(code)

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
