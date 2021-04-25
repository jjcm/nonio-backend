package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Tag - code representation of a single tag
type Tag struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	UserID    int       `db:"user_id" json:"-"`
	Count     int       `db:"count" json:"count"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// MarshalJSON custom JSON builder for Tag structs
func (t *Tag) MarshalJSON() ([]byte, error) {
	// return the custom JSON for this post
	return json.Marshal(&struct {
		Tag   string `json:"tag"`
		Count int    `json:"count"`
	}{
		Tag:   t.Name,
		Count: t.Count,
	})
}

// ToJSON - get a string representation of this Tag in JSON
func (t *Tag) ToJSON() string {
	jsonData, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}

/************************************************/
/******************** CREATE ********************/
/************************************************/

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

/************************************************/
/********************* READ *********************/
/************************************************/

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

// GetTags - get tags out of the database offset by an integer
func GetTags(offset int, limit int) ([]Tag, error) {
	tags := []Tag{}
	err := DBConn.Select(&tags, "SELECT id, name, count FROM tags WHERE count > 0 ORDER BY count DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return tags, err
	}

	return tags, nil
}

// GetTagsByPrefix - get tags that match a specific prefix
func GetTagsByPrefix(prefix string) ([]Tag, error) {
	tags := []Tag{}
	err := DBConn.Select(&tags, "SELECT name, count FROM tags WHERE name LIKE ? ORDER BY count DESC LIMIT 100", fmt.Sprintf("%v%%", prefix))
	if err != nil {
		return tags, err
	}

	return tags, nil
}

/************************************************/
/******************** UPDATE ********************/
/************************************************/
// Not needed

/************************************************/
/******************** DELETE ********************/
/************************************************/
// TODO
