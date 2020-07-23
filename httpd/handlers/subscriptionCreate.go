package handlers

import (
	"encoding/json"
	"net/http"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// SubscriptionAdditionRequest is the shape of the JSON request that is needed to add a sub for a tag
type SubscriptionAdditionRequest struct {
	TagName string `json:"tag"`
}

// CreateSubscription adds a sub for a tag
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("You can only POST to the AddSubscription route"), 405)
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

	// if the Subscription already exists, then just return true
	subscription.FindSubscription(tag.ID, user.ID)
	if subscription.ID > 0 {
		SendResponse(w, subscription, 200)
		return
	}

	// otherwise, let's make a subscription!
	subscription, err := user.CreateSubscription(tag)
	if err != nil {

	}

	SendResponse(w, true, 200)
}
