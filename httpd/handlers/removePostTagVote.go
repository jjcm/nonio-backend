package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jjcm/soci-backend/models"
)

// PostTagVoteRemoveRequest this is the shape of the JSON request that is needed to
// create a vote for post tag
type PostTagVoteRemoveRequest struct {
	PostURL string `json:"post"`
	TagName string `json:"tag"`
}

// RemovePostTagVote - protected http handler
// the user associated with the passed auth token can create a new post-tag
func RemovePostTagVote(w http.ResponseWriter, r *http.Request) {
	// any non GET handlers need to attach CORS headers. I always forget about that
	CorsAdjustments(&w)

	// silly AJAX prflight, here's where we can put in the CORS requirements
	if r.Method == "OPTIONS" {
		SendResponse(w, "", 200)
		return
	}

	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to RemovePostTagVote route"), 405)
		return
	}

	// decode the request parameters 'post_id' and 'tag_id'
	var request PostTagVoteRemoveRequest
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

	// check if there is PostTagVote for the user id, post id and tag id
	postTagVote := &models.PostTagVote{}
	if err := postTagVote.FindByUK(post.ID, tag.ID, user.ID); err != nil {
		sendSystemError(w, fmt.Errorf("Query post-tag-vote: %v", err))
		return
	}
	// if there is not PostTagVote, just return directly
	if postTagVote.ID == 0 {
		SendResponse(w, true, 200)
		return
	}

	// query the votes with post id and tag id
	votes, err := postTagVote.GetVotesByPostTag(post.ID, tag.ID)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Query votes: %v", err))
		return
	}
	// only there is the user's PostTagVote, so it needs to delete the PostTag
	needDelPostTag := false
	if len(votes) == 1 {
		needDelPostTag = true
	}

	// check if this is the only one PostTagVote by user for the specific post
	votes, err = postTagVote.GetVotesByPostUser(post.ID, user.ID)
	if err != nil {
		sendSystemError(w, fmt.Errorf("Query votes: %v", err))
		return
	}
	needUpdatePost := false
	if len(votes) == 1 {
		needUpdatePost = true
	}

	// do many database operations with transaction
	if err = models.WithTransaction(func(tx models.Transaction) error {
		// delete the PostTagVote with unique key: user id, post id, tag id
		if err := postTagVote.DeleteByUKWithTx(tx, post.ID, tag.ID, user.ID); err != nil {
			return fmt.Errorf("Delete post-tag-vote: %v", err)
		}

		postTag := &models.PostTag{}
		if needDelPostTag {
			// delete the PostTag with unique key: post id, tag id
			if err := postTag.DeleteByUKWithTx(tx, post.ID, tag.ID); err != nil {
				return fmt.Errorf("Delete post-tag: %v", err)
			}
		} else {
			// decrement the score of the PostTag with unique key: post id, tag id
			if err := postTag.DecrementScoreWithTx(tx, post.ID, tag.ID); err != nil {
				return fmt.Errorf("Delete PostTag's score: %v", err)
			}
		}

		// if it needs to decrement the score of the Post with post id
		if needUpdatePost {
			if err := post.DecrementScoreWithTx(tx, post.ID); err != nil {
				return fmt.Errorf("Decrement Post's score: %v", err)
			}
		}

		return nil
	}); err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, true, 200)
}
