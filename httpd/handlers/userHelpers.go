package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// CheckIfUsernameIsAvailable - return a boolean value to see if a given username is
// already taken
func CheckIfUsernameIsAvailable(w http.ResponseWriter, r *http.Request) {
	requestedUsername := utils.ParseRouteParameter(r.URL.Path, "/users/username-is-available/")
	if strings.TrimSpace(requestedUsername) == "" {
		sendSystemError(w, errors.New("please pass a valid username for us to get you your requested content"))
		return
	}

	isAvailable, err := models.UsernameIsAvailable(requestedUsername)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, isAvailable, 200)
}

// ChangePassword changes the password of the user as long as the checks on the new/old passwords go through
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the password change route"), 405)
		return
	}

	type requestPayload struct {
		OldPassword     string `json:"oldPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)

	// get the user from context
	user := models.User{}
	user.FindByID(r.Context().Value("user_id").(int))

	err := user.ChangePassword(payload.OldPassword, payload.NewPassword, payload.ConfirmPassword)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	user.FindByID(r.Context().Value("user_id").(int))

	SendResponse(w, true, 200)
}

// UpdateDescription updates the description for the user
func UpdateDescription(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the description update route"), 405)
		return
	}

	type requestPayload struct {
		Description string `json:"description"`
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)

	// get the user from context
	user := models.User{}
	user.FindByID(r.Context().Value("user_id").(int))

	err := user.UpdateDescription(payload.Description)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	user.FindByID(r.Context().Value("user_id").(int))

	SendResponse(w, user.Description, 200)
}
