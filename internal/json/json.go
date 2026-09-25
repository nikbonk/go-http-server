package json

import (
	"encoding/json"
	"log"
	"net/http"
)

type Request struct {
	Body string `json:"body"`
}

type Response struct {
	Error       string `json:"error"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, errorMessage string, err error) {
	if err != nil {
		log.Println(err)
	}
	if statusCode > 499 {
		log.Println("Server error:", statusCode, errorMessage)
	}

	RespondWithJSON(w, statusCode, Response{
		Error: errorMessage,
		Valid: false,
	})
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, r Response) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(r)
}
