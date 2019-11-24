package main

import (
	"net/http"

	"github.com/jjcm/soci-backend/httpd/handlers"
)

func openRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/":                        handlers.Home,
		"/register":                handlers.Register,
		"/login":                   handlers.Login,
		"/login-social":            handlers.LoginSocial,
		"/login-social/callback":   handlers.LoginSocialCallback,
		"/info":                    handlers.Info,
		"/posts/":                  handlers.GetPostByURL,
		"/posts/url-is-available/": handlers.CheckIfURLIsAvailable,
	}

	return routes
}

func protectedRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/protected": handlers.GetTokenDetails,

		// post routes
		"/posts/new":               handlers.GetNewestPosts,
		"/posts":                   handlers.GetPosts,
		"/post/create":             handlers.CreatePost,

		// tag routes
		"/tags": handlers.GetTags,
	}

	return routes
}
