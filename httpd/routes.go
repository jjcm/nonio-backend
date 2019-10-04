package main

import (
	"net/http"

	"github.com/jjcm/soci-backend/httpd/handlers"
)

func openRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/":                      handlers.Home,
		"/register":              handlers.Register,
		"/login":                 handlers.Login,
		"/login-social":          handlers.LoginSocial,
		"/login-social/callback": handlers.LoginSocialCallback,
		"/info":                  handlers.Info,
	}

	return routes
}

func protectedRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/protected": handlers.GetTokenDetails,
		"/posts":     handlers.GetPosts,
	}

	return routes
}
