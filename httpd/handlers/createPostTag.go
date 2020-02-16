package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jjcm/soci-backend/models"
)

type PostTageVote struct {
}

// PostTagCreationRequest this is the shape of the JSON request that is needed to
// create a new post
type PostTagCreationRequest struct {
	PostId int `json:"post_id"`
	TagId  int `json:"tag_id"`
}

// CreatePostTag - protected http handler
// the user associated with the passed auth token can create a new post-tag
func CreatePostTag(w http.ResponseWriter, r *http.Request) {
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

	// decode the request parameters 'post_id' and 'tag_id'
	var request PostTagCreationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)

	// query the user by user id
	user := models.User{}
	if err := user.FindByID(r.Context().Value("user_id").(int)); err != nil {
		sendSystemError(w, fmt.Errorf("Query user: %v", err))
		return
	}

	// query the post by post id
	post := models.Post{}
	if err := post.FindByID(request.PostId); err != nil {
		sendSystemError(w, fmt.Errorf("Query post: %v", err))
		return
	}

	// query the tag by tag id
	tag := models.Tag{}
	if err := tag.FindByID(request.TagId); err != nil {
		sendSystemError(w, fmt.Errorf("Query tag: %v", err))
		return
	}

	postTag := models.PostTag{
		Post: &post,
		Tag:  &tag,
	}
	// check if the PostTag is existed in database
	item, err := postTag.FindByPostTagIds(request.PostId, request.TagId)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Query post-tag: %v", err))
		return
	}
	// if the PostTag is existed, return error
	if item != nil {
		sendSystemError(w, fmt.Errorf("PostTag is existed"))
		return
	}

	// insert the PostTag to database
	if err := postTag.CreatePostTag(); err != nil {
		sendSystemError(w, fmt.Errorf("Create PostTag: %v", err))
		return
	}

	// prepare the value for insertion
	postTagVote := &models.PostTagVote{
		Post:      &post,
		PostID:    post.ID,
		PostURL:   post.URL,
		Tag:       &tag,
		TagID:     tag.ID,
		TagName:   tag.Name,
		Voter:     &user,
		VoterID:   user.ID,
		VoterName: user.Name,
	}
	// insert the PostTagVote to database
	if err := postTagVote.CreatePostTagVote(); err != nil {
		sendSystemError(w, fmt.Errorf("Create PostTagVote: %v", err))
		return
	}

	SendResponse(w, postTagVote, 200)
}
