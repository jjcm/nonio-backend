package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	Tags      []PostTag `db:"-"`
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
		Title     string    `json:"title"`
		UserName  string    `json:"user"`
		TimeStamp int64     `json:"time"`
		URL       string    `json:"url"`
		Content   string    `json:"content"`
		Type      string    `json:"type"`
		Score     int       `json:"score"`
		Tags      []PostTag `json:"tags"`
	}{
		Title:     p.Title,
		UserName:  p.Author.GetDisplayName(),
		TimeStamp: p.CreatedAt.UnixNano() / int64(time.Millisecond),
		URL:       p.URL,
		Content:   p.Content,
		Type:      p.Type,
		Score:     p.Score,
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
	if err != nil {
		return err
	}
	return nil
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
	err := DBConn.Select(&postTags, "SELECT * FROM posts_tags where post_id = ?)", p.ID)
	if err != nil {
		return err
	}
	p.Tags = postTags
	return nil
}

// GetPostsByScoreSince - get posts from the database that have been created since
// the provided cutoff time, and offset the results
func GetPostsByScoreSince(cutoff time.Time, offset int) ([]Post, error) {
	posts := []Post{}

	err := DBConn.Select(&posts, "SELECT * FROM `posts` WHERE created_at > ? ORDER BY `score` DESC LIMIT 100 OFFSET ?", cutoff.Format("2006-01-02 15:04:05"), offset)

	return posts, err
}

// GetLatestPosts - get 100 posts ordered by creation date (newest first) and
// offset by passed in value
func GetLatestPosts(offset int) ([]Post, error) {
	posts := []Post{}

	err := DBConn.Select(&posts, "SELECT * FROM `posts` ORDER BY `created_at` DESC, `id` DESC LIMIT 100 OFFSET ?", offset)

	return posts, err
}

// Comments will return comments associated with the current post
// This method has gone back and forth a bit, and
func (p *Post) Comments(depthLimit int) ([]Comment, error) {
	var err error
	var comments []Comment

	// this is a temporary work around to let front end dev get back at it...
	err = DBConn.Select(&comments, "SELECT * FROM comments WHERE post_id = ?", p.ID)
	return comments, err

	// everything below this line won't get run because of the above return statement,
	// but it's here to pick up where I left off... :-/

	// we're going to run this query X times, where X = depthLimit
	query := "SELECT id FROM comments WHERE post_id = ? and parent_id IN (?)"
	var commentIDs []string
	parentIDs := []string{"0"}

	for index := 0; index < depthLimit; index++ {
		fmt.Println(query, p.ID, strings.Join(parentIDs, ","))
		rows, err := DBConn.Query(query, p.ID, strings.Join(parentIDs, ","))
		if err != nil {
			return comments, err
		}
		for rows.Next() {
			var id int
			rows.Scan(&id)
			fmt.Println(id)
			parentIDs = append(parentIDs, strconv.Itoa(id))
		}
		rows.Close()
	}
	fmt.Println(commentIDs)
	return comments, err

	for depth := 0; depth < depthLimit; depth++ {
		parentIDs := getUniqueCommentParentIDs(comments)
		fmt.Println(parentIDs)
		// prepare for the next loop
		// run the query
		rows, err := DBConn.Query("SELECT id, author_id, post_id, created_at, type, content, text, parent_id FROM `comments` WHERE post_id = ? AND parent_id in (?)", p.ID, parentIDs)
		for rows.Next() {
			c := Comment{}
			err = rows.Scan(
				&(c.ID),
				&(c.AuthorID),
				&(c.PostID),
				&(c.CreatedAt),
				&(c.Type),
				&(c.Content),
				&(c.Text),
				&(c.ParentID),
			)
			if err != nil {
				return comments, err
			}
			comments = append(comments, c)
		}
	}
	return comments, err
}
