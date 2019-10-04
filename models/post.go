package models

import (
	"encoding/json"
	"strconv"
	"time"
)

// Post - struct representation of a single post
type Post struct {
	Title     string
	Author    User
	CreatedAt time.Time
	URL       string
	Tags      []Tag
}

// MarshalJSON custom JSON builder for Post structs
func (p *Post) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Title     string `json:"title"`
		UserName  string `json:"user"`
		TimeStamp string `json:"time"`
		URL       string `json:"url"`
		Tags      []Tag  `json:"tags"`
	}{
		Title:     p.Title,
		UserName:  p.Author.Name,
		TimeStamp: strconv.Itoa(int(p.CreatedAt.Unix())),
		URL:       p.URL,
		Tags:      p.Tags,
	})
}
