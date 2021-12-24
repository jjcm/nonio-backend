package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"io/ioutil"
	"net/http"
	"os"
	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// StripeWebhook create a new customer for user
func StripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the stripe webhook route"), http.StatusMethodNotAllowed)
		return
	}

	const MaxBodyBytes = int64(65536)

	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	body, err := ioutil.ReadAll(r.Body)

	endPointSecret := os.Getenv("WEBHOOK_ENDPOINT_SECRET")

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endPointSecret)

	if err != nil {
		sendSystemError(w, err)
		return
	}

	if event.Type == "payment_intent.succeeded" {
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			sendSystemError(w, err)
			return
		}
		u := models.User{}
		if err := u.FindByCustomerId(paymentIntent.Customer.ID); err != nil {
			sendSystemError(w, fmt.Errorf("find user by id: %v", err))
			return
		}
		if u.StripeCustomerID == "" {
			sendSystemError(w, errors.New("no customer for the user"))
			return
		}
		if err := u.UpdateAccountType("supporter"); err != nil {
			sendSystemError(w, fmt.Errorf("find user by id: %v", err))
			return
		}
	} else if event.Type == "payment_intent.payment_failed" {
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		u := models.User{}
		if err := u.FindByCustomerId(paymentIntent.Customer.ID); err != nil {
			sendSystemError(w, fmt.Errorf("find user by id: %v", err))
			return
		}
		if u.StripeCustomerID == "" {
			sendSystemError(w, errors.New("no customer for the user"))
			return
		}
		if err := u.UpdateAccountType("free"); err != nil {
			sendSystemError(w, fmt.Errorf("find user by id: %v", err))
			return
		}
	}
	SendResponse(w, true, 200)
}
