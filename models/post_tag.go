package models

import (
	"database/sql"
	"time"
)

// PostTag - struct representation of a single post-tag
type PostTag struct {
	ID        int           `db:"id" json:"-"`
	Post      *Post         `db:"-" json:"-"`
	PostID    int           `db:"post_id" json:"-"`
	PostURL   string        `db:"-" json:"post"`
	Tag       *Tag          `db:"-" json:"-"`
	TagName   string        `db:"-" json:"tag"`
	TagID     int           `db:"tag_id" json:"-"`
	Score     int           `db:"score" json:"score"`
	CreatedAt time.Time     `db:"created_at" json:"-"`
	Votes     []PostTagVote `db:"-" json:"-"`
}

// FindByID - find a given PostTag in the database by its primary key
func (p *PostTag) FindByID(id int) error {
	dbPostTag := PostTag{}
	err := DBConn.Get(&dbPostTag, "SELECT * FROM tags WHERE id = ?", id)
	if err != nil {
		return err
	}

	*p = dbPostTag
	return nil
}

// FindByPostTagIds - query the PostTag by post id and tag id
func (p *PostTag) FindByPostTagIds(postId int, tagId int) (*PostTag, error) {
	dbPostTag := PostTag{}
	err := DBConn.Get(&dbPostTag, "select * from posts_tags where post_id = ? and tag_id = ?", postId, tagId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &dbPostTag, nil
}

// CreatePostTag - create the PostTag with post and tag information
func (p *PostTag) CreatePostTag() error {
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a new PostTag association
	_, err := DBConn.Exec("INSERT INTO posts_tags (post_id, tag_id, score, created_at) VALUES (?, ?, 1, ?)", p.Post.ID, p.Tag.ID, now)
	if err != nil {
		return err
	}
	return nil
}
