package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// GetCommentsForPost will return all comments for a specific post
func GetCommentsForPost(w http.ResponseWriter, r *http.Request) {
	postSlug := strings.ToLower(utils.ParseRouteParameter(r.URL.Path, "/comments/post/"))
	p := models.Post{}
	p.FindByURL(postSlug)
	if p.ID == 0 {
		sendNotFound(w, errors.New("Post with url '"+postSlug+"' not found"))
		return
	}

	// query the comments for the post order by lineage score
	comments, err := models.GetCommentsByPost(p.ID)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	output := map[string]interface{}{
		"comments": comments,
	}
	SendResponse(w, output, 200)
}

// GetCommentsForUser will return all comments for a specific user
func GetCommentsForUser(w http.ResponseWriter, r *http.Request) {
	username := strings.ToLower(utils.ParseRouteParameter(r.URL.Path, "/comments/user/"))
	u := models.User{}
	u.FindByUsername(username)
	if u.ID == 0 {
		sendNotFound(w, errors.New("User with username '"+username+"' not found"))
		return
	}

	// query the comments for the post order by lineage score
	comments, err := u.GetComments()
	if err != nil {
		sendSystemError(w, err)
		return
	}

	output := map[string]interface{}{
		"comments": comments,
	}
	SendResponse(w, output, 200)
}

// GetComments - get the comments from database with different url parameters
func GetComments(w http.ResponseWriter, r *http.Request) {
	Log.Info(r.URL)
	params := &models.CommentQueryParams{}
	// parse the url parameters
	r.ParseForm()

	// ?offset=NUMBER
	// Offsets the responses by a set number.
	formOffset := strings.TrimSpace(r.FormValue("offset"))
	if formOffset != "" {
		offset, err := strconv.Atoi(formOffset)
		if err != nil {
			sendSystemError(w, fmt.Errorf("string to int: %v", err))
			return
		}
		params.Offset = offset
	}

	// ?time=all|day|week|month|year
	// Only returns comments that were created within a specific time period.
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
	default: // default time
		cutoff = time.Now().AddDate(-50, 0, 0)
	}
	params.Since = cutoff.Format("2006-01-02 15:04:05")

	// ?post=POST
	// Only returns results that match a specific post url.
	post := strings.TrimSpace(r.FormValue("post"))
	if post != "" {
		p := &models.Post{}
		err := p.FindByURL(post)
		if err != nil {
			sendSystemError(w, fmt.Errorf("query comments by post %s: %v", post, err))
			return
		}
		params.PostID = p.ID
	}

	// ?sort=popular|top|new
	// Returns comments sorted by a particular algorithm.
	formSort := strings.TrimSpace(r.FormValue("sort"))
	params.Sort = "top"
	if formSort == "new" {
		params.Sort = "new"
	}

	// ?user=USER
	// Only returns results comments that were made by a specific user.
	formUser := strings.TrimSpace(r.FormValue("user"))
	if formUser != "" {
		author := models.User{}
		// query the user by user name
		if err := author.FindByUsername(formUser); err != nil {
			sendSystemError(w, fmt.Errorf("query user by name %s: %v", formUser, err))
			return
		}
		if author.ID == 0 {
			sendNotFound(w, errors.New("user's name: "+formUser))
			return
		}
		params.UserID = author.ID
	}

	// query the comments by the url parameters
	comments, err := models.GetCommentsByParams(params)
	if err != nil {
		sendSystemError(w, fmt.Errorf("query comments by parameters: %v", err))
		return
	}

	output := map[string]interface{}{
		"comments": comments,
	}
	SendResponse(w, output, 200)
}
