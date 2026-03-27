package handlers

import (
	"errors"
	"log"
	"net/http"
	"smply/service"
	"smply/utils"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	alias := r.FormValue("alias")

	if url == "" {
		Error(w, http.StatusUnprocessableEntity, "'url' is a required field")
		return
	}

	if !utils.IsValidURL(url) {
		Error(w, http.StatusUnprocessableEntity, "Url not valid")
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
			Error(w, http.StatusConflict, msg)
			return
		}

		Error(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	Success(w, http.StatusCreated, result)
}

func StatsApi(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	stat, err := service.GetStats(r.Context(), code)

	if err != nil {
		log.Println(err)
		Error(w, http.StatusNotFound, "Not found")
		return
	}

	Success(w, http.StatusOK, stat)
}

func RedirectApi(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	stat, err := service.GetByShort(r.Context(), code)

	if err != nil {
		log.Println(err)
		Error(w, http.StatusNotFound, "Not found")
		return
	}

	Success(w, http.StatusOK, stat)
}
