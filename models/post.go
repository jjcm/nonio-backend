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

/************************************************/
/******************** CREATE ********************/
/************************************************/

// CreatePost - create a new post in the database
func (u *User) CreatePost(title, url, content, postType string, width int, height int) (Post, error) {
	p := Post{}
	now := time.Now().Format("2006-01-02 15:04:05")

	postURL := url
	// TODO: this needs some backend verification that it's only /^[0-9A-Za-z-_]+$/

	// perpare and truncate the title if necessary
	postTitle := title
	if len(title) > 256 {
		postTitle = title[0:255]
	}

	// set the type if it's blank
	if strings.TrimSpace(postType) == "" {
		postType = "image"
	}

	// try and create the post in the DB
	result, err := DBConn.Exec("INSERT INTO posts (title, url, user_id, thumbnail, score, content, type, width, height, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", postTitle, postURL, u.ID, "", 0, content, postType, width, height, now, now)
	// check for specific error of an existing post URL
	// if we hit this error, run a second DB call to see how many of the posts have a similar URL alias and then tack on a suffix to make this one unique
	if err != nil && err.Error()[0:10] == "Error 1062" {
		var countSimilarAliases int
		// the uniquie URL that we are testing for might not be long enough for the following LIKE query to work, so let's check that here
		urlToCheckFor := postURL
		if len(urlToCheckFor) > 240 {
			urlToCheckFor = urlToCheckFor[0:240] // 240 is safe enough, since this is such an edge case that it's probably never going to happen
		}
		DBConn.Get(&countSimilarAliases, "SELECT COUNT(*) FROM posts WHERE url LIKE ?", urlToCheckFor+"%")
		suffix := "-" + strconv.Itoa(countSimilarAliases+1)
		newURL := postURL + suffix
		if len(newURL) > 255 {
			newURL = postURL[0:len(postURL)-len(suffix)] + suffix
		}

		// now let's try it again with the updated post URL
		result, err = DBConn.Exec("INSERT INTO posts (title, url, user_id, thumbnail, score, content, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", postTitle, newURL, u.ID, "", 0, content, postType, now, now)
		if err != nil {
			return p, err
		}
	}

	if err != nil {
		return p, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return p, err
	}
	p.FindByID(int(insertID))
	p.Author = *u

	return p, nil
}

/************************************************/
/********************* READ *********************/
/************************************************/

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

// GetPosts will return all Posts that were authored by this user.
// limit and offset will adjust the SQL query to return a smaller subset
// pass -1 for limit and you can return all
func (u *User) GetPosts(limit, offset int) ([]Post, error) {
	posts := []Post{}

	// run the correct sql query
	var query string
	if limit == -1 {
		query = "SELECT * FROM posts WHERE user_id = ?"
		err := DBConn.Select(&posts, query, u.ID)
		if err != nil {
			return posts, err
		}
	} else {
		query = "SELECT * FROM posts WHERE user_id = ? LIMIT ? OFFSET ?"
		err := DBConn.Select(&posts, query, u.ID, limit, offset)
		if err != nil {
			return posts, err
		}
	}

	return posts, nil
}

/************************************************/
/******************** UPDATE ********************/
/************************************************/

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

/************************************************/
/******************** DELETE ********************/
/************************************************/
// TODO

/************************************************/
/******************** HELPER ********************/
/************************************************/

// GetCreatedAtTimestamp - get the created at timestamp in the predetermined format
func (p *Post) GetCreatedAtTimestamp() int64 {
	return p.CreatedAt.UnixNano() / int64(time.Millisecond)
}

func intSlice2Str(vals []int) string {
	var s []string

	for _, v := range vals {
		s = append(s, strconv.Itoa(v))
	}
	return strings.TrimSuffix(strings.Join(s, ","), ",")
}
