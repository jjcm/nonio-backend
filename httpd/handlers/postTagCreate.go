package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"soci-backend/httpd/utils"
	"soci-backend/models"
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
		tempTag, err := models.TagFactory(tagName, user)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Create tag: %v", err)
		}

		tag = tempTag
	}

	return &user, &post, &tag, nil
}

// CreatePostTag - protected http handler
// the user associated with the passed auth token can create a new post-tag
func CreatePostTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("You can only POST to the CreatePostTag route"), 405)
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

	postTag := &models.PostTag{}
	// check if the PostTag is existed in database
	if err := postTag.FindByUK(post.ID, tag.ID); err != nil {
		sendSystemError(w, fmt.Errorf("Query PostTag: %v", err))
		return
	}
	// if the PostTag is existed, return error
	if postTag.PostID > 0 {
		sendSystemError(w, fmt.Errorf("PostTag exists"))
		return
	}

	postTagVote := &models.PostTagVote{}
	// check if this is the first PostTagVote by user for the specific post
	votes, err := postTagVote.GetVotesByPostUser(post.ID, user.ID)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Query votes: %v", err))
		return
	}
	needUpdatePost := true
	if len(votes) > 0 {
		needUpdatePost = false
	}

	// prepare the PostTagVote for insertion
	postTagVote = &models.PostTagVote{
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
		postTag.PostID = post.ID
		postTag.TagID = tag.ID
		if err := postTag.CreatePostTagWithTx(tx); err != nil {
			return fmt.Errorf("Create PostTag: %v", err)
		}

		// insert the PostTagVote to database
		if err := postTagVote.CreatePostTagVoteWithTx(tx); err != nil {
			return fmt.Errorf("Create PostTagVote: %v", err)
		}

		// check if it needs to increment the score of post
		if needUpdatePost {
			// increment the score of Post
			if err := post.IncrementScoreWithTx(tx, post.ID); err != nil {
				return fmt.Errorf("Increment Post's score: %v", err)
			}
		}

		return nil
	}); err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, postTagVote, 200)
}
