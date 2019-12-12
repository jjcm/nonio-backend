package models

import (
	"encoding/json"
	"time"
)

// Post - struct representation of a single post
type Post struct {
	ID        int       `db:"id" json:"-"`
	Title     string    `db:"title" json:"title"`
	URL       string    `db:"url" json:"url"`
	Author    User      `db:"-" json:"user"`
	AuthorID  int       `db:"user_id" json:"-"`
	Thumbnail string    `db:"thumbnail" json:"thumbnail"`
	Score     int       `db:"score" json:"score"`
	Content   string    `db:"content" json:"content"`
	Type      string    `db:"type" json:"type"`
	CreatedAt time.Time `db:"created_at" json:"date"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Tags      []Tag     `db:"-"`
}

// MarshalJSON custom JSON builder for Post structs
func (p *Post) MarshalJSON() ([]byte, error) {
	// build tag array for JS if the tag list is currently empty
	if len(p.Tags) < 1 {
		p.getTags()
	}

	// populate user if it currently isn't hydrated
	if p.Author.ID == 0 {
		p.Author.FindByID(p.AuthorID)
	}

	// return the custom JSON for this post
	return json.Marshal(&struct {
		Title     string `json:"title"`
		UserName  string `json:"user"`
		TimeStamp int64  `json:"time"`
		URL       string `json:"url"`
		Tags      []Tag  `json:"tags"`
	}{
		Title:     p.Title,
		UserName:  p.Author.Name,
		TimeStamp: p.CreatedAt.UnixNano() / int64(time.Millisecond),
		URL:       p.URL,
		Tags:      p.Tags,
	})
}

// ToJSON - get a string representation of this Post in JSON
func (p *Post) ToJSON() string {
	jsonData, err := json.Marshal(p)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}

// FindByURL - find a given post in the database by its URL
func (p *Post) FindByURL(url string) error {
	dbPost := Post{}
	err := DBConn.Get(&dbPost, "SELECT * FROM posts WHERE url = ?", url)
	if err != nil {
		return err
	}

	*p = dbPost
	return nil
}

// FindByID - find a given post in the database by its primary kye
func (p *Post) FindByID(id int) error {
	dbPost := Post{}
	err := DBConn.Get(&dbPost, "SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}

	*p = dbPost
	return nil
}

// GetCreatedAtTimestamp - get the created at timestamp in the predetermined format
func (p *Post) GetCreatedAtTimestamp() string {
	return p.CreatedAt.Format("2006-01-02 03:04PM")
}

// AddTag - associate a post with an existing tag
func (p *Post) AddTag(t Tag) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a new tag association
	_, err := DBConn.Exec("INSERT INTO posts_tags (post_id, tag_id, created_at) VALUES (?, ?, ?)", p.ID, t.ID, now)
	if err != nil {
		return err
	}

	err = p.getTags()
	if err != nil {
		return err
	}
	return nil
}

func (p *Post) getTags() error {
	tags := []Tag{}
	err := DBConn.Select(&tags, "SELECT * FROM `tags` WHERE id IN (SELECT `tag_id` FROM posts_tags WHERE post_id = ?)", p.ID)
	if err != nil {
		return err
	}
	p.Tags = tags
	return nil
}

// GetPostsByScoreSince - get posts from the database that have been created since
// the provided cutoff time, and offset the results
func GetPostsByScoreSince(cutoff time.Time, offset int) ([]Post, error) {
	posts := []Post{}

	err := DBConn.Select(&posts, "SELECT * FROM `posts` WHERE created_at > ? ORDER BY `score` LIMIT 100 OFFSET ?", cutoff.Format("2006-01-02 15:04:05"), offset)

	return posts, err
}

// GetLatestPosts - get 100 posts ordered by creation date (newest first) and
// offset by passed in value
func GetLatestPosts(offset int) ([]Post, error) {
	posts := []Post{}

	err := DBConn.Select(&posts, "SELECT * FROM `posts` ORDER BY `created_at` DESC, `id` DESC LIMIT 100 OFFSET ?", offset)

	return posts, err
}
