package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

// PostTagVote - struct representation of a single post-tag-vote
type PostTagVote struct {
	ID        int       `db:"id" json:"-"`
	Post      *Post     `db:"-" json:"-"`
	PostID    int       `db:"post_id" json:"postID"`
	PostURL   string    `db:"-" json:"-"`
	Tag       *Tag      `db:"-" json:"-"`
	TagName   string    `db:"-" json:"-"`
	TagID     int       `db:"tag_id" json:"tagID"`
	Voter     *User     `db:"-" json:"-"`
	VoterName string    `db:"-" json:"-"`
	VoterID   int       `db:"voter_id" json:"-"`
	Creator   *User     `db:"-" json:"-"`
	CreatorID int       `db:"creator_id" json:"-"`
	Tallied   bool      `db:"tallied" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

// FindByUK - find a given PostTagVote in the database by unique keys
func (v *PostTagVote) FindByUK(postID int, tagID int, userID int) error {
	dbPostTagVote := PostTagVote{}
	err := DBConn.Get(&dbPostTagVote, "SELECT * FROM posts_tags_votes WHERE post_id = ? and tag_id = ? and voter_id = ?", postID, tagID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	*v = dbPostTagVote
	return nil
}

// DeleteByUKWithTx - delete a PostTagVote in the database by unique keys
func (v *PostTagVote) DeleteByUKWithTx(tx Transaction, postID int, tagID int, userID int) error {
	_, err := tx.Exec("delete from posts_tags_votes where post_id = ? and tag_id = ? and voter_id = ?", postID, tagID, userID)
	return err
}

// CreatePostTagVote - create the PostTagVote with post and tag information
func (u *User) CreatePostTagVote(postID int, tagID int) error {
	post := &Post{}
	post.FindByID(postID)
	_, err := DBConn.Exec("INSERT INTO posts_tags_votes (post_id, tag_id, voter_id, creator_id) VALUES (?, ?, ?, ?)", postID, tagID, u.ID, post.AuthorID)
	return err
}

// CreatePostTagVote - create the PostTagVote with post and tag information
func (v *PostTagVote) CreatePostTagVote() error {
	post := &Post{}
	post.FindByID(v.PostID)
	// create a new PostTag association
	_, err := DBConn.Exec("INSERT INTO posts_tags_votes (post_id, tag_id, voter_id) VALUES (?, ?, ?)", v.PostID, v.TagID, v.VoterID)
	return err
}

// CreatePostTagVoteWithTx - create the PostTagVote with post and tag information
func (v *PostTagVote) CreatePostTagVoteWithTx(tx Transaction) error {
	// create a new PostTag association
	_, err := tx.Exec("INSERT INTO posts_tags_votes (post_id, tag_id, voter_id) VALUES (?, ?, ?)", v.PostID, v.TagID, v.VoterID)
	if err != nil {
		return err
	}
	return nil
}

// GetVotesByPostUser - query the rows from posts_tags_votes with post id and user id
func (v *PostTagVote) GetVotesByPostUser(postID int, userID int) ([]PostTagVote, error) {
	votes := []PostTagVote{}

	err := DBConn.Select(&votes, "select * from posts_tags_votes where post_id = ? and voter_id = ?", postID, userID)
	if err == sql.ErrNoRows {
		return votes, nil
	}
	return votes, err
}

// GetUntalliedVotes - query the rows from posts_tags_votes for votes for a specific user that haven't been tallied yet for payout
func (u *User) GetUntalliedVotes(before time.Time) ([]PostTagVote, error) {
	votes := []PostTagVote{}

	timestring := before.UTC().Format("2006-01-02 03:04:05 -0700 MST")
	err := DBConn.Select(&votes, "select * from posts_tags_votes where voter_id = ? AND created_at <= ? AND tallied = ? AND creator_id != ?", u.ID, timestring, 0, u.ID)
	return votes, err
}

// GetUntalliedVotes - returns all untallied votes in the system.
func (v *PostTagVote) GetUntalliedVotes() ([]PostTagVote, error) {
	votes := []PostTagVote{}

	err := DBConn.Select(&votes, "select * from posts_tags_votes where tallied = ?", 0)
	if err == sql.ErrNoRows {
		return votes, nil
	}
	return votes, err
}

// MarkVotesAsTallied - Mark all of the votes in an array as being tallied.
func (v *PostTagVote) MarkVotesAsTallied(votes []PostTagVote) error {
	var ids []int
	for _, vote := range votes {
		ids = append(ids, vote.ID)
	}

	query := "UPDATE posts_tags_votes SET tallied = 1 where id IN = (?)"
	generatedQuery, args, err := sqlx.In(query, 1, ids)
	_, err = DBConn.Exec(generatedQuery, args...)
	if err == sql.ErrNoRows {
		return nil
	}
	return err
}

// GetVotesByPostTag - query the rows from posts_tags_votes with post id and tag id
func (v *PostTagVote) GetVotesByPostTag(postID int, tagID int) ([]PostTagVote, error) {
	votes := []PostTagVote{}

	err := DBConn.Select(&votes, "select * from posts_tags_votes where post_id = ? and tag_id = ?", postID, tagID)
	if err == sql.ErrNoRows {
		return votes, nil
	}
	return votes, err
}
