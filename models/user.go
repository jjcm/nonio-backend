package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User this is a standard struct that represents a user in the system
type User struct {
	ID        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
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

	if dbUser.Email == "" {
		return errors.New("User not found")
	}

	// if an email was found, then let's hydrate the current User struct with
	// the found one
	*u = dbUser
	return nil
}

// Login a user if their password matches the stored hash
func (u *User) Login(password string) error {
	if !checkPasswordHash(password, u.Password) {
		return errors.New("Not a match")
	}
	return nil
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
