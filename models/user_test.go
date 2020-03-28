package models

import (
	"strconv"
	"testing"
)

func TestCanCreateAUser(t *testing.T) {
	setupTestingDB()

	err := CreateUser("user@example.com", "", "password")
	if err != nil {
		t.Errorf("Creating a user should not have errors. Error: " + err.Error())
	}
}

func TestAUsersPasswordCanBeChecked(t *testing.T) {
	setupTestingDB()

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

func TestFindingAUserByTheirUsername(t *testing.T) {
	setupTestingDB()

	u := User{}
	err := u.FindByUsername("")
	if err == nil {
		t.Errorf("Searching for a user by an empty username should have thrown an error")
	}

	// now let's create a user and search with an invalid username, we are expecting an error
	CreateUser("user@example.com", "radUser123", "anything")
	err = u.FindByUsername("radUser")
	if err == nil {
		t.Errorf("The username 'radUser' should not exist in the database so an error should have been thrown")
	}
	err = u.FindByUsername("radUser123")
	if err != nil {
		t.Errorf("The username 'radUser123' should exist in the database so an error should not have been thrown")
	}
	if u.ID == 0 {
		t.Error("We should have hyradted the user struct, but it's not hydrated")
	}
}

func TestFindingAUserByTheirPrimaryKeyID(t *testing.T) {
	setupTestingDB()

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

func TestWeCanCheckIfAUsernameIsAvailable(t *testing.T) {
	setupTestingDB()

	isAvaiable, _ := UsernameIsAvailable("anything")
	if !isAvaiable {
		t.Errorf("Because the database is empty, any username should be available")
	}

	// create our user with a username of "anything"
	CreateUser("user@example.com", "anything", "password")

	// now the username "anything" shouldn't be available
	isAvaiable, _ = UsernameIsAvailable("anything")
	if isAvaiable {
		t.Errorf("The username should be taken, but the system said it's available")
	}
}

func TestWeCanGetTheUsersPreferredDisplayName(t *testing.T) {
	u := User{
		ID:       88,
		Email:    "user@example.com",
		Username: "sociuser",
		Name:     "Soci User",
	}

	// default is to show the user by their username
	expected := "sociuser"
	if u.GetDisplayName() != expected {
		t.Errorf("If the username isn't an empty string, then it should be the value returned.\nExpected: %s\n  Actual: %s\n", expected, u.GetDisplayName())
	}

	u.Username = ""
	expected = "Soci User"
	if u.GetDisplayName() != expected {
		t.Errorf("If the username is an empty string, then it should be the name that is returned.\nExpected: %s\n  Actual: %s\n", expected, u.GetDisplayName())
	}

	u.Name = ""
	expected = "user@example.com"
	if u.GetDisplayName() != expected {
		t.Errorf("If the username and name are both empty, then the user's email should be returned.\nExpected: %s\n  Actual: %s\n", expected, u.GetDisplayName())
	}

	u.Email = ""
	expected = "User88"
	if u.GetDisplayName() != expected {
		t.Errorf("If the username, name, and email address are all empty, then we need to have some sort of fallback. Why not their Primary Key?\nExpected: %s\n  Actual: %s\n", expected, u.GetDisplayName())
	}
}

func TestWeCanGetAllThePostsFromAUser(t *testing.T) {
	setupTestingDB()

	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	// the author shouldn't have any posts at this point
	posts, err := author.MyPosts(-1, 0)
	if len(posts) != 0 {
		t.Errorf("Expected the User to not have any posts")
	}
	if err != nil {
		t.Error(err)
	}

	// create a post
	author.CreatePost("It's my post! Yay!", "post-title", "lorem ipsum", "image")
	posts, err = author.MyPosts(-1, 0)
	if err != nil {
		t.Error(err)
	}
	if len(posts) == 0 {
		t.Errorf("Expected a user to have posts")
		t.Errorf("%v", posts)
	}

	// now let's create a ton of posts to test the offset and limit behavior
	for index := 1; index < 200; index++ {
		indexAsString := strconv.Itoa(index)
		author.CreatePost("Post Title "+indexAsString, "post-title-"+indexAsString, "lorem ipsum", "image")
	}

	posts, err = author.MyPosts(-1, 0)
	if err != nil {
		t.Error(err)
	}
	if len(posts) != 200 {
		t.Errorf("Expected a user to have 200 posts:\nExpected: %v\n  Actual: %v", 200, len(posts))
		return
	}

	limitedPosts, err := author.MyPosts(100, 0)
	if err != nil {
		t.Error(err)
	}
	if limitedPosts[0].Title != "It's my post! Yay!" {
		t.Errorf("Expected the first Post to be the one we created at the beginning of this test")
	}
	if len(limitedPosts) != 100 {
		t.Errorf("Expected that we could limit the total number of Posts returned.\nExpected: %v\n  Actual: %v", 100, len(limitedPosts))
	}
}
