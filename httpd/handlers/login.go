package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jjcm/soci-backend/models"
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

	u := models.User{}
	err := u.FindByEmail(requestUser.Email)
	if err != nil {
		SendResponse(w, MakeError("Those credentials do not match our records"), 404)
		return
	}

	err = u.Login(requestUser.Password)
	if err != nil {
		SendResponse(w, MakeError("Those credentials do not match our records"), 404)
		return
	}

	token, err := tokenCreator(u.Email)
	if err != nil {
		SendResponse(w, MakeError("There was an error signing your JWT token: "+err.Error()), 500)
		return
	}

	response := map[string]string{
		"token": token,
	}
	SendResponse(w, response, 200)
}
