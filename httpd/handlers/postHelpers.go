package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// CheckIfURLIsAvailable - return a boolean value to see if a given URL is
// already taken
func CheckIfURLIsAvailable(w http.ResponseWriter, r *http.Request) {
	requestedURL := utils.ParseRouteParameter(r.URL.Path, "/post/url-is-available/")
	if strings.TrimSpace(requestedURL) == "" {
		sendSystemError(w, errors.New("please pass a valid URL for us to get you your requested content"))
		return
	}

	Log.Info("Checking if URL is available: " + requestedURL)
	isAvailable, err := models.URLIsAvailable(requestedURL)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, isAvailable, 200)
}

func CheckExternalURLTitle(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		URL string `json:"url"`
	}

	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the check external url route"), 405)
		return
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	if payload.URL == "" {
		sendSystemError(w, errors.New("`url` cannot be empty"))
	}

	title, err := models.ParseExternalURL(strings.TrimSpace(payload.URL))
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, title, 200)
}
