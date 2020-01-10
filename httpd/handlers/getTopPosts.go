package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jjcm/soci-backend/models"
)

// GetTopPosts - get posts from the database ordered by score and return them as
// JSON. Optionally limit them to a certian timeframe
func GetTopPosts(w http.ResponseWriter, r *http.Request) {
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

	var cutoff time.Time
	timeframeParam := strings.ToLower(parseRouteParameter(r.URL.Path, "/posts/top/"))
	switch timeframeParam {
	case "day":
		cutoff = time.Now().AddDate(0, 0, -1)
	case "week":
		cutoff = time.Now().AddDate(0, 0, -7)
	case "month":
		cutoff = time.Now().AddDate(0, -1, 0)
	case "year":
		cutoff = time.Now().AddDate(-1, 0, 0)
	}

	posts, err := models.GetPostsByScoreSince(cutoff, offset)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	output := map[string]interface{}{
		"timeframe": timeframeParam,
		"posts":     posts,
	}
	SendResponse(w, output, 200)
}
