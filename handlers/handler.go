package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"snipfyi/service"
	"snipfyi/utils"
	"snipfyi/views"
)

func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
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

func Stats(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	stat, err := service.Retrieve(code)

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

func Shorten(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("url")

	if input == "" {
		Error(w, http.StatusUnprocessableEntity, "'url' is a required field")
		return
	}

	if !utils.IsValidURL(input) {
		Error(w, http.StatusUnprocessableEntity, "Url not valid")
		return
	}

	result, err := service.StoreUrl(input)
	if err != nil {
		Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	result.ShortToUrl()

	Success(w, http.StatusCreated, result)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	url, err := service.Retrieve(code)

	if err != nil {
		if err == sql.ErrNoRows {
			Error(w, http.StatusNotFound, "Not found")
			return
		}

		Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	go func() {
		err := service.IncrementVisited(url.Id)

		if err != nil {
			log.Println(err)
		}
	}()

	http.Redirect(w, r, url.Original, http.StatusFound)
}
