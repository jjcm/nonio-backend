package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jjcm/soci-backend/models"
)

// PostCreationRequest this is the shape of the JSON request that is needed to
// create a new post
type PostCreationRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

// CreatePost - protected http handler
// the user associated with the passed auth token can create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	// any non GET handlers need to attach CORS headers. I always forget about that
	CorsAdjustments(&w)

	// silly AJAX prflight, here's where we can put in the CORS requirements
	if r.Method == "OPTIONS" {
		SendResponse(w, "", 200)
		return
	}

	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to the post creation route"), 405)
		return
	}

	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	var postRequest PostCreationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&postRequest)
	newPost, err := u.CreatePost(postRequest.Title, postRequest.Content, postRequest.Type)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, &newPost, 200)
}
