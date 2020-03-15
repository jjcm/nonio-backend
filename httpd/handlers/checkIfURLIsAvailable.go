package handlers

import (
	"errors"
	"net/http"
	"strings"

	"soci-backend/models"
)

// CheckIfURLIsAvailable - return a boolean value to see if a given URL is
// already taken
func CheckIfURLIsAvailable(w http.ResponseWriter, r *http.Request) {
	requestedURL := parseRouteParameter(r.URL.Path, "/posts/url-is-available/")
	if strings.TrimSpace(requestedURL) == "" {
		sendSystemError(w, errors.New("Please pass a valid URL for us to get you your requested content"))
		return
	}

	isAvailable, err := models.URLIsAvailable(requestedURL)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, isAvailable, 200)
}
