package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

type PostQueryResponse struct {
	Response  []byte
	CreatedAt time.Time
}

var PostCache map[string]PostQueryResponse = make(map[string]PostQueryResponse)

// GetPostByURL find a specific post in the database and send back a JSON
// representation of it
func GetPostByURL(w http.ResponseWriter, r *http.Request) {
	url := utils.ParseRouteParameter(r.URL.Path, "/posts/")
	if strings.TrimSpace(url) == "" {
		sendSystemError(w, errors.New("please pass a valid URL for us to get you your requested content"))
		return
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

// fill the tags for those posts
func fillPostTags(posts []*models.Post) error {
	for _, post := range posts {
		tags, err := models.GetPostTags(post.ID)
		if err != nil {
			return err
		}
		post.Tags = tags
	}
	return nil
}

// GetPosts - get the posts from database with different url parameters
func GetPosts(w http.ResponseWriter, r *http.Request) {

	// check our cache first
	var cacheResponse PostQueryResponse
	// Use if we ever want to invalidate the cache after a period of time.
	//if cacheResponse, ok := PostCache[r.URL.String()]; ok && time.Now().Add(time.Minute*-1).Before(cacheResponse.CreatedAt) {
	if cacheResponse, ok := PostCache[r.URL.String()]; ok {
		Log.Info("cache hit")
		Log.Info(r.URL)
		SendJSONResponse(w, cacheResponse.Response, 200)
		return
	}

	Log.Info(r.URL)
	params := &models.PostQueryParams{}
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
	// Only returns posts that were created within a specific time period.
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

	// ?tag=TAG
	// Only returns results that match a specific tag. Multiple tags can be listed by separating tags with a +
	tag := strings.TrimSpace(r.FormValue("tag"))
	if tag != "" {
		t := &models.Tag{}
		err := t.FindByTagName(tag)
		if err != nil {
			sendSystemError(w, fmt.Errorf("query posts by tag %s: %v", tag, err))
			return
		}
		params.TagID = t.ID
	}

	// ?sort=popular|top|new
	// Returns posts sorted by a particular algorithm.
	formSort := strings.TrimSpace(r.FormValue("sort"))
	Log.Info("formSort: %s", formSort)
	params.Sort = "popular"
	switch formSort {
	case "popular":
		// This is unused for now, but will eventually be to show "popular things since the users last login"
		// get the user id from context
		userID := r.Context().Value("user_id").(int)
		// query the user by user id
		user := &models.User{}
		if err := user.FindByID(userID); err != nil {
			sendSystemError(w, fmt.Errorf("query user: %v", err))
			return
		}

		// check duration of 24 hours vs last login
		cutoff = user.LastLogin
		oneDayAgo := time.Now().AddDate(0, 0, -1)
		if user.LastLogin.After(oneDayAgo) {
			cutoff = oneDayAgo
		}
		params.Since = cutoff.Format("2006-01-02 15:04:05")
		params.Sort = "popular"

	case "new":
		// sort by the create time
		params.Sort = "new"
	case "top":
		params.Sort = "top"
	default:
		params.Sort = "popular"
	}

	// ?user=USER
	// Only returns results posts that were made by a specific user.
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

	// query the posts by the url parameters
	posts, err := models.GetPostsByParams(params)
	if err != nil {
		sendSystemError(w, fmt.Errorf("query posts by parameters: %v", err))
		return
	}

	// fill the tags for the posts
	if err := fillPostTags(posts); err != nil {
		sendSystemError(w, fmt.Errorf("query tags by posts: %v", err))
		return
	}

	output := map[string]interface{}{
		"posts": posts,
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		SendResponse(w, err.Error(), 500)
		return
	}
	SendJSONResponse(w, jsonData, 200)

	// Add the query to our cache
	cacheResponse.Response = jsonData
	cacheResponse.CreatedAt = time.Now()
	PostCache[r.URL.String()] = cacheResponse
}
