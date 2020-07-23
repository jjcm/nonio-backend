package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"soci-backend/models"
)

// SubscriptionDeletionRequest is the shape of the JSON request that is needed to remove a sub for a tag
type SubscriptionDeletionRequest struct {
	TagName string `json:"tag"`
}

// DeleteSubscription removes a sub for a tag
func DeleteSubscription(w http.ResponseWriter, r *http.Request) {
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

	// figure out what tag they're trying to remove
	tag := models.Tag{}
	tag.FindByTagName(request.TagName)

	// check if the Subscription exists in the db
	if err := user.DeleteSubscription(tag); err != nil {
		sendSystemError(w, fmt.Errorf("Error deleting the subscription: %v", err))
		return
	}

	SendResponse(w, true, 200)
}
