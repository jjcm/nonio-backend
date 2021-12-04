package models

import (
	"testing"
	"time"
)

// TODO need to update the users to set the subscription level correctly
func TestWeCanCalculatePayouts(t *testing.T) {
	setupTestingDB()

	//user1, _ := UserFactory("example1@example.com", "ralph", "password", 10+ServerFee)
	//user2, _ := UserFactory("example2@example.com", "joey", "password", 20+ServerFee)
	//user3, _ := UserFactory("example3@example.com", "bobby", "password", 4+ServerFee)
	user1, _ := UserFactory("example1@example.com", "ralph", "password")
	user2, _ := UserFactory("example2@example.com", "joey", "password")
	user3, _ := UserFactory("example3@example.com", "bobby", "password")

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

	currentTime := time.Now()
	payouts, err := calculatePayouts(currentTime)
	if err != nil {
		t.Errorf("Payout calculation failed: %v\n", err)
	}
	if len(payouts) != 5 {
		t.Errorf("Returned %v payouts instead of the 5 expected.", len(payouts))
	}
}

func TestWeCanAllocatePayouts(t *testing.T) {
	setupTestingDB()
	//user1, _ := UserFactory("example1@example.com", "ralph", "password", 10+ServerFee)
	//user2, _ := UserFactory("example2@example.com", "joey", "password", 20+ServerFee)
	//user3, _ := UserFactory("example3@example.com", "bobby", "password", 4+ServerFee)
	user1, _ := UserFactory("example1@example.com", "ralph", "password")
	user2, _ := UserFactory("example2@example.com", "joey", "password")
	user3, _ := UserFactory("example3@example.com", "bobby", "password")

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

	AllocatePayouts()

	financialData1, err := user1.GetFinancialData()
	if err != nil {
		t.Errorf("Couldn't get financial data: %v", err)
	}
	if financialData1.Cash != 2 {
		t.Errorf("User 1's cash was incorrect. Got %v expected 2.", financialData1.Cash)
	}

	financialData2, _ := user2.GetFinancialData()
	if financialData2.Cash != 7 {
		t.Errorf("User 2's cash was incorrect. Got %v expected 7.", financialData2.Cash)
	}

	financialData3, _ := user3.GetFinancialData()
	if financialData3.Cash != 25 {
		t.Errorf("User 3's cash was incorrect. Got %v expected 25.", financialData3.Cash)
	}

	// try allocating again to ensure double counts arent happening
	AllocatePayouts()
	financialData1, _ = user1.GetFinancialData()
	if err != nil {
		t.Errorf("Couldn't get financial data: %v", err)
	}
	if financialData1.Cash != 2 {
		t.Errorf("User 1's cash was incorrect. Got %v expected 2.", financialData1.Cash)
	}
}
