package models

import (
	"testing"

	"github.com/icrowley/fake"
)

func TestCanGetTags(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	tags, _ := GetTags(0, 100)
	if len(tags) != 0 {
		t.Errorf("There should not be any tags in a fresh DB")
	}

	// create an author for tags
	CreateUser("example@example.com", "password")
	author := User{}
	author.FindByEmail("example@example.com")

	// create a huge number of tags, so we can test our tag retriever with valid limits and offsets
	limit := 500
	index := 0
	for index < limit {
		err := CreateTag(fake.Words(), author)
		// i expect the faker library to return a bunch of duplicates, but I want 500 unique words so we will only increment the index counter if they are all unique
		if err == nil {
			index++
		}
	}

	batch1, err := GetTags(0, 100)
	if err != nil {
		t.Errorf("We should be able to get a batch of tags from the DB, instead we got an error: %v", err.Error())
	}
	if len(batch1) != 100 {
		t.Errorf("We should have only retrieved 100 tags. Instead we retrieved %v", len(batch1))
	}
}

func TestWeCanCreateAndRetrieveATag(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create our author
	CreateUser("example@example.com", "password")
	author := User{}
	author.FindByEmail("example@example.com")

	newTag := fake.Word()
	err := CreateTag(newTag, author)
	if err != nil {
		t.Errorf("You should be able to create a new tag with a given string")
	}

	tag := Tag{}
	err = tag.FindByTagName(newTag)
	if err != nil {
		t.Errorf("We should have been able to find the newly created tag in the DB")
	}
	if tag.Name != newTag {
		t.Errorf("The struct was not hydrated correctly. Expected '%v', got '%v'", newTag, tag.Name)
	}
}

func TestYouCantCreateATagWithTheNameOfAnExistingTag(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create our author
	CreateUser("example@example.com", "password")
	author := User{}
	author.FindByEmail("example@example.com")

	newTag := fake.Word()
	err := CreateTag(newTag, author)
	if err != nil {
		t.Errorf("The initial creation of a tag should have worked fine")
	}
	// try it again
	err = CreateTag(newTag, author)
	if err == nil {
		t.Errorf("This tag should already exist, so an error should have been thrown when trying to create it again")
	}
}
