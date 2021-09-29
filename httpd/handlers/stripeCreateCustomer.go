package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"

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
	params := &stripe.CustomerParams{
		Email: stripe.String(payload.Email),
	}

	c, err := customer.New(params)
	if err != nil {
		sendSystemError(w, fmt.Errorf("new customer: %v", err))
		return
	}

	// You should store the ID of the customer in your database alongside your
	// users. This sample uses cookies to simulate auth.
	http.SetCookie(w, &http.Cookie{
		Name:  "customer",
		Value: c.ID,
	})

	output := map[string]interface{}{
		"customer": c,
	}
	SendResponse(w, output, 200)
}
