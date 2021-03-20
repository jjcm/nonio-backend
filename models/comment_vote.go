package models

import (
	"database/sql"
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
func (u *User) CreateCommentVote(commentID int, upvote bool) (CommentVote, error) {
	vote := CommentVote{}
	// Check if a vote already exists
	vote.FindByUK(commentID, u.ID)
	if vote.ID != 0 {
		if vote.Upvote == upvote {
			// it's the same type of vote, so we don't need to do anything
			return vote, nil
		} else {
			// the vote is different, so lets remove the other vote first then add the new one
			u.DeleteCommentVote(commentID)
			return vote, nil
		}
	} else {
		// The comment vote doesn't yet exist
	}

	// TODO
	return vote, nil
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

	// run the correct sql query
	var query = "SELECT * FROM comment_votes WHERE voter_id = ? AND post_id = ?"
	err := DBConn.Select(&votes, query, u.ID, postID)
	if err != nil {
		return votes, err
	}

	return votes, nil
}

/************************************************/
/******************** UPDATE ********************/
/************************************************/
// TODO

/************************************************/
/******************** DELETE ********************/
/************************************************/

// RemoveVote removes a user's vote on a comment, if it exists. It will change the lineage score for the comment and its ancestors.
func (u *User) DeleteCommentVote(commentID int) error {
	return nil
}

// DeleteByUKWithTx - delete a PostTagVote in the database by unique keys
func (v *CommentVote) DeleteByUKWithTx(tx Transaction, commentID int, userID int) error {
	_, err := tx.Exec("delete from posts_tags_votes where post_id = ? and tag_id = ? and voter_id = ?", commentID, userID)
	return err
}
