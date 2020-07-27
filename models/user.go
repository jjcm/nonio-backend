package models

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User this is a standard struct that represents a user in the system
type User struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Username  string    `db:"username" json:"username"`
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

// FindByUsername find a user by searching the DB
func (u *User) FindByUsername(email string) error {
	dbUser := User{}
	err := DBConn.Get(&dbUser, "SELECT * FROM users WHERE username = ?", email)
	if err != nil {
		return err
	}

	// if a record was found, then let's hydrate the current User struct with
	// the found one
	if dbUser.ID != 0 {
		*u = dbUser
	}
	return nil
}

// FindByID find a user by searching the DB
func (u *User) FindByID(id int) error {
	dbUser := User{}
	if id == 0 {
		dbUser.Username = "Anonymous coward"
		*u = dbUser
		return nil
	}
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

// CommentOnPost will try and create a comment in the database
func (u *User) CommentOnPost(post Post, parent *Comment, content string) (Comment, error) {
	c := Comment{}
	now := time.Now().Format("2006-01-02 15:04:05")

	if u.ID == 0 || post.ID == 0 {
		return c, errors.New("Can't create a comment for an invalid user or post")
	}

	var commentParentID int
	if parent != nil {
		commentParentID = parent.ID
	}

	result, err := DBConn.Exec("INSERT INTO comments (author_id, post_id, created_at, content, parent_id) VALUES (?, ?, ?, ?, ?)", u.ID, post.ID, now, content, commentParentID)
	if err != nil {
		return c, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return c, err
	}

	c.FindByID(int(insertID))
	return c, err
}

// AbandonComment removes the user from the comment, but leaves the content
func (u *User) AbandonComment(comment *Comment) error {
	if u.ID == 0 || comment.ID == 0 {
		return errors.New("Can't abandon a comment for an invalid user or comment")
	}

	_, err := DBConn.Exec("UPDATE comments SET author_id = NULL WHERE id = ?", comment.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment removes it from the db
func (u *User) DeleteComment(comment *Comment) error {
	if u.ID == 0 || comment.ID == 0 {
		return errors.New("Can't delete a comment for an invalid user or comment")
	}

	_, err := DBConn.Exec("DELETE FROM comments WHERE id = ?", comment.ID)
	if err != nil {
		return err
	}

	return nil
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
func createUser(email, username, password string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	hashedPassword, err := hashPassword(password)
	_, err = DBConn.Exec("INSERT INTO users (email, username, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", email, username, hashedPassword, now, now)
	if err != nil {
		return err
	}
	return nil
}

// UserFactory will create and return an instance of a user
func UserFactory(email, username, password string) (User, error) {
	u := User{}
	err := createUser(email, username, password)
	if err != nil {
		return u, err
	}
	err = u.FindByEmail(email)

	return u, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ChangePassword changes the password of the user, assuming correct args
func (u *User) ChangePassword(oldPassword string, newPassword string, confirmPassword string) error {
	// Check if both the new password and its confirmation password matches
	if newPassword != confirmPassword {
		return errors.New("New password and confirmation do not match")
	}

	// Check if the password has the required amount of entropy. In this case the min is 2^40 combinations
	const minEntropy float64 = 40
	if getEntropy(newPassword) < minEntropy {
		return errors.New("New password does not meet the entropy requirement")
	}

	// Make sure the old password isn't incorrect
	if !checkPasswordHash(oldPassword, u.Password) {
		return errors.New("Old password is incorrect")
	}

	// Generate a hash from the new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return errors.New("Error hashing password")
	}

	// If the checks look good, change the password
	_, err = DBConn.Exec("UPDATE users set password = ? where ID = ?", hashedPassword, u.ID)

	return err
}

// getEntropy returns the log base 2 of the entropy of a password
func getEntropy(password string) float64 {
	type entropyDictionary struct {
		lowercase bool
		uppercase bool
		numbers   bool
		specials  bool
	}
	var charsets entropyDictionary
	var lowercaseRe = regexp.MustCompile(`[a-z]`).MatchString
	var uppercaseRe = regexp.MustCompile(`[A-Z]`).MatchString
	var numbersRe = regexp.MustCompile(`[0-9]`).MatchString
	for _, char := range password {
		if lowercaseRe(string(char)) {
			charsets.lowercase = true
		} else if uppercaseRe(string(char)) {
			charsets.uppercase = true
		} else if numbersRe(string(char)) {
			charsets.numbers = true
		} else {
			charsets.specials = true
		}
	}

	var entropyBase = 0
	if charsets.lowercase {
		entropyBase += 26
	}
	if charsets.uppercase {
		entropyBase += 26
	}
	if charsets.numbers {
		entropyBase += 10
	}
	if charsets.specials {
		entropyBase += 50
	}

	var entropy = math.Log2(math.Pow(float64(entropyBase), float64(len(password))))
	return entropy
}

// UsernameIsAvailable - check the database to see if a certian username is
// already taken
func UsernameIsAvailable(username string) (bool, error) {
	var total int
	err := DBConn.Get(&total, "SELECT COUNT(*) FROM users WHERE username = ?", username)
	if err != nil {
		return false, err
	}
	if total != 0 {
		return false, nil
	}
	return true, nil
}

// GetDisplayName - return a string that shows the user's preferred display name
func (u *User) GetDisplayName() string {
	if u.Username != "" {
		return u.Username
	}
	if u.Name != "" {
		return u.Name
	}
	if u.Email != "" {
		return u.Email
	}
	return "User" + strconv.Itoa(u.ID)
}

// MyPosts will return all Posts that were authored by this user.
// limit and offset will adjust the SQL query to return a smaller subset
// pass -1 for limit and you can return all
func (u *User) MyPosts(limit, offset int) ([]Post, error) {
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

// MyVotes will return every posttag the user has voted on.
func (u *User) MyVotes() ([]PostTagVote, error) {
	votes := []PostTagVote{}

	// run the correct sql query
	var query = "SELECT * FROM posts_tags_votes WHERE voter_id = ?"
	err := DBConn.Select(&votes, query, u.ID)
	if err != nil {
		return votes, err
	}

	return votes, nil
}

// MySubscriptions will return the user's tag subscriptions
func (u *User) MySubscriptions() ([]Subscription, error) {
	subscriptions := []Subscription{}

	// run the correct sql query
	var query = "SELECT * FROM subscriptions WHERE user_id = ?"
	err := DBConn.Select(&subscriptions, query, u.ID)
	if err != nil {
		return subscriptions, err
	}

	return subscriptions, nil
}
