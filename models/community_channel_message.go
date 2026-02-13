package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// CommunityChannelMessage - a single message in a text channel
type CommunityChannelMessage struct {
	ID        int       `db:"id" json:"id"`
	ChannelID int       `db:"channel_id" json:"channelID"`
	AuthorID  int       `db:"author_id" json:"authorID"`
	Content   string    `db:"content" json:"content"`
	ImageURL  string    `db:"image_url" json:"imageUrl,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Author    User      `db:"-" json:"-"`
}

// MarshalJSON custom JSON for API
func (m *CommunityChannelMessage) MarshalJSON() ([]byte, error) {
	username := ""
	if m.Author.ID != 0 {
		username = m.Author.GetDisplayName()
	} else if m.AuthorID > 0 {
		u := User{}
		_ = u.FindByID(m.AuthorID)
		username = u.GetDisplayName()
	}
	return json.Marshal(&struct {
		ID        int    `json:"id"`
		ChannelID int    `json:"channelID"`
		AuthorID  int    `json:"authorID"`
		User      string `json:"user"`
		Content   string `json:"content"`
		ImageURL  string `json:"imageUrl,omitempty"`
		Date      int64  `json:"date"`
	}{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		AuthorID:  m.AuthorID,
		User:      username,
		Content:   m.Content,
		ImageURL:  m.ImageURL,
		Date:      m.CreatedAt.UnixNano() / int64(time.Millisecond),
	})
}

// CreateChannelMessage - insert a message (caller must verify channel is text and user is member)
func (u *User) CreateChannelMessage(channelID int, content, imageURL string) (CommunityChannelMessage, error) {
	msg := CommunityChannelMessage{}
	now := time.Now().Format("2006-01-02 15:04:05")

	ch := CommunityChannel{}
	if err := ch.FindByID(channelID); err != nil {
		return msg, err
	}
	if ch.ID == 0 {
		return msg, errors.New("channel not found")
	}
	if ch.Kind != ChannelKindText {
		return msg, fmt.Errorf("channel is not a text channel")
	}

	if content == "" && imageURL == "" {
		return msg, errors.New("content or image is required")
	}

	result, err := DBConn.Exec(
		"INSERT INTO community_channel_messages (channel_id, author_id, content, image_url, created_at) VALUES (?, ?, ?, ?, ?)",
		channelID, u.ID, content, imageURL, now,
	)
	if err != nil {
		return msg, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return msg, err
	}
	err = msg.FindByID(int(insertID))
	return msg, err
}

// FindByID - load message by id
func (m *CommunityChannelMessage) FindByID(id int) error {
	row := CommunityChannelMessage{}
	err := DBConn.Get(&row, "SELECT * FROM community_channel_messages WHERE id = ?", id)
	if err != nil {
		return err
	}
	*m = row
	return nil
}

// ChannelMessageQueryParams - list params for messages
type ChannelMessageQueryParams struct {
	ChannelID int
	BeforeID  int
	Limit     int
}

// GetChannelMessages - list messages for a channel, newest first; if BeforeID > 0, return older than that message
func GetChannelMessages(params *ChannelMessageQueryParams) ([]*CommunityChannelMessage, error) {
	limit := params.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	query := "SELECT * FROM community_channel_messages WHERE channel_id = ?"
	args := []interface{}{params.ChannelID}
	if params.BeforeID > 0 {
		query += " AND id < ?"
		args = append(args, params.BeforeID)
	}
	query += " ORDER BY id DESC LIMIT ?"
	args = append(args, limit)

	msgs := []*CommunityChannelMessage{}
	err := DBConn.Select(&msgs, query, args...)
	return msgs, err
}
