package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jjcm/soci-backend/models"
)

// RegisterPayload This is the shape of the JSON payload that will be sent to
// the API to register a new user
type RegisterPayload struct {
	Email    string `json:"email"`
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
		SendResponse(w, "You can only POST to the registration route", 405)
		return
	}

	var payload RegisterPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	Log.Info(payload)

	// let's check and see if the registered email is already registered
	u := models.User{}
	err := u.FindByEmail(payload.Email)
	if err == nil {
		// err is nil, meaning there was not a problem looking up this user, so one was found
		SendResponse(w, "This email has already been registered", 500)
		return
	}

	Log.Info("now creating new user")
	models.CreateUser(payload.Email, payload.Password)

	// send a token
	token, err := tokenCreator(u.Email)
	if err != nil {
		SendResponse(w, "There was an error signing your JWT token: "+err.Error(), 500)
		return
	}

	SendResponse(w, token, 200)
}
