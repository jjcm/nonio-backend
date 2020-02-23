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
		"/posts/":                       handlers.GetPostByURL,
		//"/posts/url/":                 handlers.GetPostByURL,
		"/posts/url-is-available/":      handlers.CheckIfURLIsAvailable,
		"/users/username-is-available/": handlers.CheckIfUsernameIsAvailable,
		"/comments/post/":               handlers.GetCommentsForPost,
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
		//"/posts/popular":       handlers.GetPosts,
		"/posts/user/": handlers.GetPostsByAuthor,

		// tag routes
		"/tags": handlers.GetTags,
		/*
		"/tags/popular/": handlers.GetPostsByTag,
		"/tags/top/day/": handlers.GetPostsByTag,
		"/tags/top/week/": handlers.GetPostsByTag,
		"/tags/top/month/": handlers.GetPostsByTag,
		"/tags/top/year/": handlers.GetPostsByTag,
		"/tags/top/all/": handlers.GetPostsByTag,
		"/tags/new/": handlers.GetPostsByTag,
		*/


		// comment routes
		"/comments/create": handlers.CommentOnPost,
		/*
			"/comments/add-vote": handlers.AddCommentVote,
			"/comments/remove-vote": handlers.RemoveCommentVote,
		*/

		// posttag routes
		"/posttags/create":   handlers.CreatePostTag,
		"/posttags/add-vote": handlers.AddPostTagVote,
		/*
			"/posttags/remove-vote": handlers.RemovePostTagVote,
		*/
	}

	return routes
}
