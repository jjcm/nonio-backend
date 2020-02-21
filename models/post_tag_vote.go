package models

// PostTagVote - struct representation of a single post-tag-vote
type PostTagVote struct {
	ID        int    `db:"id" json:"-"`
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

// FindByID - find a given PostTagVote in the database by its primary key
func (v *PostTagVote) FindByID(id int) error {
	dbPostTagVote := PostTagVote{}
	err := DBConn.Get(&dbPostTagVote, "SELECT * FROM posts_tags_votes WHERE id = ?", id)
	if err != nil {
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
