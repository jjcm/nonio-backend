package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// AbandonComment will remove the user from the comment, but leave the content
func AbandonComment(w http.ResponseWriter, r *http.Request) {
	type requestPayload struct {
		ID *int `json:"id"`
	}
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("You can only POST to the abandon comment route"), 405)
		return
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)
	// before we even check for the existance of the related items, let's verify this comment payload is even valid
	if payload.ID == nil {
		sendSystemError(w, errors.New("Abandoning a comment requires the `id` of the comment to be present"))
		return
	}

	// second, find the user that is trying to write the post
	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	// make sure the comment we're deleting exists
	comment := models.Comment{}
	comment.FindByID(*(payload.ID))

	// make sure the owner of the comment is the user who's making the request
	if int(comment.AuthorID.Int32) != u.ID {
		SendResponse(w, utils.MakeError("You can only delete comments you own"), 401)
		return
	}

	err := u.AbandonComment(&comment)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Abandon comment: %v", err))
		return
	}

	// status 201 for "created"
	w.Header().Set("Access-Control-Allow-Origin", "*") // this should be locked down before launch
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(201)
	w.Write([]byte("true"))
}
