package models

import (
	"database/sql"
	"fmt"
	"time"
)

// CommentVote - struct representation of a comment vote
type CommentVote struct {
	ID        int       `db:"id" json:"-"`
	Comment   *Comment  `db:"-" json:"-"`
	CommentID int       `db:"comment_id" json:"comment_id"`
	Voter     *User     `db:"-" json:"-"`
	VoterName string    `db:"-" json:"-"`
	VoterID   int       `db:"voter_id" json:"-"`
	Upvote    bool      `db:"upvote" json:"-"`
	Post      *Post     `db:"-" json:"-"`
	PostID    int       `db:"post_id" json:"-"`
	PostURL   string    `db:"-" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

/************************************************/
/******************** CREATE ********************/
/************************************************/

// CreateCommentVote adds a vote for the user for a comment. It will change the lineage score for the comment and its ancestors.
func (u *User) CreateCommentVoteWithTx(tx Transaction, commentID int, upvote bool) error {
	vote := CommentVote{}
	// Check if a vote already exists
	vote.FindByUK(commentID, u.ID)
	if vote.ID != 0 {
		if vote.Upvote == upvote {
			// it's the same type of vote, so we don't need to do anything
			return nil
		} else {
			// the vote is different, so lets remove the other vote first then add the new one
			u.DeleteCommentVoteWithTx(tx, commentID)
			return u.CreateCommentVoteWithTx(tx, commentID, upvote)
		}
	} else {
		// The comment vote doesn't yet exist, so we need to create it.

		// Find the comment we're voting on
		comment := Comment{}

		if err := comment.FindByID(commentID); comment.ID == 0 {
			// Comment does not exist, so throw an error
			return fmt.Errorf("Comment does not exist: %v", err)
		}

		// Create the comment vote, return if there's an error
		if _, err := tx.Exec("INSERT INTO comment_votes (comment_id, voter_id, upvote, post_id) VALUES (?, ?, ?, ?)", commentID, u.ID, upvote, comment.PostID); err != nil {
			return err
		}

		if upvote {
			// The comment vote is an upvote
			if err := comment.AddUpvoteWithTx(tx); err != nil {
				return err
			}
		} else {
			// The comment vote is a downvote
			if err := comment.AddDownvoteWithTx(tx); err != nil {
				return err
			}
		}
	}

	return nil
}

/************************************************/
/********************* READ *********************/
/************************************************/

// FindByID - find a given CommentVote in the database by ID
func (v *CommentVote) FindByID(id int) error {
	dbCommentVote := CommentVote{}
	err := DBConn.Get(&dbCommentVote, "SELECT * FROM comment_votes WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	*v = dbCommentVote
	return nil
}

// FindByUK - find a given CommentVote in the database by unique keys
func (v *CommentVote) FindByUK(commentID int, userID int) error {
	dbCommentVote := CommentVote{}
	err := DBConn.Get(&dbCommentVote, "SELECT * FROM comment_votes WHERE voter_id = ? AND comment_id = ?", userID, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	*v = dbCommentVote
	return nil
}

// GetCommentVotesForPost will return every comment the user has voted on for a specific post
func (u *User) GetCommentVotesForPost(postID int) ([]CommentVote, error) {
	votes := []CommentVote{}

	err := DBConn.Select(&votes, "SELECT * FROM comment_votes WHERE voter_id = ? AND post_id = ?", u.ID, postID)
	if err != nil {
		return votes, err
	}

	return votes, nil
}

/************************************************/
/******************** UPDATE ********************/
/************************************************/
// Not needed

/************************************************/
/******************** DELETE ********************/
/************************************************/

// RemoveVote removes a user's vote on a comment, if it exists, via a transaction. It will change the lineage score for the comment and its ancestors.
func (u *User) DeleteCommentVoteWithTx(tx Transaction, commentID int) error {
	// Select the comment vote
	commentVote := CommentVote{}
	commentVote.FindByUK(commentID, u.ID)

	comment := Comment{}
	comment.FindByID(commentID)

	if commentVote.Upvote {
		// The comment vote is an upvote
		if err := comment.RemoveUpvoteWithTx(tx); err != nil {
			return err
		}
	} else {
		// The comment vote is a downvote
		if err := comment.RemoveDownvoteWithTx(tx); err != nil {
			return err
		}
	}

	// Then delete the comment vote
	_, err := tx.Exec("delete from comment_votes where comment_id = ? and voter_id = ?", commentID, u.ID)
	return err
}
