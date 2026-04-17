package handler

import (
	"log"
	"net/http"
	"smply/internal/queue"
	"smply/internal/render"
	"smply/internal/service"
	"text/template"

	"github.com/jackc/pgx/v5"
)

func RequestApiKey(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	data := render.ViewData{
		Title: "API",
		Page:  "api",
	}

	if email == "" {
		data.Error = "'email' is a required field"
		render.Page(w, "api.html", data)
		return
	}

	token, err := service.CreateMagicToken(r.Context(), email)
	if err != nil {
		log.Printf("Error creating magic token for %s: %v", email, err)
		data.Error = "Unable to process your request. Please try again later."
		render.Page(w, "api.html", data)
		return
	}

	err = queue.EnqueueAPIKeyMagicLinkEmail(r.Context(), email, token)
	if err != nil {
		log.Printf("failed to enqueue email: %v", err)
	}

	data.Data = map[string]string{"SentTo": email}
	render.Page(w, "api.html", data)
}

func CreateApiKey(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		render.ErrorPage(w, http.StatusUnprocessableEntity, "Cannot find token")
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
