package handlers

import (
	"errors"
	"log"
	"net/http"
	"smply/service"
	"smply/utils"
	"strings"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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

	result, err := service.StoreUrl(url, alias)
	if err != nil {
		log.Println(err)

		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			msg := "A record with this value already exists"
			if strings.Contains(err.Error(), "urls.short") {
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

	stat, err := service.GetStats(code)

	if err != nil {
		log.Println(err)
		Error(w, http.StatusNotFound, "Not found")
		return
	}

	Success(w, http.StatusOK, stat)
}

func RedirectApi(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	stat, err := service.GetByShort(code)

	if err != nil {
		log.Println(err)
		Error(w, http.StatusNotFound, "Not found")
		return
	}

	Success(w, http.StatusOK, stat)
}
