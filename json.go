package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type request struct {
	Body string `json:"body"`
}

type response struct {
	Error       string `json:"error"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, statusCode int, errorMessage string, err error) {
	if err != nil {
		log.Println(err)
	}
	if statusCode > 499 {
		log.Println("Server error:", statusCode, errorMessage)
	}

	respondWithJSON(w, statusCode, response{
		Error: errorMessage,
		Valid: false,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, r response) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}
