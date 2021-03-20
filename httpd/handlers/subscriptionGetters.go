package handlers

import (
	"net/http"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// GetSubscriptions - gets a user's posttagvotes.
func GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	subscriptions, err := u.GetSubscriptions()
	if err != nil {
		SendResponse(w, utils.MakeError(err.Error()), 500)
		return
	}

	output := map[string]interface{}{
		"subscriptions": subscriptions,
	}
	SendResponse(w, output, 200)
}
