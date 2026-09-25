package http

import (
	"encoding/json"
	"net/http"

	nbjson "github.com/nikbonk/go-http-server/internal/json"
	"github.com/nikbonk/go-http-server/internal/profanity"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// imagine system check before just returning a 200 OK
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func ValidateBodyHandler(w http.ResponseWriter, r *http.Request) {
	const maxChars = 140
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	request := nbjson.Request{}
	err := decoder.Decode(&request)
	if err != nil {
		nbjson.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(request.Body) > maxChars {
		nbjson.RespondWithError(w, http.StatusBadRequest, "140 character limit reached.", nil)
		return
	}

	containsProfanity, cleanedBody := profanity.ContainsProfanity(request.Body)
	nbjson.RespondWithJSON(w, http.StatusOK, nbjson.Response{
		Error:       "",
		Valid:       !containsProfanity,
		CleanedBody: cleanedBody,
	})

}
