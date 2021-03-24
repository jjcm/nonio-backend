package handlers

import (
	"net/http"

	"errors"
	"soci-backend/httpd/utils"
	"soci-backend/models"
	"strings"
)

// GetCommentVotesForPost - gets a user's votes for a post.
func GetCommentVotesForPost(w http.ResponseWriter, r *http.Request) {
	url := utils.ParseRouteParameter(r.URL.Path, "/comment-votes/post/")
	if strings.TrimSpace(url) == "" {
		sendSystemError(w, errors.New("a post url is needed to get the comment votes for it"))
		return
	}

	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	p := models.Post{}
	err := p.FindByURL(url)
	if err != nil {
		sendNotFound(w, errors.New("we couldn't find a post with the url `"+url+"`"))
		return
	}

	votes, err := u.GetCommentVotesForPost(p.ID)
	if err != nil {
		SendResponse(w, utils.MakeError(err.Error()), 500)
		return
	}

	output := map[string]interface{}{
		"votes": votes,
	}
	SendResponse(w, output, 200)
}
