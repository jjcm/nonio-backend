package handlers

import (
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gopkg.in/danilopolani/gocialite.v1"
)

var gocial = gocialite.NewDispatcher()

// LoginSocial try and authenticate a user over oauth
func LoginSocial(w http.ResponseWriter, r *http.Request) {
	// other providers documented here
	// https://github.com/danilopolani/gocialite/wiki/Multi-provider-example

	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     os.Getenv("OAUTH_ID"),
			"clientSecret": os.Getenv("OAUTH_SECRET"),
			"redirectURL":  "http://localhost:9000/login-social/callback",
		},
	}
	providerScopes := map[string][]string{
		"github": []string{""},
	}

	providerData := providerSecrets["github"]
	actualScopes := providerScopes["github"]
	authURL, err := gocial.New().
		Driver("github").
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	if err != nil {
		Log.Error(err.Error())
		return
	}

	http.Redirect(w, r, authURL, http.StatusSeeOther)
}

// LoginSocialCallback callback for oauth app
func LoginSocialCallback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	code := r.FormValue("code")
	state := r.FormValue("state")
	Log.Info(code)
	Log.Info(state)

	// Handle callback and check for errors
	user, _, err := gocial.Handle(state, code)
	if err != nil {
		Log.Error(err.Error())
		return
	}

	// Print in terminal user information
	Log.Info(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":     user.Email,
		"expiresAt": time.Now().Add(time.Minute * 10).Unix(), // tokens are valid for 10 minutes?
	})
	tokenString, err := token.SignedString(HmacSampleSecret)
	if err != nil {
		SendResponse(w, "There was an error signing your JWT token: "+err.Error(), 500)
		return
	}

	SendResponse(w, tokenString, 200)
}
