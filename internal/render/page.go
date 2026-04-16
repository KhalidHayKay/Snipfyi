package render

import (
	"html/template"
	"log"
	"net/http"
)

type ViewData struct {
	Title string
	Page  string
	Data  any
	Error string
}

// Page renders an HTML template with the given data
func Page(w http.ResponseWriter, page string, data ViewData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles(
		"templates/layouts/layout.html",
		"templates/pages/"+page,
	)
	if err != nil {
		log.Printf("ERROR: Failed to parse templates: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("ERROR: Failed to execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SinglePage(w http.ResponseWriter, page string, data ViewData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	t, err := template.ParseFiles("templates/pages/" + page)
	if err != nil {
		log.Printf("ERROR: Failed to parse template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("ERROR: Failed to execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ErrorPage renders an error HTML page
func ErrorPage(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	data := map[string]any{
		"Title":   "Error",
		"Message": message,
		"Status":  status,
	}

	// Determine which error template to use based on status code
	errorTemplate := "templates/error/error.html"
	if status == http.StatusNotFound {
		errorTemplate = "templates/error/not-found.html"
	}

	t, err := template.ParseFiles(
		"templates/layouts/layout.html",
		errorTemplate,
	)
	if err != nil {
		// Fallback to plain text error if template fails
		log.Printf("ERROR: Failed to parse error template: %v", err)
		w.Write([]byte(message))
		return
	}

	if err := t.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("ERROR: Failed to render error page: %v", err)
		w.Write([]byte(message))
	}
}
