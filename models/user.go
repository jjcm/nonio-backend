package models

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"soci-backend/httpd/utils"
	"strconv"
	"time"

	b64 "encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

// User this is a standard struct that represents a user in the system
type User struct {
	ID                 int       `db:"id" json:"id"`
	Email              string    `db:"email" json:"email"`
	Username           string    `db:"username" json:"username"`
	Name               string    `db:"name" json:"name"`
	Password           string    `db:"password" json:"password"`
	StripeCustomerID   string    `db:"stripe_customer_id" json:"stripe_customer_id"`
	Description        string    `db:"description" json:"description"`
	SubscriptionAmount float64   `db:"subscription_amount" json:"subscriptionAmount"`
	Cash               float64   `db:"cash" json:"cash"`
	LastLogin          time.Time `db:"last_login" json:"-"`
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time `db:"updated_at" json:"updatedAt"`
}

type TempPassword struct {
	ID                 int       `db:"id" json:"-"`
	Email              string    `db:"email" json:"-"`
	TempPassword       string    `db:"temp_password" json:"-"`
	TempPasswordExpiry time.Time `db:"temp_password_expiry" json:"-"`
}

/************************************************/
/******************** CREATE ********************/
/************************************************/

// createUser try and create a new user
func createUser(email, username, password string, subscriptionAmount float64) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	_, err = DBConn.Exec("INSERT INTO users (email, username, password, subscription_amount, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)", email, username, hashedPassword, subscriptionAmount, now, now)
	if err != nil {
		return err
	}
	return nil
}

// UserFactory will create and return an instance of a user
func UserFactory(email, username, password string, subscriptionAmount float64) (User, error) {
	u := User{}
	err := createUser(email, username, password, subscriptionAmount)
	if err != nil {
		return u, err
	}
	err = u.FindByEmail(email)

	return u, err
}

/************************************************/
/********************* READ *********************/
/************************************************/

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
func (u *User) FindByUsername(username string) error {
	dbUser := User{}
	err := DBConn.Get(&dbUser, "SELECT * FROM users WHERE username = ?", username)
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

// GetAll gets all users.
func (u *User) GetAll() ([]User, error) {
	users := []User{}
	err := DBConn.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	// if an email was found, then let's hydrate the current User struct with
	// the found one
	return users, nil
}

type UserFinancialData struct {
	SubscriptionAmount float64 `db:"subscription_amount" json:"subscription_amount"`
	Cash               float64 `db:"cash" json:"cash"`
	StripeCustomerID   float64 `db:"stripe_customer_id" json:"stripe_customer_id"`
}

// GetFinancialData will return the user's tag subscriptions
func (u *User) GetFinancialData() (UserFinancialData, error) {
	financialData := UserFinancialData{}

	// run the correct sql query
	var query = "SELECT cash, subscription_amount, stripe_customer_id FROM users WHERE id = ?"
	err := DBConn.Get(&financialData, query, u.ID)
	if err != nil {
		return financialData, err
	}

	return financialData, nil
}

type UserInfo struct {
	Description  string `db:"description" json:"description"`
	Posts        int    `db:"-" json:"posts"`
	Karma        int    `db:"-" json:"karma"`
	Comments     int    `db:"-" json:"comments"`
	CommentKarma int    `db:"-" json:"comment_karma"`
}

// GetFinancialData will return the user's tag subscriptions
func (u *User) GetInfo() (UserInfo, error) {
	userInfo := UserInfo{}

	// get our description first
	var query = "SELECT description FROM users WHERE id = ?"
	err := DBConn.Get(&userInfo, query, u.ID)
	if err != nil {
		return userInfo, err
	}

	// get the number of posts for the user
	err = DBConn.Get(&userInfo.Posts, "SELECT COUNT(*) FROM posts WHERE user_id = ?", u.ID)
	if err != nil {
		return userInfo, err
	}

	// get the karma for the user
	err = DBConn.Get(&userInfo.Karma, "SELECT SUM(score) FROM posts WHERE user_id = ?", u.ID)
	if err != nil {
		return userInfo, err
	}

	// get the number of comments for the user
	err = DBConn.Get(&userInfo.Comments, "SELECT COUNT(*) FROM comments WHERE author_id = ?", u.ID)
	if err != nil {
		return userInfo, err
	}

	// get the karma for the user
	err = DBConn.Get(&userInfo.CommentKarma, "SELECT SUM(upvotes - downvotes) FROM comments WHERE author_id = ?", u.ID)
	if err != nil {
		if err.Error() == "sql: Scan error on column index 0, name \"SUM(upvotes - downvotes)\": converting NULL to int is unsupported" {
			// they don't yet have any comments, so their comment karma should just be 0
			userInfo.CommentKarma = 0
		} else {
			return userInfo, err
		}
	}

	return userInfo, nil
}

/************************************************/
/******************** UPDATE ********************/
/************************************************/

// ChangePassword changes the password of the user, assuming correct args
func (u *User) ChangePassword(oldPassword string, newPassword string, confirmPassword string) error {
	// Check if both the new password and its confirmation password matches
	if newPassword != confirmPassword {
		return errors.New("new password and confirmation do not match")
	}

	// Check if the password has the required amount of entropy. In this case the min is 2^40 combinations
	const minEntropy float64 = 40
	passwordEntropy := getEntropy(newPassword)
	if passwordEntropy < minEntropy {
		return fmt.Errorf("new password does not meet the entropy requirement. Password entropy: %v. Required: %v. Password: %v", passwordEntropy, minEntropy, newPassword)
	}

	// Make sure the old password isn't incorrect
	if !checkPasswordHash(oldPassword, u.Password) {
		return errors.New("old password is incorrect")
	}

	// Generate a hash from the new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return errors.New("error hashing password")
	}

	// If the checks look good, change the password
	u.Password = hashedPassword
	err = u.update()

	return err
}

// UpdateDescription updates the description for the user
func (u *User) UpdateDescription(description string) error {
	_, err := DBConn.Exec("UPDATE users SET description = ? WHERE id = ?", description, u.ID)

	return err
}

// UpdateStripeCustomerID updates the stripe customer id for the user
func (u *User) UpdateStripeCustomerID(id string) error {
	_, err := DBConn.Exec("UPDATE users SET stripe_customer_id = ? WHERE id = ?", id, u.ID)

	return err
}

// ChangeForgottenPassword changes the password of the user, using an emailed token as verification
func (u *User) ChangeForgottenPassword(token string, newPassword string, confirmPassword string) error {
	// Check if both the new password and its confirmation password matches
	if newPassword != confirmPassword {
		return errors.New("new password and confirmation do not match")
	}

	// Check if the password has the required amount of entropy. In this case the min is 2^40 combinations
	const minEntropy float64 = 40
	passwordEntropy := getEntropy(newPassword)
	if passwordEntropy < minEntropy {
		return fmt.Errorf("new password does not meet the entropy requirement. Password entropy: %v. Required: %v. Password: %v", passwordEntropy, minEntropy, newPassword)
	}

	// Check everything is kosher
	tempPassword := TempPassword{}
	err := DBConn.Get(&tempPassword, "SELECT * from user_temp_passwords where temp_password = ?", token)
	if err != nil {
		return errors.New("error retrieving temp password from DB")
	}
	if tempPassword.ID == 0 {
		return errors.New("token not found in db")
	}
	if tempPassword.TempPasswordExpiry.Before(time.Now()) {
		return errors.New("token found but expired")
	}

	// If all is good, then find our user and update their password
	u.FindByEmail(tempPassword.Email)

	// Generate a hash from the new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return errors.New("error hashing password")
	}

	// If the checks look good, change the password
	u.Password = hashedPassword
	err = u.update()
	if err != nil {
		Log.Errorf("Error updating user password: %v\n", err)
		return err
	}

	_, err = DBConn.Exec("DELETE FROM user_temp_passwords WHERE email = ?", u.Email)
	if err != nil {
		Log.Errorf("Error removing previous forgot password request: %v\n", err)
		return err
	}

	return err
}

/************************************************/
/******************** DELETE ********************/
/************************************************/

// TODO

/************************************************/
/******************** HELPER ********************/
/************************************************/

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

// Login a user if their password matches the stored hash
func (u *User) Login(password string) error {
	if !checkPasswordHash(password, u.Password) {
		return errors.New("incorrect password")
	}
	u.LastLogin = time.Now()
	err := u.update()
	return err
}

func (u *User) update() error {
	_, err := DBConn.Exec("UPDATE users SET email = ?, name = ?, last_login = ?, password = ?, updated_at = ? WHERE id = ?", u.Email, u.Name, u.LastLogin, u.Password, time.Now(), u.ID)
	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

func (u *User) ForgotPasswordRequest(email string) error {
	user := User{}
	user.FindByEmail(email)

	// If the email isn't associated with an account, do nothing
	if user.ID == 0 {
		return nil
	}

	token := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
	rand.Read(token)
	encodedToken := b64.StdEncoding.EncodeToString(token)

	// Remove the + / = from the string to keep it URL safe
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		Log.Errorf("Regex failed to compile: %v\n", err)
		return err
	}

	encodedToken = re.ReplaceAllString(encodedToken, "")
	expiry := time.Now().Add(time.Hour).Format("2006-01-02 15:04:05")

	// delete any previous request
	_, err = DBConn.Exec("DELETE FROM user_temp_passwords WHERE email = ?", email)
	if err != nil {
		Log.Errorf("Error removing previous forgot password request: %v\n", err)
		return err
	}

	_, err = DBConn.Exec("INSERT INTO user_temp_passwords (email, temp_password, temp_password_expiry) VALUES (?, ?, ?)", email, encodedToken, expiry)
	if err != nil {
		Log.Errorf("Error inserting temp password: %v\n", err)
		return err
	}

	// TODO - make the host an environment variable
	utils.SendEmailOAUTH2(
		email,
		"Nonio - Forgot Password Request",
		fmt.Sprintf("Please click the following link to set a new password: %v/admin/change-forgotten-password?token=%v", WebHost, encodedToken),
	)

	return nil
}
