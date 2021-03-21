package models

import (
	"testing"
)

func TestWeCanVoteOnComments(t *testing.T) {
	setupTestingDB()

	user1, _ := UserFactory("example1@example.com", "bob", "password", 0)
	user2, _ := UserFactory("example2@example.com", "ralph", "password", 0)

	post, _ := user1.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)

	comment1, _ := user2.CreateComment(post, nil, "Test comment from user 2 on user 1's post")
	comment2, _ := user2.CreateComment(post, &comment1, "Test comment replying to comment 1")
	comment3, _ := user2.CreateComment(post, &comment2, "Test comment replying to comment 2")

	if err := WithTransaction(func(tx Transaction) error {
		err := user1.CreateCommentVoteWithTx(tx, comment1.ID, true)
		err = user1.CreateCommentVoteWithTx(tx, comment2.ID, true)
		err = user1.CreateCommentVoteWithTx(tx, comment3.ID, true)
		return err
	}); err != nil {
		t.Errorf("Couldn't perform a transactional update of the comment votes: %v\n", err)
	}

	comment1.FindByID(comment1.ID)
	if comment1.LineageScore != 2 {
		t.Errorf("Comment's lineage score should be 2. Got %v instead.\n", comment1.LineageScore)
	}
	/*
		if err != nil {
			t.Errorf("We should have been able to create a comment. Error recieved: %s", err)
		}
		if commentVote.ID == 0 {
			t.Error("The created comment should have been instantiated correctly.")
		}
	*/
}
