package handlers

import (
	"encoding/json"
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

	type requestPayload struct {
		Email string `json:"email"`
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		sendSystemError(w, fmt.Errorf("decode request payload: %v", err))
		return
	}

	u := models.User{}
	if err := u.FindByEmail(payload.Email); err != nil {
		sendSystemError(w, fmt.Errorf("find user by email: %v", err))
		return
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(payload.Email),
	}

	var c *stripe.Customer
	var err error
	if u.StripeCustomerID != "" {
		c, err = customer.Get(u.StripeCustomerID, params)
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
		if err := u.UpdateStripCustomerID(u.StripeCustomerID); err != nil {
			sendSystemError(w, fmt.Errorf("update strip customer id: %v", err))
			return
		}
	}

	output := map[string]interface{}{
		"customer": c,
	}
	SendResponse(w, output, 200)
}
