package models

import (
	"testing"
)

func TestWeCanCreateComments(t *testing.T) {
	setupTestingDB()

	author, _ := UserFactory("example@example.com", "", "password", 0)

	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)

	// create comment
	comment, err := author.CreateComment(post, nil, "This is a dumb post")

	if err != nil {
		t.Errorf("We should have been able to create a comment. Error recieved: %s", err)
	}
	if comment.ID == 0 {
		t.Error("The created comment should have been instantiated correctly.")
	}
}

func TestWeCanGetCommentsForAPost(t *testing.T) {
	setupTestingDB()

	author, _ := UserFactory("example@example.com", "", "password", 0)

	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)

	// create comment
	comment, _ := author.CreateComment(post, nil, "This is a dumb post")
	author.CreateComment(post, &comment, "This is a reply")

	comments, err := GetCommentsByPost(post.ID)

	if err != nil {
		t.Errorf("We should have been able to get comments. Error recieved: %s", err)
	}
	if len(comments) != 2 {
		t.Errorf("Should have found two comments. Instead recieved: %v", len(comments))

	}
}

func TestWeCanDeleteComments(t *testing.T) {
	setupTestingDB()

	author, _ := UserFactory("example@example.com", "", "password", 0)

	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)

	// create comment
	comment, _ := author.CreateComment(post, nil, "This is a dumb post")
	err := author.DeleteComment(&comment)

	if err != nil {
		t.Errorf("We should have been able to delete a comment. Error recieved: %s", err)
	}
}

func TestWeCanAbandonComments(t *testing.T) {
	setupTestingDB()

	author, _ := UserFactory("example@example.com", "", "password", 0)

	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)

	// create comment
	comment, _ := author.CreateComment(post, nil, "This is a dumb post")
	err := author.AbandonComment(&comment)

	if err != nil {
		t.Errorf("We should have been able to abandon a comment. Error recieved: %s", err)
	}

	// make sure we can get comments for our post still
	_, err2 := GetCommentsByPost(post.ID)
	if err2 != nil {
		t.Errorf("We should have been able to get comments for a post. Error recieved: %s", err)
	}

}
