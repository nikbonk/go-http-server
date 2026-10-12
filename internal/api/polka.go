package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (c *ApiConfig) PolkaUpgradeUserEventHandler(w http.ResponseWriter, r *http.Request) {
	polkaExpectedEvent := "user.upgraded"

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
