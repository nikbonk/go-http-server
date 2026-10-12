package api

import (
	"encoding/json"
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

	rowsAffected, err := c.db.UpdateUserMembership(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
