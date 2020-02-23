package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jjcm/soci-backend/models"
)

// GetPopularPosts - get posts from the database ordered by score and return them as
// JSON. Optionally limit them to a certian timeframe
func GetPopularPosts(w http.ResponseWriter, r *http.Request) {
	// parse the tag name from URL path
	nameParam := strings.ToLower(parseRouteParameter(r.URL.Path, "/tags/popular/"))
	name := strings.ReplaceAll(nameParam, "/", "")

	// get the user id from context
	userID := r.Context().Value("user_id").(int)

	// query the user by user id
	user := &models.User{}
	if err := user.FindByID(userID); err != nil {
		sendSystemError(w, fmt.Errorf("Query user: %v", err))
		return
	}

	cutoff := user.LastLogin
	oneDayAgo := time.Now().AddDate(0, 0, -1)
	if user.LastLogin.After(oneDayAgo) {
		cutoff = oneDayAgo
	}

	posts, err := models.GetPostsByPostTagScoreSince(name, cutoff)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Get popular posts: %v", err))
		return
	}

	output := map[string]interface{}{
		"posts": posts,
	}
	SendResponse(w, output, 200)
}
