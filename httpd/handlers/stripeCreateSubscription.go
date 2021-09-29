package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

// StripeCreateSubscription create a new subscription with fixed price for user
func StripeCreateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the stripe create customer route"), http.StatusMethodNotAllowed)
		return
	}

	type requestPayload struct {
		PriceID string `json:"priceId"`
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendSystemError(w, fmt.Errorf("decode request payload: %v", err))
		return
	}

	// read customer from cookie to simulate auth
	cookie, _ := r.Cookie("customer")
	// Create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(cookie.Value),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(payload.PriceID),
			},
		},
		PaymentBehavior: stripe.String("default_incomplete"),
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	newSub, err := sub.New(subscriptionParams)
	if err != nil {
		sendSystemError(w, fmt.Errorf("new subscription: %v", err))
		return
	}

	output := map[string]interface{}{
		"subscriptionId": newSub.ID,
		"clientSecret":   newSub.LatestInvoice.PaymentIntent.ClientSecret,
	}
	SendResponse(w, output, 200)
}
