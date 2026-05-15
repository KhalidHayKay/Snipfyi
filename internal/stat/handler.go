package stat

import (
	"net/http"
	"smply/app/render"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Page(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	stats, err := h.service.Get(r.Context(), alias)
	if err != nil {
		render.ErrorPage(w, http.StatusNotFound, "Not found")
		return
	}

	render.Page(w, "stats.html", render.ViewData{
		Title: "URL Stats",
		Page:  "stats",
		Data: map[string]any{
			"Stats": stats,
		},
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	stats, err := h.service.Get(r.Context(), alias)
	if err != nil {
		render.ErrorJSON(w, http.StatusNotFound, "Not found")
		return
	}

	render.JSON(w, http.StatusOK, stats)
}

func (h *Handler) GetForAdmin(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetAdmin(r.Context())
	if err != nil {
		render.ErrorPage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	render.AdminPage(w, "admin/stats.html", render.ViewData{
		Title: "Admin Stats",
		Page:  "admin_stats",
		Data: map[string]any{
			"Stats": stats,
		},
	})
}
