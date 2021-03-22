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

	if err := WithTransaction(func(tx Transaction) error {
		err := user1.CreateCommentVoteWithTx(tx, comment1.ID, true)
		return err
	}); err != nil {
		t.Errorf("Couldn't perform a transactional update of the comment votes: %v\n", err)
	}
}

func TestWeCanAdjustLineageScore(t *testing.T) {
	setupTestingDB()

	user1, _ := UserFactory("example1@example.com", "bob", "password", 0)
	user2, _ := UserFactory("example2@example.com", "ralph", "password", 0)

	post, _ := user1.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)

	comment1, _ := user2.CreateComment(post, nil, "Test comment from user 2 on user 1's post")
	comment2, _ := user2.CreateComment(post, &comment1, "Test comment replying to comment 1")
	comment3, _ := user2.CreateComment(post, &comment2, "Test comment replying to comment 2")

	if err := WithTransaction(func(tx Transaction) error {
		if err := user1.CreateCommentVoteWithTx(tx, comment1.ID, true); err != nil {
			return err
		}
		if err := user1.CreateCommentVoteWithTx(tx, comment2.ID, true); err != nil {
			return err
		}
		if err := user1.CreateCommentVoteWithTx(tx, comment3.ID, true); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Errorf("Couldn't perform a transactional update of the comment votes: %v\n", err)
	}

	comment1.FindByID(comment1.ID)
	if comment1.LineageScore != 3 {
		t.Errorf("Comment's lineage score should be 3. Got %v instead.\n", comment1.LineageScore)
	}

	comment2.FindByID(comment2.ID)
	if comment2.LineageScore != 2 {
		t.Errorf("Comment's lineage score should be 2. Got %v instead.\n", comment2.LineageScore)
	}

	comment3.FindByID(comment3.ID)
	if comment3.LineageScore != 1 {
		t.Errorf("Comment's lineage score should be 1. Got %v instead.\n", comment3.LineageScore)
	}

	comment4, _ := user2.CreateComment(post, nil, "Test comment from user 2 on user 1's post")
	comment5, _ := user2.CreateComment(post, &comment4, "Test comment replying to comment 4")
	comment6, _ := user2.CreateComment(post, &comment5, "Test comment replying to comment 5")

	if err := WithTransaction(func(tx Transaction) error {
		if err := user1.CreateCommentVoteWithTx(tx, comment4.ID, false); err != nil {
			return err
		}
		if err := user1.CreateCommentVoteWithTx(tx, comment5.ID, false); err != nil {
			return err
		}
		if err := user1.CreateCommentVoteWithTx(tx, comment6.ID, false); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Errorf("Couldn't perform a transactional update of the comment votes: %v\n", err)
	}

	comment4.FindByID(comment4.ID)
	if comment4.LineageScore != -3 {
		t.Errorf("Comment's lineage score should be -3. Got %v instead.\n", comment4.LineageScore)
	}

	comment5.FindByID(comment5.ID)
	if comment5.LineageScore != -2 {
		t.Errorf("Comment's lineage score should be -2. Got %v instead.\n", comment5.LineageScore)
	}

	comment6.FindByID(comment6.ID)
	if comment6.LineageScore != -1 {
		t.Errorf("Comment's lineage score should be -1. Got %v instead.\n", comment6.LineageScore)
	}

	comment7, _ := user2.CreateComment(post, nil, "Test comment from user 2 on user 1's post")
	comment8, _ := user2.CreateComment(post, &comment7, "Test comment replying to comment 7")

	if err := WithTransaction(func(tx Transaction) error {
		if err := user1.CreateCommentVoteWithTx(tx, comment7.ID, true); err != nil {
			return err
		}
		if err := user1.CreateCommentVoteWithTx(tx, comment8.ID, false); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Errorf("Couldn't perform a transactional update of the comment votes: %v\n", err)
	}

	comment7.FindByID(comment7.ID)
	if comment7.LineageScore != 0 {
		t.Errorf("Comment's lineage score should be 0. Got %v instead.\n", comment7.LineageScore)
	}
}

func TestWeCanGetUpvotesAndDownvotes(t *testing.T) {
	setupTestingDB()

	bob, _ := UserFactory("example1@example.com", "bob", "password", 0)
	ralph, _ := UserFactory("example2@example.com", "ralph", "password", 0)
	joe, _ := UserFactory("example3@example.com", "joe", "password", 0)
	wanda, _ := UserFactory("example4@example.com", "wanda", "password", 0)

	post, _ := bob.CreatePost("Post Title", "post-title", "lorem ipsum", "image", 0, 0)

	comment1, _ := ralph.CreateComment(post, nil, "Test comment from ralph on bob's post")

	if err := WithTransaction(func(tx Transaction) error {
		if err := bob.CreateCommentVoteWithTx(tx, comment1.ID, true); err != nil {
			return err
		}
		if err := joe.CreateCommentVoteWithTx(tx, comment1.ID, true); err != nil {
			return err
		}
		if err := wanda.CreateCommentVoteWithTx(tx, comment1.ID, true); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Errorf("Couldn't perform a transactional update of the comment votes: %v\n", err)
	}
}
