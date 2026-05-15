package url

import (
	"errors"
	"log"
	"net/http"
	"smply/app/render"
	"smply/internal/queue"
	"smply/utils"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ShortenPage(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "shorten.html", render.ViewData{
		Title: "Shorten URL",
		Page:  "shorten",
	})
}

func (h *Handler) HandleShortenForm(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	alias := r.FormValue("alias")

	file := "shorten.html"

	data := render.ViewData{
		Title: "Shorten URL",
		Page:  "shorten",
	}

	if url == "" {
		data.Error = "'url' is a required field"
		render.Page(w, file, data)
		return
	}

	if !utils.IsValidURL(url) {
		data.Error = "URL is not valid"
		render.Page(w, file, data)
		return
	}

	result, err := h.service.Store(r.Context(), url, alias)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			msg := "A record with this value already exists"
			if pgErr.ConstraintName == "urls_short_key" || strings.Contains(pgErr.Message, "urls.short") {
				msg = "This alias is already taken"
			}
			data.Error = msg
			render.Page(w, file, data)
			return
		}

		data.Error = "Internal server error"
		render.Page(w, file, data)
		return
	}

	data.Data = result
	render.Page(w, file, data)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	alias := r.FormValue("alias")

	if url == "" {
		render.ErrorJSON(w, http.StatusUnprocessableEntity, "'url' is a required field")
		return
	}

	if !utils.IsValidURL(url) {
		render.ErrorJSON(w, http.StatusUnprocessableEntity, "Url not valid")
		return
	}

	shortUrl, err := h.service.Store(r.Context(), url, alias)
	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			msg := "A record with this value already exists"
			if pgErr.ConstraintName == "urls_short_key" || strings.Contains(pgErr.Message, "urls.short") {
				msg = "This alias is already taken"
			}
			render.ErrorJSON(w, http.StatusConflict, msg)
			return
		}

		render.ErrorJSON(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	render.JSON(w, http.StatusCreated, shortUrl)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	url, err := h.service.GetByAlias(r.Context(), alias)

	if err != nil {
		render.ErrorJSON(w, http.StatusNotFound, "Not found")
		return
	}

	render.JSON(w, http.StatusOK, url)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	url, err := h.service.GetByAlias(r.Context(), alias)
	if err != nil {
		if err == pgx.ErrNoRows {
			render.ErrorPage(w, http.StatusNotFound, "Not found")
			return
		}

		render.ErrorPage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// stats update called here because `GetByAlias` might be called without redirect happening.
	// For example an API call.
	err = queue.EnqueueStatsUpdate(
		r.Context(),
		url.Alias,
		r.Referer(),
		r.UserAgent(),
		time.Now(),
	)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, url.Original, http.StatusFound)
}
