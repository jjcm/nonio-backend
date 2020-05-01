package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"soci-backend/models"
)

// CommentOnPost will read in the JSON payload to add a comment to a given post
func CommentOnPost(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		PostURL  string `json:"post"`
		Content  string `json:"content"`
		Type     string `json:"type"`
		ParentID *int   `json:"parent"`
	}
	// any non GET handlers need to attach CORS headers. I always forget about that
	CorsAdjustments(&w)
	// silly AJAX prflight, here's where we can put in the CORS requirements
	if r.Method == "OPTIONS" {
		SendResponse(w, "", 200)
		return
	}
	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to the registration route"), 405)
		return
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	// before we even check for the existance of the related items, let's verify this comment payload is even valid
	if payload.Content == "" {
		sendSystemError(w, errors.New("Can not send us an empty comment. `content` key is required"))
	}
	if payload.Type == "" {
		payload.Type = "text" // sensible default
	}

	// first, find the post we are commenting on
	post := models.Post{}
	post.FindByURL(payload.PostURL)
	if post.ID == 0 {
		sendNotFound(w, errors.New("Post with URL '"+payload.PostURL+"' was not found"))
		return
	}

	// second, find the user that is trying to write the post
	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	// are we commenting on a comment, or directly on the post itself?
	parentComment := models.Comment{}
	if payload.ParentID != nil {
		parentComment.FindByID(*(payload.ParentID))
	}

	comment, err := u.CommentOnPost(post, &parentComment, payload.Type, payload.Text)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	// status 201 for "created"
	SendResponse(w, &comment, 201)
}
