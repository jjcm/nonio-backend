package models

import (
	"encoding/json"
	"time"
)

// Subscription - code representation of a users subscription to a tag
type Subscription struct {
	Tag       *Tag      `db:"-" json:"-"`
	TagName   string    `db:"-" json:"tag"`
	TagID     int       `db:"tag_id" json:"tagID"`
	User      User      `db:"-" json:"user"`
	UserID    int       `db:"user_id" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// MarshalJSON custom JSON builder for Tag structs
func (s *Subscription) MarshalJSON() ([]byte, error) {
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

// createSubscription - create a subscription in the database for a tag
func createSubscription(tag Tag, user User) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := DBConn.Exec("INSERT INTO subscriptions (tag_id, user_id, created_at) VALUES (?, ?, ?)", tag.ID, user.ID, now)
	if err != nil {
		return err
	}
	return nil
}

// SubscriptionFactory will create and return an instance of a tag
func SubscriptionFactory(tag Tag, user User) (Subscription, error) {
	s := Subscription{}
	err := createSubscription(tag, user)
	if err != nil {
		return s, err
	}
	err = s.FindSubscription(tag, user)

	return s, err
}

// FindSubscription - find a given tag in the database by its primary key
func (s *Subscription) FindSubscription(tag Tag, user User) error {
	dbSubscription := Subscription{}
	err := DBConn.Get(&dbSubscription, "SELECT * FROM subscriptions WHERE tag_id = ? AND user_id = ?", tag.ID, user.ID)
	if err != nil {
		return err
	}

	*s = dbSubscription
	return nil
}
