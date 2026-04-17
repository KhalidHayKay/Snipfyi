package handler

import (
	"errors"
	"log"
	"net/http"
	"smply/internal/render"
	"smply/internal/service"
	"smply/utils"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
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

	result, err := service.StoreUrl(r.Context(), url, alias)
	if err != nil {
		log.Println(err)

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

	render.JSON(w, http.StatusCreated, result)
}

func Stats(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	stat, err := service.GetStats(r.Context(), alias)

	if err != nil {
		log.Println(err)
		render.ErrorJSON(w, http.StatusNotFound, "Not found")
		return
	}

	render.JSON(w, http.StatusOK, stat)
}

func RedirectAPI(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	url, err := service.GetByAlias(r.Context(), alias)

	if err != nil {
		log.Println(err)
		render.ErrorJSON(w, http.StatusNotFound, "Not found")
		return
	}

	render.JSON(w, http.StatusOK, url)
}
