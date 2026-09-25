package main

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	// imagine system check before just returning a 200 OK
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	websiteHits := cfg.fileserverHits.Load()

	tmpl, err := template.ParseFiles("./html/metrics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, map[string]int32{
		"websiteHits": websiteHits,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func validateBodyHandler(w http.ResponseWriter, r *http.Request) {
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

	containsProfanity, cleanedBody := containsProfanity(request.Body)
	respondWithJSON(w, http.StatusOK, response{
		Error:       "",
		Valid:       !containsProfanity,
		CleanedBody: cleanedBody,
	})

}
