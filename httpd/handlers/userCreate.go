package handlers

import (
	"encoding/json"
	"net/http"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// Register saves a new user in the DB
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the registration route"), 405)
		return
	}

	type requestPayload struct {
		Email              string  `json:"email"`
		Username           string  `json:"username"`
		Password           string  `json:"password"`
		SubscriptionAmount float64 `json:"subscriptionAmount"`
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)

	if payload.Username == "" || payload.Password == "" || payload.Email == "" || payload.SubscriptionAmount == 0 {
		SendResponse(w, utils.MakeError("username, password, email, and subscription amount are all required"), 400)
		return
	}

	// let's check and see if the registered email is already taken
	u := models.User{}
	err := u.FindByEmail(payload.Email)
	if err == nil {
		// err is nil, meaning there was not a problem looking up this user, so one was found
		SendResponse(w, utils.MakeError("this email has already been registered"), 500)
		return
	}
	// let's check and see if the registered username is already taken
	err = u.FindByUsername(payload.Username)
	if err == nil {
		// err is nil, meaning there was not a problem looking up this user, so one was found
		SendResponse(w, utils.MakeError("this username has already been registered"), 500)
		return
	}

	Log.Info("now creating new user")
	_, err = models.UserFactory(payload.Email, payload.Username, payload.Password, payload.SubscriptionAmount)
	if err != nil {
		// err is nil, meaning there was not a problem looking up this user, so one was found
		SendResponse(w, utils.MakeError("error registering user: "+err.Error()), 500)
		return
	}

	// send a token
	token, err := utils.TokenCreator(payload.Email)
	if err != nil {
		SendResponse(w, utils.MakeError("there was an error signing your JWT token: "+err.Error()), 500)
		return
	}

	response := map[string]string{
		"token":    token,
		"username": payload.Username,
	}
	SendResponse(w, response, 200)
}
