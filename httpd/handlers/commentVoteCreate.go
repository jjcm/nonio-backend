package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// PostCommentVoteAdditionRequest defines the parameters for adding the comment vote
type PostCommentVoteAdditionRequest struct {
	ID     int  `json:"id"`
	Upvote bool `json:"upvoted"`
}

// AddCommentVote - protected http handler
func AddCommentVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("You can only POST to AddCommentVote route"), 405)
		return
	}

	Log.Info(r.Body)
	// decode the request parameters
	var request PostCommentVoteAdditionRequest
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&request)
	Log.Info("comment vote time")
	Log.Info(fmt.Sprintf("comment id: %v", request.ID))
	Log.Info(fmt.Sprintf("comment upvoted: %v", request.Upvote))

	// do many database operations with transaction
	if err := models.WithTransaction(func(tx models.Transaction) error {

		id := request.ID
		// increment the lineage_score for the comment, until the parent is zero
		for {
			comment := &models.Comment{}
			err := comment.FindByID(id)
			if err != nil {
				return err
			}

			if comment.ID == 0 {
				err = fmt.Errorf("Comment does not exist")
				return err
			}

			Log.Info(fmt.Sprintf("Comment %v is being upvoted? %v", id, request.Upvote))
			if request.Upvote {
				err = comment.IncrementLineageScoreWithTx(tx, id)
			} else {
				err = comment.DecrementLineageScoreWithTx(tx, id)
			}

			if err != nil {
				return fmt.Errorf("increment lineage score: %v", err)
			}

			if comment.ParentID == 0 {
				break
			}

			id = comment.ParentID
		}

		return nil
	}); err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, true, 200)
}
