package models

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// Post - struct representation of a single post
type Post struct {
	ID        int       `db:"id" json:"ID"`
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
	Tags      []PostTag `db:"-"`
	Width     int       `db:"width" json:"width"`
	Height    int       `db:"height" json:"height"`
}

// PostQueryParams - structure represents the parameters for querying posts
type PostQueryParams struct {
	TagID  int
	Since  string
	Offset int
	UserID int
	Sort   string
}

// MarshalJSON custom JSON builder for Post structs
func (p *Post) MarshalJSON() ([]byte, error) {
	// build tag array for JS if the tag list is currently empty
	if len(p.Tags) < 1 {
		p.getPostTags()
	}

	// populate user if it currently isn't hydrated
	if p.Author.ID == 0 {
		p.Author.FindByID(p.AuthorID)
	}

	// return the custom JSON for this post
	return json.Marshal(&struct {
		ID        int       `json:"ID"`
		Title     string    `json:"title"`
		UserName  string    `json:"user"`
		TimeStamp int64     `json:"time"`
		URL       string    `json:"url"`
		Content   string    `json:"content"`
		Type      string    `json:"type"`
		Score     int       `json:"score"`
		Tags      []PostTag `json:"tags"`
		Width     int       `json:"width"`
		Height    int       `json:"height"`
	}{
		ID:        p.ID,
		Title:     p.Title,
		UserName:  p.Author.GetDisplayName(),
		TimeStamp: p.CreatedAt.UnixNano() / int64(time.Millisecond),
		URL:       p.URL,
		Content:   p.Content,
		Type:      p.Type,
		Score:     p.Score,
		Tags:      p.Tags,
		Width:     p.Width,
		Height:    p.Height,
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

// FindByID - find a given post in the database by its primary key
func (p *Post) FindByID(id int) error {
	dbPost := Post{}
	err := DBConn.Get(&dbPost, "SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		return err
	}

	*p = dbPost
	return nil
}

// IncrementScore - increment the score by post id
func (p *Post) IncrementScore(id int) error {
	_, err := DBConn.Exec("update posts set score=score+1 where id = ?", id)
	return err
}

// IncrementScoreWithTx - increment the score by post id
func (p *Post) IncrementScoreWithTx(tx Transaction, id int) error {
	_, err := tx.Exec("update posts set score=score+1 where id = ?", id)
	return err
}

// DecrementScoreWithTx - decrement the score by post id
func (p *Post) DecrementScoreWithTx(tx Transaction, id int) error {
	_, err := tx.Exec("update posts set score=score-1 where id = ?", id)
	return err
}

// GetCreatedAtTimestamp - get the created at timestamp in the predetermined format
func (p *Post) GetCreatedAtTimestamp() int64 {
	return p.CreatedAt.UnixNano() / int64(time.Millisecond)
}

// AddTag - associate a post with an existing tag
func (p *Post) AddTag(t Tag) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	// create a new tag association
	_, err := DBConn.Exec("INSERT INTO posts_tags (post_id, tag_id, score, created_at) VALUES (?, ?, 1, ?)", p.ID, t.ID, now)
	if err != nil {
		return err
	}

	// get the post tags from 'posts_tags'
	err = p.getPostTags()
	if err != nil {
		return err
	}
	return nil
}

// get the post tags
func (p *Post) getPostTags() error {
	postTags := []PostTag{}
	err := DBConn.Select(&postTags, "SELECT * FROM posts_tags where post_id = ?", p.ID)
	if err != nil {
		return err
	}
	p.Tags = postTags
	return nil
}

// GetPostTags - get the tags with post id
func GetPostTags(id int) ([]PostTag, error) {
	tags := []PostTag{}

	err := DBConn.Select(&tags, "SELECT * FROM posts_tags where post_id = ?", id)
	if err != nil {
		return nil, err
	}

	// query the tag name with id
	for i, item := range tags {
		tag := Tag{}
		if err = tag.FindByID(item.TagID); err != nil {
			return nil, err
		}
		tags[i].TagName = tag.Name
	}

	return tags, err
}

func intSlice2Str(vals []int) string {
	var s []string

	for _, v := range vals {
		s = append(s, strconv.Itoa(v))
	}
	return strings.TrimSuffix(strings.Join(s, ","), ",")
}

// GetPostsByParams - get the posts by parameters
func GetPostsByParams(params *PostQueryParams) ([]*Post, error) {
	args := []interface{}{}

	query := "select * from posts where created_at > ?"
	// time range
	args = append(args, params.Since)

	// special user
	if params.UserID > 0 {
		query = query + " and user_id = ?"
		args = append(args, params.UserID)
	}

	// tags
	if params.TagID > 0 {
		query = query + " and id in (SELECT post_id from posts_tags where tag_id = ?)"
		args = append(args, params.TagID)
	}

	// orders
	if params.Sort == "popular" || params.Sort == "top" {
		query = query + " order by score desc"
	}
	if params.Sort == "new" {
		query = query + " order by created_at desc"
	}

	// offset
	query = query + " limit 100 offset ?"
	args = append(args, params.Offset)

	posts := []*Post{}
	// exec the query string
	if err := DBConn.Select(&posts, query, args...); err != nil {
		return nil, err
	}

	return posts, nil
}

// Comments will return comments associated with the current post
// This method has gone back and forth a bit, and
func (p *Post) Comments(depthLimit int) ([]Comment, error) {
	var err error
	var comments []Comment

	// this is a temporary work around to let front end dev get back at it...
	err = DBConn.Select(&comments, "SELECT * FROM comments WHERE post_id = ?", p.ID)
	return comments, err
}
