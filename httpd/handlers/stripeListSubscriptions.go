package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"
	"soci-backend/models"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
)

// StripeListSubscriptions list the subscriptions
func StripeListSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendResponse(w, utils.MakeError("you can only GET to the list subscriptions route"), http.StatusMethodNotAllowed)
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

	params := &stripe.SubscriptionListParams{
		Customer: u.StripeCustomerID,
		Status:   "all",
	}
	params.AddExpand("data.default_payment_method")

	iter := sub.List(params)
	output := map[string]interface{}{
		"subscriptions": iter.SubscriptionList(),
	}
	SendResponse(w, output, 200)
}
