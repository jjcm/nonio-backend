package handlers

import (
	"net/http"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// GetNotifications - gets a user's notifications as a list of comments
func GetNotifications(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	notifications, err := u.GetNotifications()
	if err != nil {
		SendResponse(w, utils.MakeError(err.Error()), 500)
		return
	}

	output := map[string]interface{}{
		"notifications": notifications,
	}
	SendResponse(w, output, 200)
}

// GetUnreadNotificationCount - gets a user's count of unread notifications
func GetUnreadNotificationCount(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	notificationCount, err := u.GetUnreadNotificationCount()
	if err != nil {
		SendResponse(w, utils.MakeError(err.Error()), 500)
		return
	}

	SendResponse(w, notificationCount, 200)
}
