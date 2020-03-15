package handlers

import (
	"errors"
	"net/http"
	"strings"

	"soci-backend/models"
)

// GetCommentsForPost will return all comments for a specific post
func GetCommentsForPost(w http.ResponseWriter, r *http.Request) {
	postSlug := strings.ToLower(parseRouteParameter(r.URL.Path, "/comments/post/"))
	p := models.Post{}
	p.FindByURL(postSlug)
	if p.ID == 0 {
		sendNotFound(w, errors.New("Post with url '"+postSlug+"' not found"))
		return
	}

	// spec say to get comments 5 levels deep to not get too far into a recursive mess
	comments, err := p.Comments(5)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	output := map[string]interface{}{
		"comments": models.StructureComments(comments),
	}
	SendResponse(w, output, 200)
}
