package httpd

import (
	"net/http"

	"soci-backend/httpd/handlers"
)

// OpenRoutes - routes that don't require auth
func OpenRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/":                      handlers.Home,
		"/register":              handlers.Register,
		"/login":                 handlers.Login,
		"/login-social":          handlers.LoginSocial,
		"/login-social/callback": handlers.LoginSocialCallback,
		"/info":                  handlers.Info,
		"/posts/":                handlers.GetPostByURL,

		//change this to /post/url-is-available
		"/posts/url-is-available/": handlers.CheckIfURLIsAvailable,

		//change this to /user/username-is-available
		"/users/username-is-available/": handlers.CheckIfUsernameIsAvailable,

		"/comments/post/": handlers.GetCommentsForPost,
		// TODO "/comments/user/":               handlers.GetCommentsForUser,
		// TODO "/comments/comment/":            handlers.GetCommentsForComment,
	}

	return routes
}

// ProtectedRoutes - routes that require auth
func ProtectedRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/protected": handlers.GetTokenDetails,

		// POST ROUTES
		"/posts":       handlers.GetPosts,
		"/post/create": handlers.CreatePost,
		// TODO "/post/delete/": handlers.DeletePost,

		// TAG ROUTES
		"/tags": handlers.GetTags,

		// COMMENT ROUTES
		"/comment/create":      handlers.CommentOnPost,
		"/comment/delete":      handlers.DeleteComment,
		"/comment/abandon":     handlers.AbandonComment,
		"/comment/add-vote":    handlers.AddCommentVote,
		"/comment/remove-vote": handlers.RemoveCommentVote,

		// POSTTAG ROUTES
		"/posttag/create":      handlers.CreatePostTag,
		"/posttag/add-vote":    handlers.AddPostTagVote,
		"/posttag/remove-vote": handlers.RemovePostTagVote,

		// SUBSCRIPTION ROUTES
		"/subscriptions":       handlers.GetSubscriptions,
		"/subscription/create": handlers.CreateSubscription,
		"/subscription/delete": handlers.DeleteSubscription,

		// VOTES ROUTES
		"/votes": handlers.GetVotes,
	}

	return routes
}
