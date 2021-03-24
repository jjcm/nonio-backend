package handlers

import (
	"errors"
	"net/http"
	"strings"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// CheckIfURLIsAvailable - return a boolean value to see if a given URL is
// already taken
func CheckIfURLIsAvailable(w http.ResponseWriter, r *http.Request) {
	requestedURL := utils.ParseRouteParameter(r.URL.Path, "/posts/url-is-available/")
	if strings.TrimSpace(requestedURL) == "" {
		sendSystemError(w, errors.New("please pass a valid URL for us to get you your requested content"))
		return
	}

	isAvailable, err := models.URLIsAvailable(requestedURL)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, isAvailable, 200)
}
