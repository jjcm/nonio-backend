package handlers

import (
	"encoding/json"
	"net/http"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// Login try and log a user in, if successful generate a JWT token and return that
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("You can only POST to the login route"), 405)
		return
	}

	requestUser := models.User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)
	if requestUser.Email == "" {
		SendResponse(w, utils.MakeError("both password and email are required"), 400)
		return
	}

	u := models.User{}
	err := u.FindByEmail(requestUser.Email)
	if err != nil {
		Log.Errorf("error finding user by email: %v", requestUser.Email)
		sendNotFound(w, err)
		return
	}

	err = u.Login(requestUser.Password)
	if err != nil {
		// Info - there's an error logging in on the mobile browser, but I'm not sure why. I'm logging passwords IN PLAINTEXT only for myself here. I'm lazy so I'm testing on prod, but I also want anyone reading this to know I'm only doing it for myself.
		if requestUser.Email == "j@jjcm.org" {
			Log.Errorf("error logging in with password: %v", requestUser.Password)
			Log.Error(err)
		}
		sendNotFound(w, err)
		return
	}

	token, err := utils.TokenCreator(u.Email)
	if err != nil {
		SendResponse(w, utils.MakeError("there was an error signing your JWT token: "+err.Error()), 500)
		return
	}

	response := map[string]string{
		"token":    token,
		"username": u.Username,
	}
	SendResponse(w, response, 200)
}
