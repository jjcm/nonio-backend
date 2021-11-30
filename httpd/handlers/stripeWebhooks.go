package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"
)

// StripeWebhook create a new customer for user
func StripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the stripe webhook route"), http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON payload
	type requestPayload struct {
		Type string `json:"type"`
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendSystemError(w, fmt.Errorf("decode request payload: %v", err))
		return
	}

	Log.Info(payload.Type)

	SendResponse(w, true, 200)
}
