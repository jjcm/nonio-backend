package httpd

import (
	"net/http"

	"soci-backend/httpd/handlers"
)

// OpenRoutes - routes that don't require auth
func OpenRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/": handlers.Home,

		// these should be user routes
		"/register": handlers.Register,
		"/login":    handlers.Login,
		//"/login-social":          handlers.LoginSocial,
		//"/login-social/callback": handlers.LoginSocialCallback,

		"/posts/":                 handlers.GetPostByURL,
		"/post/url-is-available/": handlers.CheckIfURLIsAvailable,

		"/user/username-is-available/":  handlers.CheckIfUsernameIsAvailable,
		"/user/forgot-password-request": handlers.ForgotPasswordRequest,
		"/users/":                       handlers.GetUser,

		"/comments": handlers.GetComments,
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
		"/tags":  handlers.GetTags,
		"/tags/": handlers.GetTagsByPrefix,

		// COMMENT ROUTES
		"/comment/create":      handlers.CommentOnPost,
		"/comment/edit":        handlers.EditComment,
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

		// USER ROUTES
		"/user/change-password":    handlers.ChangePassword,
		"/user/update-description": handlers.UpdateDescription,
		"/user/get-financials":     handlers.GetFinancials,

		// VOTES ROUTES
		"/votes": handlers.GetVotes,

		// COMMENT VOTES ROUTES
		"/comment-votes": handlers.GetCommentVotes,
		//"/comment-votes/post/": handlers.GetCommentVotesForPost,
		// TODO "/comment-votes/user/":               handlers.GetCommentVotesForUser,
	}

	return routes
}
