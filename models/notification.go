package models

import (
	"encoding/json"
	"time"
)

// Notification - code representation of responses to a user's comments or posts
type Notification struct {
	ID        int       `db:"id" json:"-"`
	User      int       `db:"user_id" json:"-"`
	Comment   Comment   `db:"-" json:"-"`
	CommentID int       `db:"comment_id" json:"commentID"`
	Read      bool      `db:"read" json:"read"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}

// MarshalJSON custom JSON builder for Tag structs
func (n *Notification) MarshalJSON() ([]byte, error) {
	// hydrate the comment
	if n.Comment.ID == 0 {
		n.Comment.FindByID(n.CommentID)
	}

	// hydrate the comment's post
	if n.Comment.Post.ID == 0 {
		n.Comment.Post.FindByID(n.Comment.PostID)
	}

	// return the custom JSON for this post
	return json.Marshal(&struct {
		ID          int    `json:"id"`
		CommentID   int    `json:"commentID"`
		CommentDate int64  `json:"commentDate"`
		Post        string `json:"post"`
		PostTitle   string `json:"postTitle"`
		Content     string `json:"content"`
		User        string `json:"user"`
		Upvotes     int    `json:"upvotes"`
		Downvotes   int    `json:"downvotes"`
		Parent      int    `json:"parent"`
		Edited      bool   `json:"edited"`
		Read        bool   `json:"read"`
	}{
		ID:          n.ID,
		CommentID:   n.CommentID,
		CommentDate: n.CreatedAt.UnixNano() / int64(time.Millisecond),
		Post:        n.Comment.Post.URL,
		PostTitle:   n.Comment.Post.Title,
		Content:     n.Comment.Content,
		User:        n.Comment.Author.GetDisplayName(),
		Upvotes:     n.Comment.Upvotes,
		Downvotes:   n.Comment.Downvotes,
		Parent:      n.Comment.ParentID,
		Edited:      n.Comment.Edited,
		Read:        n.Read,
	})
}

/************************************************/
/******************** CREATE ********************/
/************************************************/

// CreateNotification - create a notification in the database for a comment
func (u *User) CreateNotification(comment Comment) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := DBConn.Exec("INSERT INTO notifications (comment_id, user_id, created_at) VALUES (?, ?, ?)", comment.ID, u.ID, now)
	if err != nil {
		return err
	}

	return nil
}

/************************************************/
/********************* READ *********************/
/************************************************/

// GetNotifications - get all notifications for a user
func (u *User) GetNotifications() ([]Notification, error) {
	notifications := []Notification{}
	return notifications, nil

	/* TODO
	// run the correct sql query
	var query = "SELECT * FROM subscriptions WHERE user_id = ?"
	err := DBConn.Select(&subscriptions, query, u.ID)
	if err != nil {
		return subscriptions, err
	}

	return subscriptions, nil
	*/
}
