package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"soci-backend/models"
)

// DeleteComment will delete the comment matching the ID submitted
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		ID *int `json:"id"`
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
	if payload.ID == nil {
		sendSystemError(w, errors.New("Deleting a comment requires the `id` of the comment to be present"))
	}

	// second, find the user that is trying to write the post
	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	// are we commenting on a comment, or directly on the post itself?
	comment := models.Comment{}
	comment.FindByID(*(payload.ID))

	/*
		comment, err := u.CommentOnPost(post, &parentComment, payload.Content)
		if err != nil {
			sendSystemError(w, fmt.Errorf("Create comment: %v", err))
			return
		}
	*/

	// status 201 for "created"
	SendResponse(w, &comment, 201)
}
