package models

import (
	"fmt"
	"testing"
)

func TestWeCanAllocatePayouts(t *testing.T) {
	setupTestingDB()

	user1, _ := UserFactory("example1@example.com", "ralph", "password", 10+ServerFee)
	user2, _ := UserFactory("example2@example.com", "joey", "password", 20+ServerFee)
	user3, _ := UserFactory("example3@example.com", "bobby", "password", 4+ServerFee)

	post1, _ := user1.CreatePost("Post Title", "test-post-1", "lorem ipsum", "image", 0, 0)
	post2, _ := user2.CreatePost("Post Title", "test-post-2", "lorem ipsum", "image", 0, 0)
	post3, _ := user3.CreatePost("Post Title", "test-post-3", "lorem ipsum", "image", 0, 0)

	// User 1 votes on all 3 posts (including their own). Expected payout is $5 each
	user1.CreatePostTagVote(post1.ID, 1)
	user1.CreatePostTagVote(post2.ID, 1)
	user1.CreatePostTagVote(post3.ID, 1)

	// User 2 votes on only user 3's post. Expected payout is $20
	user2.CreatePostTagVote(post3.ID, 1)

	// User 3 votes on user 1 and user 2's posts. They vote for two tags on user 2's posts. Expected payout is $2 each
	user3.CreatePostTagVote(post1.ID, 1)
	user3.CreatePostTagVote(post2.ID, 1)
	user3.CreatePostTagVote(post2.ID, 2)

	payouts, err := calculatePayouts()
	AllocatePayouts()
	if err != nil {
		t.Errorf("Payout calculation failed: %v\n", err)
	}
	if len(payouts) != 5 {
		t.Errorf("Returned %v payouts instead of the 5 expected.", len(payouts))
	}

	for _, payout := range payouts {
		fmt.Printf("user %v, payout: %v\n", payout.UserID, payout.Payout)
	}

	t.Errorf("Returned %v payouts instead of the 5 expected.", len(payouts))
}
