package models

import (
	"testing"
)

func TestWeCanCreateComments(t *testing.T) {
	setupTestingDB()

	author, _ := UserFactory("example@example.com", "", "password")

	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	// create comment
	comment, err := author.CommentOnPost(post, nil, "This is a dumb post")

	if err != nil {
		t.Errorf("We should have been able to create a comment. Error recieved: %s", err)
	}
	if comment.ID == 0 {
		t.Error("The created comment should have been instantiated correctly.")
	}
}

func TestWeCanGetCommentsForAPost(t *testing.T) {
	setupTestingDB()

	author, _ := UserFactory("example@example.com", "", "password")

	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	// create comment
	comment, _ := author.CommentOnPost(post, nil, "This is a dumb post")
	author.CommentOnPost(post, &comment, "This is a reply")

	comments, err := GetCommentsByPost(post.ID)

	if err != nil {
		t.Errorf("We should have been able to get comments. Error recieved: %s", err)
	}
	if len(comments) != 2 {
		t.Errorf("Should have found two comments. Instead recieved: %v", len(comments))

	}
}
