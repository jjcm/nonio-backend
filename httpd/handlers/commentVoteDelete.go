package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soci-backend/models"
)

// PostCommentVoteDeletionRequest defines the parameters for deleting the comment vote
type PostCommentVoteDeletionRequest struct {
	ID int `json:"id"`
}

func decrementLineageScore(tx models.Transaction, id int) (parent int, err error) {
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

	// decrement lineage source for the comment
	if err = comment.DecrementLineageScoreWithTx(tx, id); err != nil {
		return
	}

	return comment.ParentID, nil
}

// RemoveCommentVote - protected http handler
func RemoveCommentVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, MakeError("You can only POST to RemoveCommentVote route"), 405)
		return
	}

	// decode the request parameters
	var request PostCommentVoteDeletionRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)

	// do many database operations with transaction
	if err := models.WithTransaction(func(tx models.Transaction) error {

		id := request.ID
		// decrement the lineage_score for the comment, until the parent is zero
		for {
			parent, err := decrementLineageScore(tx, id)
			if err != nil {
				return fmt.Errorf("decrement lineage score: %v", err)
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
