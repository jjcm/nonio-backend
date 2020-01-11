package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/jjcm/soci-backend/models"
)

// GetPostsByAuthor will return json to the response writer for all posts by the
// author with the matching username.
// i.e. /posts/user/lapubell
// will return 100 posts where the author's username is lapubell
// param offset will allow us to paginate the results
func GetPostsByAuthor(w http.ResponseWriter, r *http.Request) {
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

	authorUsername := strings.ToLower(parseRouteParameter(r.URL.Path, "/posts/user/"))
	user := models.User{}
	err := user.FindByUsername(authorUsername)
	if user.ID == 0 {
		sendNotFound(w, errors.New("We couldn't find a user with the username "+authorUsername))
		return
	}
	if err != nil {
		sendSystemError(w, err)
		return
	}

	posts, err := user.MyPosts(100, offset)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	output := map[string][]models.Post{
		"posts": posts,
	}
	SendResponse(w, output, 200)
}
