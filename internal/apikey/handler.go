package apikey

import (
	"html/template"
	"net/http"
	"smply/app/render"

	"github.com/jackc/pgx/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Page(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "api.html", render.ViewData{
		Title: "API",
		Page:  "api",
	})
}

func (h *Handler) Request(w http.ResponseWriter, r *http.Request) {
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

	err := h.service.RequestNew(r.Context(), email)
	if err != nil {
		data.Error = "Unable to process your request. Please try again later."
		render.Page(w, "api.html", data)
		return
	}

	data.Data = map[string]string{"SentTo": email}
	render.Page(w, "api.html", data)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		render.ErrorPage(w, http.StatusUnprocessableEntity, "Cannot find token")
		return
	}

	key, err := h.service.Create(r.Context(), token)
	if err != nil {
		if err == pgx.ErrNoRows {
			render.ErrorPage(w, http.StatusUnprocessableEntity, "Invalid or expired token")
			return
		}

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
