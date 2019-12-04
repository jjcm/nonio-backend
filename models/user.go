package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User this is a standard struct that represents a user in the system
type User struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	Password  string    `db:"password" json:"password"`
	LastLogin time.Time `db:"last_login" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// FindByEmail find a user by searching the DB
func (u *User) FindByEmail(email string) error {
	dbUser := User{}
	err := DBConn.Get(&dbUser, "SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return err
	}

	// if an email was found, then let's hydrate the current User struct with
	// the found one
	*u = dbUser
	return nil
}

// FindByID find a user by searching the DB
func (u *User) FindByID(id int) error {
	dbUser := User{}
	err := DBConn.Get(&dbUser, "SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	*u = dbUser
	return nil
}

// CreatePost - create a new post in the database and set the current User as
// the author
func (u *User) CreatePost(title, url, content, postType string) (Post, error) {
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
	result, err := DBConn.Exec("INSERT INTO posts (title, url, user_id, thumbnail, score, content, type, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", postTitle, postURL, u.ID, "", 0, content, postType, now, now)
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

// Login a user if their password matches the stored hash
func (u *User) Login(password string) error {
	if !checkPasswordHash(password, u.Password) {
		return errors.New("Not a match")
	}
	u.LastLogin = time.Now()
	err := u.update()
	return err
}

func (u *User) update() error {
	_, err := DBConn.Exec("UPDATE users SET email = ?, name = ?, last_login = ?, updated_at = ? WHERE id = ?", u.Email, u.Name, u.LastLogin, time.Now(), u.ID)
	return err
}

// CreateUser try and create a new user
func CreateUser(email, password string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	hashedPassword, err := hashPassword(password)
	_, err = DBConn.Exec("INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, ?, ?)", email, hashedPassword, now, now)
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
