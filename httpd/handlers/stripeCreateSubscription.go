package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"
	"soci-backend/models"

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

	uid := r.Context().Value("user_id").(int)

	u := models.User{}
	if err := u.FindByID(uid); err != nil {
		sendSystemError(w, fmt.Errorf("find user by id: %v", err))
		return
	}
	if u.StripeCustomerID == "" {
		sendSystemError(w, errors.New("no customer for the user"))
		return
	}

	// check if the subscription is existed
	listParams := &stripe.SubscriptionListParams{
		Customer: u.StripeCustomerID,
		Status:   "all",
	}
	listParams.AddExpand("data.default_payment_method")

	iter := sub.List(listParams)
	subscriptions := iter.SubscriptionList().Data
	if len(subscriptions) > 0 {
		sendSystemError(w, errors.New("subscription already exists"))
		return
	}

	// create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(u.StripeCustomerID),
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
