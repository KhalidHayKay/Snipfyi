package handlers

import (
	"log"
	"net/http"
	"smply/service"
	"text/template"
)

func RequestApiKey(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	if email == "" {
		Error(w, http.StatusUnprocessableEntity, "'email' is a required field")
		return
	}

	err := service.RequestApiKey(r.Context(), email)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusNotFound, "Not found")
		return
	}

	Success(w, http.StatusOK, nil)
}

func CreateApiKey(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token == "" {
		Error(w, http.StatusNotFound, "cannot find token")
		return
	}

	key, err := service.CreateApiKey(r.Context(), token)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusNotFound, "Not found")
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

// func ApiKeyPage(w http.ResponseWriter, r *http.Request) {
//     // validate token from query param, look up and delete key from DB

//     // render as a standalone file, not through the layout
//     t := template.Must(template.ParseFiles("templates/pages/api-key.html"))
//     t.Execute(w, data)
// }
