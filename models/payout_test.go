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

	user1.UpdateAccountType("supporter")
	user2.UpdateAccountType("supporter")
	user3.UpdateAccountType("supporter")

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

	user1.SetCash(20)
	user2.SetCash(20)
	user3.SetCash(20)

	currentTime := time.Now()
	payouts, err := calculatePayouts(currentTime)
	if err != nil {
		t.Errorf("Payout calculation failed: %v\n", err)
	}
	if len(payouts) != 3 {
		t.Errorf("Returned %v payouts instead of the 3 expected.", len(payouts))
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

	user1.UpdateAccountType("supporter")
	user2.UpdateAccountType("supporter")
	user3.UpdateAccountType("supporter")

	user1.UpdateSubscriptionAmount(10)
	user2.UpdateSubscriptionAmount(10)
	user3.UpdateSubscriptionAmount(10)

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

func TestLedgerEntries(t *testing.T) {
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

	user1.SetCash(20)
	user2.SetCash(20)
	user3.SetCash(20)

	ledgerEntries, err := calculatePayouts(time.Now())
	if err != nil {
		t.Errorf("Ledger calculation failed")
	}

	if len(ledgerEntries) != 5 {
		t.Errorf("5 ledger entries expected. Got %v", len(ledgerEntries))
	}

	if ledgerEntries[0].authorId == 2 && ledgerEntries[0].amount != 4.5 {
		t.Errorf("Contribution of User 1 to User 2 of amount 4.5 missing")
	}

	if ledgerEntries[1].authorId == 3 && ledgerEntries[1].amount != 4.5 {
		t.Errorf("Contribution of User 1 to User 3 of amount 4.5 missing")
	}

	if ledgerEntries[2].authorId == 3 && ledgerEntries[2].amount != 9 {
		t.Errorf("Contribution of User 2 to User 3 of amount 9 missing")
	}

	if ledgerEntries[3].authorId == 1 && ledgerEntries[3].amount != 4.5 {
		t.Errorf("Contribution of User 3 to User 1 of amount 4.5 missing")
	}
	if ledgerEntries[4].authorId == 2 && ledgerEntries[4].amount != 4.5 {
		t.Errorf("Contribution of User 3 to User 2 of amount 4.5 missing")
	}
}
