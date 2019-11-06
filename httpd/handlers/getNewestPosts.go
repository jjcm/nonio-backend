package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jjcm/soci-backend/models"
)

// GetNewestPosts - get 100 of the latest posts and return them as JSON
func GetNewestPosts(w http.ResponseWriter, r *http.Request) {
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
	posts, err := models.GetLatestPosts(offset)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	output := map[string]interface{}{
		"posts": posts,
	}
	SendResponse(w, output, 200)
}
