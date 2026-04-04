package handler

import (
	"log"
	"net/http"
	"smply/internal/render"
	"smply/service"
	"text/template"

	"github.com/jackc/pgx/v5"
)

func RequestApiKey(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	if email == "" {
		render.ErrorJSON(w, http.StatusUnprocessableEntity, "'email' is a required field")
		return
	}

	err := service.RequestApiKey(r.Context(), email)
	if err != nil {
		log.Println(err)
		render.ErrorJSON(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	render.JSON(w, http.StatusOK, nil)
}

func CreateApiKey(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		render.ErrorPage(w, http.StatusUnprocessableEntity, "cannot find token")
		return
	}

	key, err := service.CreateApiKey(r.Context(), token)
	if err != nil {
		if err == pgx.ErrNoRows {
			render.ErrorPage(w, http.StatusUnprocessableEntity, "Invalid or expired token")
			return
		}

		log.Println(err)
		render.ErrorPage(w, http.StatusUnprocessableEntity, "Unable to create API key")
		return
	}

	data := map[string]any{
		"Title":  "Your API Key",
		"Page":   "api-key",
		"ApiKey": key,
	}

	// render as a standalone file, not through the layout
	t := template.Must(template.ParseFiles("templates/pages/api-key.html"))
	t.Execute(w, data)
}
