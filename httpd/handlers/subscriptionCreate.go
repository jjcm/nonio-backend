package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"soci-backend/models"
)

// SubscriptionAdditionRequest is the shape of the JSON request that is needed to add a sub for a tag
type SubscriptionAdditionRequest struct {
	TagName string `json:"tag"`
}

// CreateSubscription adds a sub for a tag
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	// any non GET handlers need to attach CORS headers
	CorsAdjustments(&w)

	// silly AJAX prflight, here's where we can put in the CORS requirements
	if r.Method == "OPTIONS" {
		SendResponse(w, "", 200)
		return
	}

	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to the AddSubscription route"), 405)
		return
	}

	// decode the request parameters 'tag'
	var request SubscriptionAdditionRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)

	// get the user from context
	user := models.User{}
	user.FindByID(r.Context().Value("user_id").(int))

	tag := models.Tag{}
	tag.FindByTagName(request.TagName)

	subscription := models.Subscription{}
	// check if the Subscription exists in the db
	if err := subscription.FindSubscription(tag.ID, user.ID); err != nil {
		sendSystemError(w, fmt.Errorf("Query subscription: %v", err))
		return
	}

	// if the Subscription already exists, then just return true
	if subscription.ID > 0 {
		SendResponse(w, subscription, 200)
		return
	}

	// otherwise, let's make a subscription!
	subscription, err := user.CreateSubscription(tag)
	if err != nil {

	}

	SendResponse(w, subscription, 200)
}
