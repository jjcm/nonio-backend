package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jjcm/soci-backend/models"
)

// PostTagCreationRequest this is the shape of the JSON request that is needed to
// create a new tag for post
type PostTagCreationRequest struct {
	PostURL string `json:"post"`
	TagName string `json:"tag"`
}

// find the structure of user, post, tag with user id, post url, tag name from database
func findUserPostTag(userID int, postURL string, tagName string) (*models.User, *models.Post, *models.Tag, error) {
	// query the user by user id
	user := models.User{}
	if err := user.FindByID(userID); err != nil {
		return nil, nil, nil, fmt.Errorf("Query user: %v", err)
	}

	// query the post by post id
	post := models.Post{}
	if err := post.FindByURL(postURL); err != nil {
		return nil, nil, nil, fmt.Errorf("Query post: %v", err)
	}

	// query the tag by tag id
	tag := models.Tag{}
	if err := tag.FindByTagName(tagName); err != nil {
		return nil, nil, nil, fmt.Errorf("Query tag: %v", err)
	}
	// if there is no rows about the tag name, insert a new one
	if tag.ID == 0 {
		id, err := models.CreateTag(tagName, user)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Create tag: %v", err)
		}
		// update the Tag structure
		tag.ID = int(id)
		tag.Author = user
		tag.Name = tagName
	}

	return &user, &post, &tag, nil
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

	// decode the request parameters 'post' and 'tag'
	var request PostTagCreationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)

	// get the user id from context
	userID := r.Context().Value("user_id").(int)

	// find the structure of user, post, tag with user id, post url and tag name
	user, post, tag, err := findUserPostTag(userID, request.PostURL, request.TagName)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	postTag := models.PostTag{}
	// check if the PostTag is existed in database
	if err := postTag.FindByPostTagIds(post.ID, tag.ID); err != nil {
		sendSystemError(w, fmt.Errorf("Query post-tag: %v", err))
		return
	}
	// if the PostTag is existed, return error
	if postTag.PostID > 0 {
		sendSystemError(w, fmt.Errorf("PostTag is existed"))
		return
	}

	// prepare the PostTagVote for insertion
	postTagVote := &models.PostTagVote{
		Post:      post,
		PostID:    post.ID,
		PostURL:   post.URL,
		Tag:       tag,
		TagID:     tag.ID,
		TagName:   tag.Name,
		Voter:     user,
		VoterID:   user.ID,
		VoterName: user.Name,
	}

	// do many database operations with transaction
	if err = models.WithTransaction(func(tx models.Transaction) error {
		// insert the PostTag to database
		if err := postTag.CreatePostTag(); err != nil {
			return fmt.Errorf("Create PostTag: %v", err)
		}

		// insert the PostTagVote to database
		if err := postTagVote.CreatePostTagVote(); err != nil {
			return fmt.Errorf("Create PostTagVote: %v", err)
		}

		return nil
	}); err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, postTagVote, 200)
}
