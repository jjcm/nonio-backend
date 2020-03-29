package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"soci-backend/models"
)

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
	params := &models.PostQueryParams{}
	// parse the url parameters
	r.ParseForm()

	// ?offset=NUMBER
	// Offsets the responses by a set number.
	formOffset := strings.TrimSpace(r.FormValue("offset"))
	if formOffset != "" {
		offset, err := strconv.Atoi(formOffset)
		if err != nil {
			sendSystemError(w, fmt.Errorf("String to int: %v", err))
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
	formTag := strings.TrimSpace(r.FormValue("tags"))
	// query the post ids from 'posts_tags' by 'tags', sorted by 'score' default
	if formTag != "" {
		tags := strings.Replace(strings.Trim(formTag, "+"), "+", ",", -1)

		pt := &models.PostTag{}
		ids, err := pt.GetPostsByTags(tags)
		if err != nil {
			sendSystemError(w, fmt.Errorf("Query posts by tags %s: %v", tags, err))
			return
		}
		params.TagIDs = ids
	}

	// sort by the post score default
	params.SortedByScore = true
	// if the tag list is not empty, sort by the score of the posttag
	if len(params.TagIDs) > 0 {
		params.SortedByScore = false
	}

	// ?sort=popular|top|new
	// Returns posts sorted by a particular algorithm.
	formSort := strings.TrimSpace(r.FormValue("sort"))
	switch formSort {
	case "popular":
		// get the user id from context
		userID := r.Context().Value("user_id").(int)
		// query the user by user id
		user := &models.User{}
		if err := user.FindByID(userID); err != nil {
			sendSystemError(w, fmt.Errorf("Query user: %v", err))
			return
		}

		// check duration of 24 hours vs last login
		cutoff = user.LastLogin
		oneDayAgo := time.Now().AddDate(0, 0, -1)
		if user.LastLogin.After(oneDayAgo) {
			cutoff = oneDayAgo
		}
		params.Since = cutoff.Format("2006-01-02 15:04:05")

	case "new":
		// sort by the create time
		params.SortedByScore = false
	}

	// ?user=USER
	// Only returns results posts that were made by a specific user.
	formUser := strings.TrimSpace(r.FormValue("user"))
	if formUser != "" {
		author := models.User{}
		// query the user by user name
		if err := author.FindByUsername(formUser); err != nil {
			sendSystemError(w, fmt.Errorf("Query user by name %s: %v", formUser, err))
			return
		}
		if author.ID == 0 {
			sendNotFound(w, errors.New("User's name: "+formUser))
			return
		}
		params.UserID = author.ID
	}

	// query the posts by the url parameters
	posts, err := models.GetPostsByParams(params)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Query posts by parameters: %v", err))
		return
	}

	// fill the tags for the posts
	if err := fillPostTags(posts); err != nil {
		sendSystemError(w, fmt.Errorf("Query tags by posts: %v", err))
		return
	}

	output := map[string]interface{}{
		"posts": posts,
	}
	SendResponse(w, output, 200)
	return
}
