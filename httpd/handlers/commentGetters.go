package handlers

import (
	"errors"
	"net/http"
	"strings"

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
