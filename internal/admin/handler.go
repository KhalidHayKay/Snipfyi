package admin

import (
	"net/http"
	"smply/app/render"
	"smply/config"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Page(w http.ResponseWriter, r *http.Request) {
	render.AdminPage(w, "admin/login.html", render.ViewData{})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	if email == "" {
		render.AdminPage(w, "admin/login.html", render.ViewData{
			Error: "Email is required",
		})
		return
	}

	err := h.service.AttemptMagicLinkLogin(r.Context(), email)
	if err != nil {
		render.AdminPage(w, "admin/login.html", render.ViewData{
			Error: err.Error(),
		})
		return
	}

	render.AdminPage(w, "admin/login.html", render.ViewData{
		Data: map[string]string{
			"SentTo": email,
		},
	})
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		render.ErrorPage(w, http.StatusUnprocessableEntity, "Cannot find token")
		return
	}

	sessionId, err := h.service.Authenticate(r.Context(), token)
	if err != nil {
		render.AdminPage(w, "admin/login.html", render.ViewData{
			Error: err.Error(),
		})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     strings.ToLower(config.Env.App.Name) + "_session_id",
		Value:    sessionId,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/admin/stats", http.StatusFound)
}
