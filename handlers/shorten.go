package handlers

import (
	"encoding/json"
	"net/http"
	"shortener/service"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func Shorten(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)

	if url == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		encoder.Encode(ErrorResponse{"field url is required"})
		return
	}

	short, err := service.StoreUrl(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		encoder.Encode(ErrorResponse{err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	encoder.Encode(SuccessResponse{
		Success: true,
		Data:    short,
	})
}
