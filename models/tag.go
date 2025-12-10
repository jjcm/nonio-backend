package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Tag - code representation of a single tag
type Tag struct {
	ID          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	UserID      int       `db:"user_id" json:"-"`
	Count       int       `db:"count" json:"count"`
	CommunityID *int      `db:"community_id" json:"communityID,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
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
func createTag(tag string, author User, communityID int) error {
	if tag == "" {
		return fmt.Errorf("Tag cannot be an empty string")
	}

	if strings.ContainsAny(tag, " ") {
		return fmt.Errorf("Tag cannot contain spaces")
	}

	if strings.ContainsAny(tag, "#") {
		return fmt.Errorf("Tag cannot contain hashes")
	}

	if strings.ContainsAny(tag, "<>='\"./|\\") {
		return fmt.Errorf("Tag cannot contain html elements")
	}

	//checks the length of the TagName, if it's more than 30 characters, returns an error
	if len(tag) > 20 {
		return fmt.Errorf("Tag cannot be more than 20 characters")
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// Convert communityID to nullable value
	var communityIDPtr interface{}
	if communityID == 0 {
		communityIDPtr = nil
	} else {
		communityIDPtr = communityID
	}

	_, err := DBConn.Exec("INSERT INTO tags (name, user_id, community_id, created_at) VALUES (?, ?, ?, ?)", tag, author.ID, communityIDPtr, now)
	if err != nil {
		return err
	}
	return nil
}

// TagFactory will create and return an instance of a tag
// If no community is provided, it defaults to the root (0).
func TagFactory(tag string, author User, communityID ...int) (Tag, error) {
	t := Tag{}
	communityKey := 0
	if len(communityID) > 0 {
		communityKey = communityID[0]
	}

	err := createTag(tag, author, communityKey)
	if err != nil {
		return t, err
	}
	err = t.FindByTagName(tag, communityKey)

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
// If communityID is not provided, defaults to 0 (frontpage).
func (t *Tag) FindByTagName(name string, communityID ...int) error {
	dbTag := Tag{}
	communityKey := 0
	if len(communityID) > 0 {
		communityKey = communityID[0]
	}
	var err error
	if communityKey == 0 {
		err = DBConn.Get(&dbTag, "SELECT * FROM tags WHERE name = ? AND (community_id IS NULL OR community_id = 0)", name)
	} else {
		err = DBConn.Get(&dbTag, "SELECT * FROM tags WHERE name = ? AND community_id = ?", name, communityKey)
	}

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
// communityID defaults to 0 (frontpage) if omitted.
func GetTags(offset int, limit int, communityID ...int) ([]Tag, error) {
	tags := []Tag{}
	var err error
	communityKey := 0
	if len(communityID) > 0 {
		communityKey = communityID[0]
	}

	if communityKey == 0 {
		err = DBConn.Select(&tags, "SELECT id, name, count FROM tags WHERE count > 0 AND (community_id IS NULL OR community_id = 0) ORDER BY count DESC LIMIT ? OFFSET ?", limit, offset)
	} else {
		err = DBConn.Select(&tags, "SELECT id, name, count FROM tags WHERE count > 0 AND community_id = ? ORDER BY count DESC LIMIT ? OFFSET ?", communityKey, limit, offset)
	}

	if err != nil {
		return tags, err
	}

	return tags, nil
}

// GetTagsByPrefix - get tags that match a specific prefix
func GetTagsByPrefix(prefix string, communityID int) ([]Tag, error) {
	tags := []Tag{}
	var err error
	if communityID == 0 {
		err = DBConn.Select(&tags, "SELECT name, count FROM tags WHERE name LIKE ? AND (community_id IS NULL OR community_id = 0) ORDER BY count DESC LIMIT 100", fmt.Sprintf("%v%%", prefix))
	} else {
		err = DBConn.Select(&tags, "SELECT name, count FROM tags WHERE name LIKE ? AND community_id = ? ORDER BY count DESC LIMIT 100", fmt.Sprintf("%v%%", prefix), communityID)
	}

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
