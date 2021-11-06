package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"
	"soci-backend/models"
	"time"

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
		PaymentMethodID string `json:"paymentMethodId"`
		PriceID         string `json:"priceId"`
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

	// check if the user has any active subscriptions
	listParams := &stripe.SubscriptionListParams{
		Customer: u.StripeCustomerID,
		Status:   "all",
	}
	listParams.AddExpand("data.default_payment_method")

	iter := sub.List(listParams)
	subscriptions := iter.SubscriptionList().Data
	if len(subscriptions) > 0 {
		Log.Info(fmt.Sprintf("%v subscriptions found for %v. Cancelling...", len(subscriptions), u.Username))
		for i := 0; i < len(subscriptions); i++ {
			if subscriptions[i].Status != "canceled" && subscriptions[i].Status != "incomplete_expired" {
				Log.Info(fmt.Sprintf("Cancelling %v. Status %v", subscriptions[i].ID, subscriptions[i].Status))
				subscription, err := sub.Cancel(subscriptions[i].ID, nil)
				if err != nil {
					sendSystemError(w, fmt.Errorf("cancel subscription: %v", err))
					return
				}
				Log.Info(fmt.Sprintf("Cancelled %v", subscription.ID))
			}
		}
		Log.Info("Old subscriptions cancelled.")
		time.Sleep(5 * time.Second)
	}

	// create subscription
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(u.StripeCustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(payload.PriceID),
			},
		},
		PaymentBehavior:      stripe.String("default_incomplete"),
		DefaultPaymentMethod: stripe.String(payload.PaymentMethodID),
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
