package home

import (
	"net/http"
	"smply/app/render"
)

func Page(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "home.html", render.ViewData{
		Title: "Home",
		Page:  "home",
	})
}
