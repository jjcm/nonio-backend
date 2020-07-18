package handlers

import (
	"encoding/json"
	"net/http"

	"soci-backend/models"
)

// Login try and log a user in, if successful generate a JWT token and return that
func Login(w http.ResponseWriter, r *http.Request) {
	// any non GET handlers need to attach CORS headers. I always forget about that
	CorsAdjustments(&w)

	// silly AJAX prflight, here's where we can put in the CORS requirements
	if r.Method == "OPTIONS" {
		SendResponse(w, "", 200)
		return
	}

	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to the login route"), 405)
		return
	}

	requestUser := models.User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	if requestUser.Email == "" {
		SendResponse(w, MakeError("both password and email are required"), 400)
		return
	}

	u := models.User{}
	err := u.FindByEmail(requestUser.Email)
	if err != nil {
		sendNotFound(w, err)
		return
	}

	err = u.Login(requestUser.Password)
	if err != nil {
		sendNotFound(w, err)
		return
	}

	token, err := tokenCreator(u.Email)
	if err != nil {
		SendResponse(w, MakeError("There was an error signing your JWT token: "+err.Error()), 500)
		return
	}

	response := map[string]string{
		"token":    token,
		"username": u.Username,
	}
	SendResponse(w, response, 200)
}
