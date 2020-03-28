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
	PostURL   string        `db:"-" json:"-"`
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
	err := DBConn.Get(&dbPostTag, "SELECT * FROM posts_tags WHERE id = ?", id)
	if err != nil {
		return err
	}

	*p = dbPostTag
	return nil
}

// GetPostsByTags - query the post ids by the tags
func (p *PostTag) GetPostsByTags(tags string) ([]int, error) {
	ids := []int{}

	err := DBConn.Select(&ids, "select t1.post_id from posts_tags t1 join tags t2 on t1.tag_id = t2.id and t2.name in (?) ORDER BY score DESC", tags)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// FindByUK - query the PostTag by unique key: post id and tag id
func (p *PostTag) FindByUK(postID int, tagID int) error {
	dbPostTag := PostTag{}
	err := DBConn.Get(&dbPostTag, "select * from posts_tags where post_id = ? and tag_id = ?", postID, tagID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	*p = dbPostTag
	return nil
}

// IncrementScore - increment the score by post id and tag id
func (p *PostTag) IncrementScore(postID int, tagID int) error {
	_, err := DBConn.Exec("update posts_tags set score=score+1 where post_id = ? and tag_id = ?", postID, tagID)
	return err
}

// IncrementScoreWithTx - increment the score by post id and tag id
func (p *PostTag) IncrementScoreWithTx(tx Transaction, postID int, tagID int) error {
	_, err := tx.Exec("update posts_tags set score=score+1 where post_id = ? and tag_id = ?", postID, tagID)
	return err
}

// DecrementScoreWithTx - decrement the score by post id and tag id
func (p *PostTag) DecrementScoreWithTx(tx Transaction, postID int, tagID int) error {
	_, err := tx.Exec("update posts_tags set score=score-1 where post_id = ? and tag_id = ?", postID, tagID)
	return err
}

// DeleteByUKWithTx - delete a PostTag in the database by unique keys
func (p *PostTag) DeleteByUKWithTx(tx Transaction, postID int, tagID int) error {
	_, err := tx.Exec("delete from posts_tags where post_id = ? and tag_id = ?", postID, tagID)
	return err
}

// CreatePostTag - create the PostTag with post and tag information
func (p *PostTag) CreatePostTag() error {
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a new PostTag association
	_, err := DBConn.Exec("INSERT INTO posts_tags (post_id, tag_id, score, created_at) VALUES (?, ?, 1, ?)", p.PostID, p.TagID, now)
	if err != nil {
		return err
	}
	return nil
}

func CreatePostTagFromObjects(post Post, tag Tag) (int64, error) {
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a new PostTag association
	res, err := DBConn.Exec("INSERT INTO posts_tags (post_id, tag_id, score, created_at) VALUES (?, ?, 1, ?)", post.ID, tag.ID, now)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

// CreatePostTagWithTx - create the PostTag with post and tag information
func (p *PostTag) CreatePostTagWithTx(tx Transaction) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a new PostTag association
	_, err := tx.Exec("INSERT INTO posts_tags (post_id, tag_id, score, created_at) VALUES (?, ?, 1, ?)", p.PostID, p.TagID, now)
	if err != nil {
		return err
	}
	return nil
}
