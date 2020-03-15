package handlers

import (
	"encoding/json"
	"net/http"

	"soci-backend/models"
)

// RegisterPayload This is the shape of the JSON payload that will be sent to
// the API to register a new user
type RegisterPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register Save a new user in the DB
func Register(w http.ResponseWriter, r *http.Request) {
	// any non GET handlers need to attach CORS headers. I always forget about that
	CorsAdjustments(&w)

	// silly AJAX prflight, here's where we can put in the CORS requirements
	if r.Method == "OPTIONS" {
		SendResponse(w, "", 200)
		return
	}

	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to the registration route"), 405)
		return
	}

	var payload RegisterPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	if payload.Username == "" || payload.Password == "" || payload.Email == "" {
		SendResponse(w, MakeError("username, password and email are all required"), 400)
		return
	}

	// let's check and see if the registered email is already taken
	u := models.User{}
	err := u.FindByEmail(payload.Email)
	if err == nil {
		// err is nil, meaning there was not a problem looking up this user, so one was found
		SendResponse(w, MakeError("This email has already been registered"), 500)
		return
	}
	// let's check and see if the registered username is already taken
	err = u.FindByUsername(payload.Username)
	if err == nil {
		// err is nil, meaning there was not a problem looking up this user, so one was found
		SendResponse(w, MakeError("This username has already been registered"), 500)
		return
	}

	Log.Info("now creating new user")
	err = models.CreateUser(payload.Email, payload.Username, payload.Password)
	if err != nil {
		// err is nil, meaning there was not a problem looking up this user, so one was found
		SendResponse(w, MakeError("Error registering user: "+err.Error()), 500)
		return
	}

	// send a token
	token, err := tokenCreator(payload.Email)
	if err != nil {
		SendResponse(w, MakeError("There was an error signing your JWT token: "+err.Error()), 500)
		return
	}

	response := map[string]string{
		"token":    token,
		"username": payload.Username,
	}
	SendResponse(w, response, 200)
}
