package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jjcm/soci-backend/models"
)

// GetPosts - get all the posts in the system
func GetPosts(w http.ResponseWriter, r *http.Request) {
	// this code is what will help us set up the correct SQL Query
	posts, err := models.GetPostsFromQueryParams(r.URL.Query())

	output := map[string]interface{}{
		"posts": posts,
	}
	SendResponse(w, output, 200)
	return

	// CODE BELOW THIS LINE KEPT FOR REFERRENCE

	// check duration of 24 hours vs. last login
	user := models.User{}
	user.FindByID(r.Context().Value("user_id").(int))

	lastLogin := user.LastLogin
	twenty24HoursAgo := time.Now().AddDate(0, 0, -1)
	cutoff := lastLogin
	if lastLogin.After(twenty24HoursAgo) {
		cutoff = twenty24HoursAgo
	}

	// check for offset
	r.ParseForm()
	formOffset := r.FormValue("offset")
	offset := 0
	if strings.TrimSpace(formOffset) != "" {
		var err error
		offset, err = strconv.Atoi(formOffset)
		if err != nil {
			sendSystemError(w, err)
			return
		}
	}
	posts, err = models.GetPostsByScoreSince(cutoff, offset)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	o := map[string]interface{}{
		"posts": posts,
	}
	SendResponse(w, o, 200)
}
