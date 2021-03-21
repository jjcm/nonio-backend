package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Comment - struct representation of a single comment
type Comment struct {
	ID                     int           `db:"id" json:"-"`
	Post                   Post          `db:"-" json:"-"`
	PostID                 int           `db:"post_id" json:"-"`
	PostURL                string        `db:"-" json:"post"`
	CreatedAt              time.Time     `db:"created_at" json:"date"`
	Content                string        `db:"content" json:"content"`
	ParentID               int           `db:"parent_id" json:"-"`
	User                   string        `db:"-" json:"user"`
	Author                 User          `db:"-" json:"-"`
	AuthorID               sql.NullInt32 `db:"author_id" json:"-"`
	Upvotes                int           `db:"-" json:"upvotes"`
	Downvotes              int           `db:"-" json:"downvotes"`
	LineageScore           int           `db:"lineage_score" json:"lineage_score"`
	DescendentCommentCount int           `db:"descendent_comment_count" json:"descendent_comment_count"`
}

// MarshalJSON custom JSON builder for Comment structs
func (c *Comment) MarshalJSON() ([]byte, error) {
	// populate user if it currently isn't hydrated
	if c.Author.ID == 0 {
		if c.AuthorID.Valid {
			c.Author.FindByID(int(c.AuthorID.Int32))
		} else {
			anonymous := User{}
			anonymous.Username = "Anonymous coward"
			c.Author = anonymous
		}
	}
	if c.Post.ID == 0 {
		c.Post.FindByID(c.PostID)
	}

	// return the custom JSON for this post
	return json.Marshal(&struct {
		ID                     int    `json:"id"`
		Date                   int64  `json:"date"`
		Post                   string `json:"post"`
		Content                string `json:"content"`
		User                   string `json:"user"`
		Upvotes                int    `json:"upvotes"`
		Downvotes              int    `json:"downvotes"`
		Parent                 int    `json:"parent"`
		LineageScore           int    `json:"lineage_score"`
		DescendentCommentCount int    `json:"descendent_comment_count"`
	}{
		ID:                     c.ID,
		Date:                   c.CreatedAt.UnixNano() / int64(time.Millisecond),
		Post:                   c.Post.URL,
		Content:                c.Content,
		User:                   c.Author.GetDisplayName(),
		Upvotes:                c.Upvotes,
		Downvotes:              c.Downvotes,
		Parent:                 c.ParentID,
		LineageScore:           c.LineageScore,
		DescendentCommentCount: c.DescendentCommentCount,
	})
}

/************************************************/
/******************** CREATE ********************/
/************************************************/

// CreateComment will try and create a comment in the database
func (u *User) CreateComment(post Post, parent *Comment, content string) (Comment, error) {
	c := Comment{}
	now := time.Now().Format("2006-01-02 15:04:05")

	if u.ID == 0 || post.ID == 0 {
		return c, errors.New("Can't create a comment for an invalid user or post")
	}

	var commentParentID int
	if parent != nil {
		commentParentID = parent.ID
	}

	result, err := DBConn.Exec("INSERT INTO comments (author_id, post_id, created_at, content, parent_id) VALUES (?, ?, ?, ?, ?)", u.ID, post.ID, now, content, commentParentID)
	if err != nil {
		return c, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return c, err
	}

	c.FindByID(int(insertID))
	return c, err
}

/************************************************/
/********************* READ *********************/
/************************************************/

// FindByID - find a given comment in the database by its primary key
func (c *Comment) FindByID(id int) error {
	Log.Info(fmt.Sprintf("finding comment %v", id))
	dbComment := Comment{}
	if err := DBConn.Get(&dbComment, "SELECT * FROM comments WHERE id = ?", id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	*c = dbComment
	return nil
}

// GetCommentsByPost returns the comments for one post order by lineage score
func GetCommentsByPost(id int) ([]*Comment, error) {
	comments := []*Comment{}

	if err := DBConn.Select(&comments, "SELECT * FROM comments where post_id = ? order by lineage_score desc limit 100", id); err != nil {
		fmt.Println("this turned out bad")
		return nil, err
	}

	return comments, nil
}

// GetComments will return comments associated with the current post
func (p *Post) GetComments(depthLimit int) ([]Comment, error) {
	var err error
	var comments []Comment

	// this is a temporary work around to let front end dev get back at it...
	err = DBConn.Select(&comments, "SELECT * FROM comments WHERE post_id = ?", p.ID)
	return comments, err
}

/************************************************/
/******************** UPDATE ********************/
/************************************************/

// IncrementLineageScoreWithTx - increment the lineage score by comment id
func (c *Comment) IncrementLineageScoreWithTx(tx Transaction) error {
	id := c.ID

	for {
		comment := Comment{}
		err := comment.FindByID(id)
		if err != nil {
			return err
		}

		if comment.ID == 0 {
			err = fmt.Errorf("Comment does not exist")
			return err
		}

		_, err = tx.Exec("update comments set lineage_score=lineage_score+1 where id = ?", id)
		if err != nil {
			return fmt.Errorf("Error incrementing lineage score: %v", err)
		}

		if comment.ParentID == 0 {
			break
		}

		id = comment.ParentID
	}

	return nil
}

// DecrementLineageScoreWithTx - decrement the lineage score by comment id
func (c *Comment) DecrementLineageScoreWithTx(tx Transaction) error {
	id := c.ID

	for {
		comment := Comment{}
		err := comment.FindByID(id)
		if err != nil {
			return err
		}

		if comment.ID == 0 {
			err = fmt.Errorf("Comment does not exist")
			return err
		}

		_, err = tx.Exec("update comments set lineage_score=lineage_score-1 where id = ?", id)
		if err != nil {
			return fmt.Errorf("Error incrementing lineage score: %v", err)
		}

		if comment.ParentID == 0 {
			break
		}

		id = comment.ParentID
	}

	return nil
}

// IncrementDescendentComment - increment the descendent comment count
func (c *Comment) IncrementDescendentComment(id int) error {
	_, err := DBConn.Exec("update comments set descendent_comment_count=descendent_comment_count+1 where id = ?", id)
	return err
}

/************************************************/
/******************** DELETE ********************/
/************************************************/

// AbandonComment removes the user from the comment, but leaves the content
func (u *User) AbandonComment(comment *Comment) error {
	if u.ID == 0 || comment.ID == 0 {
		return errors.New("Can't abandon a comment for an invalid user or comment")
	}

	_, err := DBConn.Exec("UPDATE comments SET author_id = NULL WHERE id = ?", comment.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment removes it from the db
func (u *User) DeleteComment(comment *Comment) error {
	if u.ID == 0 || comment.ID == 0 {
		return errors.New("Can't delete a comment for an invalid user or comment")
	}

	_, err := DBConn.Exec("DELETE FROM comments WHERE id = ?", comment.ID)
	if err != nil {
		return err
	}

	return nil
}
