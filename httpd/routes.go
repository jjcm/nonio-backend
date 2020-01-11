package main

import (
	"net/http"

	"github.com/jjcm/soci-backend/httpd/handlers"
)

func openRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/":                             handlers.Home,
		"/register":                     handlers.Register,
		"/login":                        handlers.Login,
		"/login-social":                 handlers.LoginSocial,
		"/login-social/callback":        handlers.LoginSocialCallback,
		"/info":                         handlers.Info,
		"/posts/user/":                  handlers.GetPostsByAuthor,
		"/posts/":                       handlers.GetPostByURL,
		"/posts/url-is-available/":      handlers.CheckIfURLIsAvailable,
		"/users/username-is-available/": handlers.CheckIfUsernameIsAvailable,
	}

	return routes
}

func protectedRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/protected": handlers.GetTokenDetails,

		// post routes
		"/posts/new":   handlers.GetNewestPosts,
		"/post/create": handlers.CreatePost,
		"/posts/top/":  handlers.GetTopPosts,
		"/posts":       handlers.GetPosts,

		// tag routes
		"/tags": handlers.GetTags,
	}

	return routes
}
