package handlers

import (
	"encoding/json"
	"net/http"

	"fmt"
	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// CreateSubscription adds a sub for a tag
func CreateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the create subscription route"), 405)
		return
	}

	type requestPayload struct {
		TagName   string `json:"tag"`
		Community string `json:"community"`
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)

	// get the user from context
	user := models.User{}
	user.FindByID(r.Context().Value("user_id").(int))

	communityID := 0
	if payload.Community != "" {
		// Remove @ prefix if present
		communityURL := payload.Community
		if len(communityURL) > 0 && communityURL[0] == '@' {
			communityURL = communityURL[1:]
		}
		c := models.Community{}
		if err := c.FindByURL(communityURL); err == nil {
			communityID = c.ID
		}
	}

	tag := models.Tag{}
	tag.FindByTagName(payload.TagName, communityID)

	if tag.ID == 0 {
		SendResponse(w, utils.MakeError("tag not found"), 404)
		return
	}

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
		sendSystemError(w, fmt.Errorf("couldn't create a subscription: %v", err))
		return
	}

	SendResponse(w, true, 200)
}
