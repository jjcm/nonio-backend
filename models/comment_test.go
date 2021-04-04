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

func TestWeCanGetCommentsByParams(t *testing.T) {
	setupTestingDB()

	bill, _ := UserFactory("bill@example.com", "bill", "password", 0)
	joe, _ := UserFactory("joe@example.com", "joe", "password", 0)

	post1, _ := bill.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)
	post2, _ := bill.CreatePost("Post Title", "post-title-2", "lorem ipsum", "image", 0, 0)

	// create comments on the first post
	comment1, _ := joe.CreateComment(post1, nil, "This is a dumb post")
	bill.CreateComment(post1, &comment1, "This is a reply")

	// create comment on the second post
	joe.CreateComment(post2, nil, "This is even worse")

	commentQueryParams := CommentQueryParams{}
	var params *CommentQueryParams = &commentQueryParams

	// try getting all comments
	comments, err := GetCommentsByParams(params)
	if err != nil {
		t.Errorf("We should have been able to get comments. Error recieved: %s", err)
	}
	if len(comments) != 3 {
		t.Errorf("Should have found three comments. Instead recieved: %v", len(comments))
	}

	// try getting all of Joe's comments
	params.UserID = joe.ID
	comments, err = GetCommentsByParams(params)
	if err != nil {
		t.Errorf("We should have been able to get comments. Error recieved: %s", err)
	}
	if len(comments) != 2 {
		t.Errorf("Should have found two comments from Joe. Instead recieved: %v", len(comments))
	}

	// try getting all of Joe's comments, but only for the first post
	params.PostID = post1.ID
	comments, err = GetCommentsByParams(params)
	if err != nil {
		t.Errorf("We should have been able to get comments. Error recieved: %s", err)
	}
	if len(comments) != 1 {
		t.Errorf("Should have found one comment from Joe. Instead recieved: %v", len(comments))
	}

	// try getting all of post 1's comments
	params.UserID = 0
	comments, err = GetCommentsByParams(params)
	if err != nil {
		t.Errorf("We should have been able to get comments. Error recieved: %s", err)
	}
	if len(comments) != 2 {
		t.Errorf("Should have found two comments. Instead recieved: %v", len(comments))
	}
}
