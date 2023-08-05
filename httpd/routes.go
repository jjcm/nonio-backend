package httpd

import (
	"net/http"

	"soci-backend/httpd/handlers"
)

// OpenRoutes - routes that don't require auth
func OpenRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/": handlers.Home,

		"/posts":                  handlers.GetPosts,
		"/posts/":                 handlers.GetPostByURL,
		"/post/url-is-available/": handlers.CheckIfURLIsAvailable,

		"/user/register":                  handlers.Register,
		"/user/login":                     handlers.Login,
		"/user/username-is-available/":    handlers.CheckIfUsernameIsAvailable,
		"/user/forgot-password-request":   handlers.ForgotPasswordRequest,
		"/user/change-forgotten-password": handlers.ChangeForgottenPassword,
		"/users/":                         handlers.GetUser,

		"/comments": handlers.GetComments,

		"/stripe/webhooks": handlers.StripeWebhook,
	}

	return routes
}

// ProtectedRoutes - routes that require auth
func ProtectedRoutes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := map[string]func(http.ResponseWriter, *http.Request){
		"/protected": handlers.GetTokenDetails,

		// POST ROUTES
		"/post/create": handlers.CreatePost,
		// TODO "/post/delete/": handlers.DeletePost,
		"/post/parse-external-url": handlers.CheckExternalURLTitle,

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
		"/user/change-password":      handlers.ChangePassword,
		"/user/update-description":   handlers.UpdateDescription,
		"/user/get-financials":       handlers.GetFinancials,
		"/user/get-financial-ledger": handlers.GetFinancialLedger,
		// TODO - set up the GetSettings route or something similar to return whether the user is a subscriber or not
		// "/user/get-settings":        handlers.GetSettings,
		"/user/choose-free-account": handlers.ChooseFreeAccount,

		// NOTIFICATIONS ROUTES
		"/notifications":              handlers.GetNotifications,
		"/notifications/unread-count": handlers.GetUnreadNotificationCount,
		"/notification/mark-read":     handlers.MarkNotificationRead,

		// VOTES ROUTES
		"/votes": handlers.GetVotes,

		// COMMENT VOTES ROUTES
		"/comment-votes": handlers.GetCommentVotes,
		//"/comment-votes/post/": handlers.GetCommentVotesForPost,
		// TODO "/comment-votes/user/":               handlers.GetCommentVotesForUser,

		"/stripe/subscription/create": handlers.StripeCreateSubscription,
		"/stripe/subscription/delete": handlers.StripeCancelSubscription,
		"/stripe/subscription/edit":   handlers.StripeEditSubscription,
		"/stripe/subscription":        handlers.StripeGetSubscription,
		"/stripe/price-config":        handlers.StripeGetPriceConfig,
		"/stripe/create-customer":     handlers.StripeCreateCustomer,
		"/stripe/get-connect-link":    handlers.GetConnectLink,

		// user-ban
		"/user/ban": handlers.UserBan,
	}

	return routes
}
