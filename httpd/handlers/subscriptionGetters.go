package handlers

import (
	"net/http"
	"strings"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// GetSubscriptions - gets a user's posttagvotes.
func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	communityURL := strings.TrimSpace(r.FormValue("community"))
	communityID := 0
	if communityURL != "" {
		c := models.Community{}
		if err := c.FindByURL(communityURL); err == nil {
			communityID = c.ID
		}
	}

	subscriptions, err := u.GetSubscriptions(communityID)
	if err != nil {
		SendResponse(w, utils.MakeError(err.Error()), 500)
		return
	}

	output := map[string]interface{}{
		"subscriptions": subscriptions,
	}
	SendResponse(w, output, 200)
}
