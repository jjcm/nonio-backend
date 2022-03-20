package handlers

import (
	"encoding/json"
	"net/http"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// CreatePost - protected http handler
// the user associated with the passed auth token can create a new post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the post creation route"), 405)
		return
	}

	type requestPayload struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&payload)

	u := models.User{}
	u.FindByID(r.Context().Value("user_id").(int))

	if u.AccountType != "supporter" {
		SendResponse(w, utils.MakeError("only supporters can submit posts"), 403)
		return
	}

	newPost, err := u.CreatePost(payload.Title, payload.URL, payload.Content, payload.Type, payload.Width, payload.Height)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, &newPost, 200)

	// Nuke the cache
	PostCache = make(map[string]PostQueryResponse)
}
