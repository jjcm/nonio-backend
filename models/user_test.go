package models

import (
	"testing"
)

func TestCanCreateAUser(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	err := CreateUser("user@example.com", "", "password")
	if err != nil {
		t.Errorf("Creating a user should not have errors. Error: " + err.Error())
	}
}

func TestAUsersPasswordCanBeChecked(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	err := CreateUser("user@example.com", "", "password")
	if err != nil {
		panic(err)
	}
	user := User{}
	user.FindByEmail("user@example.com")

	err = user.Login("password")
	if err != nil {
		t.Errorf("The user should have been logged in with the correct password. Error: %v", err)
	}
	err = user.Login("wrongpassword")
	if err == nil {
		t.Errorf("The incorrect password should have thrown an error. Error: " + err.Error())
	}
}

func TestAUserCantBeCreatedIfTheEmailAlreadyExists(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	err := CreateUser("user@example.com", "", "anything")
	if err != nil {
		t.Errorf("Initial creation of a user should work fine")
	}

	// now let's try and create another user with the same email address
	err = CreateUser("user@example.com", "", "anything")
	if err == nil {
		t.Errorf("An error should have been thrown when we tried to create a user with an existing email address")
	}
}

func TestFindingAUserByTheirEmailAddress(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	u := User{}
	err := u.FindByEmail("")
	if err == nil {
		t.Errorf("Searching for a user by an empty email address should have thrown an error")
	}

	// now let's create a user and search with an invalid email address, we are expecting an error
	CreateUser("user@example.com", "", "anything")
	err = u.FindByEmail("wrong@address.com")
	if err == nil {
		t.Errorf("The email address 'wrong@address.com' should not exist in the database so an error should have been thrown")
	}
}

func TestFindingAUserByTheirPrimaryKeyID(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	u := User{}
	err := u.FindByID(1)
	if err == nil {
		t.Errorf("Trying to find by ID before any users exist should have thrown an error")
	}

	// this should be the very first user, so I should be able to find them by the ID 1
	CreateUser("user@example.com", "", "whatever")
	err = u.FindByID(1)
	if err != nil {
		t.Errorf("No error should have been thrown while looking up the newly created user")
	}
}

func TestWeCanTrackTheUsersLastLogin(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	CreateUser("user@example.com", "", "whatever")

	u := User{}
	u.FindByID(1)
	if !u.LastLogin.IsZero() {
		t.Errorf("A freshly created user should have a last login value of 0")
	}

	u.Login("whatever")

	if u.LastLogin.IsZero() {
		t.Errorf("A user that has logged in should have a non-zero timestamp stored in the DB")
	}
}
