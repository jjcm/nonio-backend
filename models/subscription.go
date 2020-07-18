package models

import (
	"encoding/json"
	"time"
)

// Subscription - code representation of a users subscription to a tag
type Subscription struct {
	ID        int       `db:"id" json:"-"`
	Tag       *Tag      `db:"-" json:"-"`
	TagName   string    `db:"-" json:"tag"`
	TagID     int       `db:"tag_id" json:"tagID"`
	User      User      `db:"-" json:"user"`
	UserID    int       `db:"user_id" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// MarshalJSON custom JSON builder for Tag structs
func (s *Subscription) MarshalJSON() ([]byte, error) {
	// hydrate the user
	if s.User.ID == 0 {
		s.User.FindByID(s.UserID)
	}

	// hydrate the tags
	if s.Tag.ID == 0 {
		s.Tag.FindByID(s.TagID)
	}
	// return the custom JSON for this post
	return json.Marshal(&struct {
		Tag  string `json:"tag"`
		User string `json:"user"`
	}{
		Tag:  s.Tag.Name,
		User: s.User.GetDisplayName(),
	})
}

// ToJSON - get a string representation of this Tag in JSON
func (s *Subscription) ToJSON() string {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}

// CreateSubscription - create a subscription in the database for a tag
func (u *User) CreateSubscription(tag Tag) (Subscription, error) {
	s := Subscription{}
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := DBConn.Exec("INSERT INTO subscriptions (tag_id, user_id, created_at) VALUES (?, ?, ?)", tag.ID, u.ID, now)
	if err != nil {
		return s, err
	}

	err = s.FindSubscription(tag.ID, u.ID)
	return s, err
}

// DeleteSubscription - create a subscription in the database for a tag
func (u *User) DeleteSubscription(tag Tag) error {
	_, err := DBConn.Exec("DELETE FROM subscriptions WHERE tag_id = ? AND user_id = ?", tag.ID, u.ID)
	return err
}

// FindSubscription - find a given tag in the database by the tag/user pairing
func (s *Subscription) FindSubscription(tagID int, userID int) error {
	dbSubscription := Subscription{}
	err := DBConn.Get(&dbSubscription, "SELECT * FROM subscriptions WHERE tag_id = ? AND user_id = ?", tagID, userID)
	if err != nil {
		return err
	}

	*s = dbSubscription
	return nil
}
