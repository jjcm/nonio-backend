package main

import (
	"net/http"

	"github.com/jjcm/soci-backend/httpd/handlers"
)

func openRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := make(map[string]func(http.ResponseWriter, *http.Request))
	routes["/"] = handlers.Home
	routes["/register"] = handlers.Register
	routes["/login"] = handlers.Login
	routes["/login-social"] = handlers.LoginSocial
	routes["/login-social/callback"] = handlers.LoginSocialCallback

	return routes
}

func protectedRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := make(map[string]func(http.ResponseWriter, *http.Request))
	routes["/protected"] = handlers.GetTokenDetails

	return routes
}
