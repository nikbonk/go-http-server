package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/nikbonk/go-http-server/internal/auth"
)

func (c *ApiConfig) PolkaUpgradeUserEventHandler(w http.ResponseWriter, r *http.Request) {
	polkaExpectedEvent := "user.upgraded"

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if apiKey != c.polkaKey {
		http.Error(w, "Action not allowed ", http.StatusUnauthorized)
		return
	}

	request := polkaWebhookRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.Event != polkaExpectedEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(request.Data.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.db.UpdateUserMembership(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
