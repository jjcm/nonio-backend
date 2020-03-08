package models

import "testing"

func TestWeCanCheckIfAUrlIsAvailable(t *testing.T) {
	setupTestingDB()

	isAvaiable, _ := URLIsAvailable("anything")
	if !isAvaiable {
		t.Errorf("Because the database is empty, any url should be available")
	}

	// cool, now let's create a post and the URL should not be available anymore
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	p, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	// now the URL for the existing post should be taken
	isAvaiable, _ = URLIsAvailable(p.URL)
	if isAvaiable {
		t.Errorf("The URL for a built post should be taken, but the system said it's available")
	}
}
