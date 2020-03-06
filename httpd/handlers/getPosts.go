package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jjcm/soci-backend/models"
)

// GetPosts - get the posts from database with different url parameters
func GetPosts(w http.ResponseWriter, r *http.Request) {
	// parse the url parameters
	r.ParseForm()

	// get the user id from context
	userID := r.Context().Value("user_id").(int)
	// query the user by user id
	user := &models.User{}
	if err := user.FindByID(userID); err != nil {
		sendSystemError(w, fmt.Errorf("Query user: %v", err))
		return
	}

	offset := 0
	// check for offset
	formOffset := strings.TrimSpace(r.FormValue("offset"))
	if formOffset != "" {
		var err error
		offset, err = strconv.Atoi(formOffset)
		if err != nil {
			sendSystemError(w, fmt.Errorf("String to int: %v", err))
			return
		}
	}

	// ?time=all|day|week|month|year
	// Only returns posts that were created within a specific time period.
	// - all DEFAULT
	//   No time constraints.
	// - day
	//   Only returns posts created in the last day.
	// - week
	//   Only returns posts created in the last week.
	// - month
	//   Only returns posts created in the last month.
	// - year
	//   Only returns posts created in the last year.
	//
	// parse the time value
	var cutoff time.Time
	switch formTime := strings.TrimSpace(r.FormValue("time")); formTime {
	case "day":
		cutoff = time.Now().AddDate(0, 0, -1)
	case "week":
		cutoff = time.Now().AddDate(0, 0, -7)
	case "month":
		cutoff = time.Now().AddDate(0, -1, 0)
	case "year":
		cutoff = time.Now().AddDate(-1, 0, 0)
	case "all":
		cutoff = time.Now().AddDate(-50, 0, 0)
	case "": // if it's empty, default cutoff
		// check duration of 24 hours vs last login
		cutoff = user.LastLogin
		oneDayAgo := time.Now().AddDate(0, 0, -1)
		if user.LastLogin.After(oneDayAgo) {
			cutoff = oneDayAgo
		}
	default:
		sendSystemError(w, fmt.Errorf("Invalid field 'time': %v", formTime))
		return
	}

	// ?sort=popular|top|new
	// Returns posts sorted by a particular algorithm.
	// - popular DEFAULT
	//   Posts are sorted by score, within a time span between now and the requesting user's last login, or 24 hours, whichever is longer.
	// - new
	//   Posts are sorted by date. Newest first.
	// - top
	//   Posts are sorted by score. Highest first.
	sort := strings.TrimSpace(r.FormValue("sort"))
	if sort != "" {
		var posts []models.Post

		var err error
		switch sort {
		case "popular":
			// parse the tag name
			formTag := strings.TrimSpace(r.FormValue("tag"))
			if formTag == "" {
				sendSystemError(w, errors.New("Tag is empty"))
				return
			}
			// query the posts by tag name
			posts, err = models.GetPostsByPostTagScoreSince(formTag, cutoff)

		case "new":
			// query the posts by the offset
			posts, err = models.GetLatestPosts(offset)

		case "top":
			// query the posts by the cutoff and offset
			posts, err = models.GetPostsByScoreSince(cutoff, offset)

		default:
			sendSystemError(w, fmt.Errorf("Invalid field 'sort': %v", sort))
			return
		}
		if err != nil {
			sendSystemError(w, fmt.Errorf("Query posts by sort=%s: %v", sort, err))
			return
		}

		// send the result back to the clients
		output := map[string]interface{}{
			"posts": posts,
		}
		SendResponse(w, output, 200)
		return
	}

	// ?offset=NUMBER
	// Offsets the responses by a set number.
	posts, err := models.GetPostsByScoreSince(cutoff, offset)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Query posts by offset: %v", err))
		return
	}

	output := map[string]interface{}{
		"posts": posts,
	}
	SendResponse(w, output, 200)
	return
}
