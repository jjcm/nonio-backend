package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jjcm/soci-backend/models"
)

// CheckIfURLIsAvailable - return a boolean value to see if a given URL is
// already taken
func CheckIfURLIsAvailable(w http.ResponseWriter, r *http.Request) {
	requestedUsername := parseRouteParameter(r.URL.Path, "/users/username-is-available/")
	if strings.TrimSpace(requestedURL) == "" {
		sendSystemError(w, errors.New("Please pass a valid username for us to get you your requested content"))
		return
	}

	isAvailable, err := models.UsernameIsAvailable(requestedUsername)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, isAvailable, 200)
}
