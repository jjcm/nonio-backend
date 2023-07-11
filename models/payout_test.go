package models

import (
	"testing"
	"time"
)

// TODO need to update the users to set the subscription level correctly

func TestWeCanAllocatePayouts(t *testing.T) {
	setupTestingDB()
	//user1, _ := UserFactory("example1@example.com", "ralph", "password", 10+ServerFee)
	//user2, _ := UserFactory("example2@example.com", "joey", "password", 20+ServerFee)
	//user3, _ := UserFactory("example3@example.com", "bobby", "password", 4+ServerFee)
	user1, _ := UserFactory("example1@example.com", "ralph", "password")
	user2, _ := UserFactory("example2@example.com", "joey", "password")
	user3, _ := UserFactory("example3@example.com", "bobby", "password")

	user1.UpdateAccountType("supporter")
	user2.UpdateAccountType("supporter")
	user3.UpdateAccountType("supporter")

	user1.UpdateSubscriptionAmount(10)
	user2.UpdateSubscriptionAmount(10)
	user3.UpdateSubscriptionAmount(10)

	post1, _ := user1.CreatePost("Post Title", "test-post-1", "", "lorem ipsum", "image", 0, 0)
	post2, _ := user2.CreatePost("Post Title", "test-post-2", "", "lorem ipsum", "image", 0, 0)
	post3, _ := user3.CreatePost("Post Title", "test-post-3", "", "lorem ipsum", "image", 0, 0)

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

	user1.SetCash(20)
	user2.SetCash(20)
	user3.SetCash(20)

	AllocatePayouts()

	cash1, err := user1.GetCash()
	if err != nil {
		t.Errorf("Couldn't get financial data: %v", err)
	}
	if cash1 != 24.5 {
		t.Errorf("User 1's cash was incorrect. Got %v expected 4.5.", cash1)
	}

	cash2, _ := user2.GetCash()
	if cash2 != 29 {
		t.Errorf("User 2's cash was incorrect. Got %v expected 9.", cash2)
	}

	cash3, _ := user3.GetCash()
	if cash3 != 33.5 {
		t.Errorf("User 3's cash was incorrect. Got %v expected 13.5.", cash3)
	}

	// try allocating again to ensure double counts arent happening
	AllocatePayouts()
	cash1, _ = user1.GetCash()
	if err != nil {
		t.Errorf("Couldn't get financial data: %v", err)
	}
	if cash1 != 4.5 {
		t.Errorf("User 1's cash was incorrect. Got %v expected 4.5.", cash1)
	}
}

func TestWeCanProcessPayouts(t *testing.T) {
	setupTestingDB()

	user1, _ := UserFactory("example1@example.com", "ralph", "password")
	user2, _ := UserFactory("example2@example.com", "joey", "password")
	user3, _ := UserFactory("example3@example.com", "bobby", "password")

	post1, _ := user1.CreatePost("Post Title", "test-post-1", "", "lorem ipsum", "image", 0, 0)
	post2, _ := user2.CreatePost("Post Title", "test-post-2", "", "lorem ipsum", "image", 0, 0)
	post3, _ := user3.CreatePost("Post Title", "test-post-3", "", "lorem ipsum", "image", 0, 0)

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

	time.Sleep(1 * time.Second)

	user1.CreateFuturePayout(10, time.Now())
	user2.CreateFuturePayout(20, time.Now())
	user3.CreateFuturePayout(4, time.Now())

	ProcessPayouts()

	if cash, _ := user1.GetCash(); cash != 2 {
		t.Errorf("Cash of user 1 expected to be $2, instead got: %v", cash)
	}

	if cash, _ := user2.GetCash(); cash != 7 {
		t.Errorf("Cash of user 2 expected to be $7, instead got: %v", cash)
	}

	if cash, _ := user3.GetCash(); cash != 25 {
		t.Errorf("Cash of user 3 expected to be $25, instead got: %v", cash)
	}
}
