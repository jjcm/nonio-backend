package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"

	"github.com/stripe/stripe-go/v72/sub"
)

// StripeCancelSubscription cancel a subscription
func StripeCancelSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the cancel subscription route"), http.StatusMethodNotAllowed)
		return
	}

	type requestPayload struct {
		SubscriptionID string `json:"subscriptionId"`
	}
	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendSystemError(w, fmt.Errorf("decode request payload: %v", err))
		return
	}

	subscription, err := sub.Cancel(payload.SubscriptionID, nil)
	if err != nil {
		sendSystemError(w, fmt.Errorf("cancel subscription: %v", err))
		return
	}

	output := map[string]interface{}{
		"subscription": subscription,
	}
	SendResponse(w, output, 200)
}
