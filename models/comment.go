package models

import (
	"encoding/json"
	"time"
)

// Comment - struct representation of a single comment
type Comment struct {
	ID        int       `db:"id" json:"-"`
	Post      Post      `db:"-" json:"-"`
	PostID    int       `db:"post_id" json:"-"`
	PostURL   string    `db:"-" json:"post"`
	CreatedAt time.Time `db:"created_at" json:"date"`
	Type      string    `db:"type" json:"type"`
	Content   string    `db:"content" json:"content"`
	Text      string    `db:"text" json:"text"`
	ParentID  int       `db:"parent_id" json:"-"`
	User      string    `db:"-" json:"user"`
	Author    User      `db:"-" json:"-"`
	AuthorID  int       `db:"author_id" json:"-"`
	UpVotes   []Vote    `db:"-" json:"upvotes"`
	DownVotes []Vote    `db:"-" json:"downvotes"`
	Children  []Comment `db:"-" json:"children"`
}

// MarshalJSON custom JSON builder for Comment structs
func (c *Comment) MarshalJSON() ([]byte, error) {
	// populate user if it currently isn't hydrated
	if c.Author.ID == 0 {
		c.Author.FindByID(c.AuthorID)
	}
	if c.Post.ID == 0 {
		c.Post.FindByID(c.PostID)
	}

	// return the custom JSON for this post
	return json.Marshal(&struct {
		ID        int       `json:"id"`
		Date      int64     `json:"date"`
		Post      string    `json:"post"`
		Type      string    `json:"type"`
		Content   string    `json:"content"`
		Text      string    `json:"text"`
		User      string    `json:"user"`
		UpVotes   int       `json:"upvotes"`
		DownVotes int       `json:"downvotes"`
		Parent    int       `json:"parent"`
		Children  []Comment `json:"children"`
	}{
		ID:        c.ID,
		Date:      c.CreatedAt.UnixNano() / int64(time.Millisecond),
		Post:      c.Post.URL,
		Type:      c.Type,
		Content:   c.Content,
		Text:      c.Text,
		User:      c.Author.GetDisplayName(),
		UpVotes:   len(c.UpVotes),
		DownVotes: len(c.DownVotes),
		Parent:    c.ParentID,
	})
}

// FindByID - find a given comment in the database by its primary key
func (c *Comment) FindByID(id int) error {
	dbComment := Comment{}
	err := DBConn.Get(&dbComment, "SELECT * FROM comments WHERE id = ?", id)
	if err != nil {
		return err
	}

	*c = dbComment
	return nil
}

// StructureComments will take in an un sorted list of comments and put them in
// the correct structure for frontend display
// TODO: this function isn't very efficient. Good enough for a proof of concept,
// but this is quite a bottle neck.
func StructureComments(comments []Comment) []Comment {
	// for _, var := range var {

	// }
	return comments
}
