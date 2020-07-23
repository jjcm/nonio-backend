package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soci-backend/models"
)

// PostCommentVoteAdditionRequest defines the parameters for adding the comment vote
type PostCommentVoteAdditionRequest struct {
	ID int `json:"id"`
}

func incrementLineageScore(tx models.Transaction, id int) (parent int, err error) {
	// check if the comment is existed
	comment := &models.Comment{}
	if err = comment.FindByID(id); err != nil {
		return
	}
	// if the comment is not existed, return error
	if comment.ID == 0 {
		err = fmt.Errorf("Comment is not existed")
		return
	}

	// increment lineage source for the comment
	if err = comment.IncrementLineageScoreWithTx(tx, id); err != nil {
		return
	}

	return comment.ParentID, nil
}

// AddCommentVote - protected http handler
func AddCommentVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to AddCommentVote route"), 405)
		return
	}

	// decode the request parameters
	var request PostCommentVoteAdditionRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)

	// do many database operations with transaction
	if err := models.WithTransaction(func(tx models.Transaction) error {

		id := request.ID
		// increment the lineage_score for the comment, until the parent is zero
		for {
			parent, err := incrementLineageScore(tx, id)
			if err != nil {
				return fmt.Errorf("increment lineage score: %v", err)
			}
			if parent == 0 {
				break
			}
			id = parent
		}

		return nil
	}); err != nil {
		sendSystemError(w, err)
		return
	}

	// find the comment by id
	comment := &models.Comment{}
	comment.FindByID(request.ID)

	SendResponse(w, comment, 200)
}
