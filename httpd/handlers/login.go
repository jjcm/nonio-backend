package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jjcm/soci-backend/models"
)

// Login try and log a user in, if successful generate a JWT token and return that
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, "You can only POST to the login route", 405)
		return
	}

	requestUser := models.User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&requestUser)

	u := models.User{}
	err := u.FindByEmail(requestUser.Email)
	if err != nil {
		SendResponse(w, "Those credentials do not match our records", 404)
		return
	}

	err = u.Login(requestUser.Password)
	if err != nil {
		SendResponse(w, "Those credentials do not match our records", 404)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     u.Email,
		"expiresAt": time.Now().Add(time.Minute * 10).Unix(), // tokens are valid for 10 minutes?
	})
	tokenString, err := token.SignedString(HmacSampleSecret)
	if err != nil {
		SendResponse(w, "There was an error signing your JWT token: "+err.Error(), 500)
		return
	}

	SendResponse(w, tokenString, 200)
}
