package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Tag - code representation of a single tag
type Tag struct {
	ID                     int       `db:"id" json:"id"`
	Name                   string    `db:"name" json:"name"`
	Author                 User      `db:"-" json:"createdBy"`
	UserID                 int       `db:"user_id" json:"-"`
	Score                  int       `db:"-"`
	CreatedAt              time.Time `db:"created_at" json:"createdAt"`
	DateAssociatedWithPost time.Time `db:"-" json:"-"`
}

// MarshalJSON custom JSON builder for Tag structs
func (t *Tag) MarshalJSON() ([]byte, error) {
	// return the custom JSON for this post
	return json.Marshal(&struct {
		Tag       string    `json:"tag"`
		UpVotes   int       `json:"upvotes"`
		DownVotes int       `json:"downvotes"`
		Date      time.Time `json:"date"`
	}{
		Tag:       t.Name,
		UpVotes:   0,                        // TODO - get this correctly
		DownVotes: 0,                        // TODO - get this correctly
		Date:      t.DateAssociatedWithPost, // TODO - figure this out
	})
}

// GetTags - get tags out of the database offset by an integer
func GetTags(offset int, limit int) ([]Tag, error) {
	tags := []Tag{}
	err := DBConn.Select(&tags, "SELECT id, name FROM tags LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return tags, err
	}

	return tags, nil
}

// createTag - create a tag in the database by a given word
func createTag(tag string, author User) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := DBConn.Exec("INSERT INTO tags (name, user_id, created_at) VALUES (?, ?, ?)", tag, author.ID, now)
	if err != nil {
		return err
	}
	return nil
}

// TagFactory will create and return an instance of a tag
func TagFactory(tag string, author User) (Tag, error) {
	t := Tag{}
	err := createTag(tag, author)
	if err != nil {
		return t, err
	}
	err = t.FindByTagName(tag)

	return t, err
}

// FindByID - find a given tag in the database by its primary key
func (t *Tag) FindByID(id int) error {
	dbTag := Tag{}
	err := DBConn.Get(&dbTag, "SELECT * FROM tags WHERE id = ?", id)
	if err != nil {
		return err
	}

	*t = dbTag
	return nil
}

// FindByTagName - find a tag in the database by it's name
func (t *Tag) FindByTagName(name string) error {
	dbTag := Tag{}
	err := DBConn.Get(&dbTag, "SELECT * FROM tags WHERE name = ?", name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	*t = dbTag
	return nil
}
