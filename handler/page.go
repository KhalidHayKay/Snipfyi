package handler

import (
	"errors"
	"log"
	"net/http"
	"smply/internal/queue"
	"smply/internal/render"
	"smply/internal/service"
	"smply/model"
	"smply/utils"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "home.html", render.ViewData{
		Title: "Home",
		Page:  "home",
	})
}

func ShortenPage(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "shorten.html", render.ViewData{
		Title: "Shorten URL",
		Page:  "shorten",
	})
}

func HomeShorten(w http.ResponseWriter, r *http.Request) {
	shortenForm(w, r, "Home", "home", "home.html")
}

func ShortenPageShorten(w http.ResponseWriter, r *http.Request) {
	shortenForm(w, r, "Shorten URL", "shorten", "shorten.html")
}

func shortenForm(w http.ResponseWriter, r *http.Request, title, page, file string) {
	url := r.FormValue("url")
	alias := r.FormValue("alias")

	data := render.ViewData{
		Title: title,
		Page:  page,
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

	result, err := service.StoreUrl(r.Context(), url, alias)
	if err != nil {
		log.Println(err)

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

func ApiPage(w http.ResponseWriter, r *http.Request) {
	render.Page(w, "api.html", render.ViewData{
		Title: "API",
		Page:  "api",
	})
}

func ResolveRedirect(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	url, err := service.GetByAlias(r.Context(), alias)

	if err != nil {
		log.Println(err)

		if err == pgx.ErrNoRows {
			render.ErrorPage(w, http.StatusNotFound, "Not found")
			return
		}

		render.ErrorPage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// stats update called here because redirect might now happen in `GetByAlias` call.
	// For example in API call.
	err = queue.EnqueueStatsUpdate(
		r.Context(),
		url.Alias,
		r.Referer(),
		r.UserAgent(),
		"", //utils.GetIPAddress(r),
		time.Now(),
	)
	if err != nil {
		log.Println(err)
	}

	http.Redirect(w, r, url.Original, http.StatusFound)
}

func StatsPage(w http.ResponseWriter, r *http.Request) {
	alias := chi.URLParam(r, "alias")

	stat, err := service.GetStats(r.Context(), alias)
	if err != nil {
		log.Println(err)
		render.ErrorPage(w, http.StatusNotFound, "Not found")
		return
	}

	render.Page(w, "stats.html", render.ViewData{
		Title: "URL Stats",
		Page:  "stats",
		Data: map[string]model.Url{
			"Stats": stat,
		},
	})
}
