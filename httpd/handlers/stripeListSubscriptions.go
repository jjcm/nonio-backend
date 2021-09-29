package handlers

import (
	"net/http"
	"soci-backend/httpd/utils"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

// StripeListSubscriptions list the subscriptions
func StripeListSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendResponse(w, utils.MakeError("you can only GET to the list subscriptions route"), http.StatusMethodNotAllowed)
		return
	}

	// Read customer from cookie to simulate auth
	cookie, _ := r.Cookie("customer")

	params := &stripe.SubscriptionListParams{
		Customer: cookie.Value,
		Status:   "all",
	}
	params.AddExpand("data.default_payment_method")

	iter := sub.List(params)
	output := map[string]interface{}{
		"subscriptions": iter.SubscriptionList(),
	}
	SendResponse(w, output, 200)
}
