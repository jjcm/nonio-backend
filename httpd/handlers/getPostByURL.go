package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jjcm/soci-backend/models"
)

// GetPostByURL find a specific post in the database and send back a JSON
// representation of it
func GetPostByURL(w http.ResponseWriter, r *http.Request) {
	url := parseRouteParamater(r.URL.Path, "/posts/")
	if strings.TrimSpace(url) == "" {
		sendSystemError(w, errors.New("Please pass a valid URL for us to get you your requested content"))
	}

	p := models.Post{}
	err := p.FindByURL(url)
	if err != nil {
		sendNotFound(w, errors.New("we couldn't find a post with the url `"+url+"`"))
		return
	}

	// pass a pointer to the post so that it runs through the custom
	// JSON marshaler
	SendResponse(w, &p, 200)
}
