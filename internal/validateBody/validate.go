package json

import (
	"encoding/json"
	"log"
	"net/http"

	profanityCheck "github.com/nikbonk/go-http-server/internal/profanityCheck"
)

type request struct {
	Body string `json:"body"`
}

type response struct {
	Error       string `json:"error"`
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func ValidateBodyHandler(w http.ResponseWriter, r *http.Request) {
	const maxChars = 140
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	request := request{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(request.Body) > maxChars {
		respondWithError(w, http.StatusBadRequest, "140 character limit reached.", nil)
		return
	}

	containsProfanity, cleanedBody := profanityCheck.ContainsProfanity(request.Body)
	respondWithJSON(w, http.StatusOK, response{
		Error:       "",
		Valid:       !containsProfanity,
		CleanedBody: cleanedBody,
	})

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
