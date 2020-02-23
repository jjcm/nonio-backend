package models

import "database/sql"

// PostTagVote - struct representation of a single post-tag-vote
type PostTagVote struct {
	Post      *Post  `db:"-" json:"-"`
	PostID    int    `db:"post_id" json:"-"`
	PostURL   string `db:"-" json:"post"`
	Tag       *Tag   `db:"-" json:"-"`
	TagName   string `db:"-" json:"tag"`
	TagID     int    `db:"tag_id" json:"-"`
	Voter     *User  `db:"-" json:"-"`
	VoterName string `db:"-" json:"user"`
	VoterID   int    `db:"user_id" json:"-"`
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

// CreatePostTagVote - create the PostTagVote with post and tag information
func (v *PostTagVote) CreatePostTagVote() error {
	// create a new PostTag association
	_, err := DBConn.Exec("INSERT INTO posts_tags_votes (post_id, tag_id, voter_id) VALUES (?, ?, ?)", v.PostID, v.TagID, v.VoterID)
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
