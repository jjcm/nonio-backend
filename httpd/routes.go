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
		"/posts": handlers.GetPosts,
		// DEPRECATE - MERGE IN WITH GETPOSTS "/posts/new":   handlers.GetNewestPosts,
		// DEPRECATE - MERGE IN WITH GETPOSTS "/posts/top/":  handlers.GetTopPosts,
		// DEPRECATE - MERGE IN WITH GETPOSTS "/posts/user/": handlers.GetPostsByAuthor,
		"/post/create/": handlers.CreatePost,
		// TODO "/post/view/": handlers.ViewPost,
		// TODO "/post/delete/": handlers.DeletePost,

		// TAG ROUTES
		"/tags": handlers.GetTags,
		// DEPRECATE - MERGE IN WITH GETPOSTS "/tags/popular/": handlers.GetPopularPosts,

		// COMMENT ROUTES
		// change to /comment/create
		"/comments/create": handlers.CommentOnPost,
		// TODO "/comment/add-vote": handlers.AddCommentVote,
		// TODO "/comment/remove-vote": handlers.RemoveCommentVote,

		// POSTTAG ROUTES
		//change this to /posttag/create
		"/posttags/create": handlers.CreatePostTag,
		//change this to /posttag/add-vote
		"/posttags/add-vote": handlers.AddPostTagVote,
		//change this to /posttag/remove-vote
		"/posttag/remove-vote": handlers.RemovePostTagVote,

		// VOTES ROUTES
		"/votes": handlers.GetVotes,
	}

	return routes
}
