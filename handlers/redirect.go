package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"shortener/service"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	short := r.PathValue("redirect")

	url, err := service.Retrieve(short)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
