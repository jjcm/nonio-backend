package handlers

import (
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"
	"soci-backend/models"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
)

// StripeCreateCustomer create a new customer for user
func StripeCreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the stripe create customer route"), http.StatusMethodNotAllowed)
		return
	}

	uid := r.Context().Value("user_id").(int)

	u := models.User{}
	if err := u.FindByID(uid); err != nil {
		sendSystemError(w, fmt.Errorf("find user by id: %v", err))
		return
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(u.Email),
	}

	var c *stripe.Customer
	var err error
	if u.StripeCustomerID != "" {
		_, err = customer.Get(u.StripeCustomerID, nil)
		if err != nil {
			sendSystemError(w, fmt.Errorf("get customer: %v", err))
			return
		}
	} else {
		c, err = customer.New(params)
		if err != nil {
			sendSystemError(w, fmt.Errorf("new customer: %v", err))
			return
		}

		// update the customer id to user
		if err := u.UpdateStripeCustomerID(c.ID); err != nil {
			sendSystemError(w, fmt.Errorf("update stripe customer id: %v", err))
			return
		}
	}

	SendResponse(w, true, 200)
}
