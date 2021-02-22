package finance

import (
	"fmt"
	"soci-backend/models"
	"testing"
)

func TestWeCanCallModelPackage(t *testing.T) {
	setupTestingDB()
	s := models.DemoReturnString()
	fmt.Println(s)
}

func TestWeCanCallFuncWithDBCall(t *testing.T) {
	models.SetupTestingDB()
	err := DemoInsertUser()
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
}

func TestWeCanCallModelFuncWithDBCall(t *testing.T) {
	err := models.DemoInsertUser()
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
}

func TestWeCanAllocatePayouts(t *testing.T) {
	setupTestingDB()

	s := models.DemoReturnString()
	fmt.Println(s)
	/*
		user, err := models.UserFactory("example@example.com", "", "password", 0)
		if err != nil {
			t.Errorf("%v", err)
		}
		fmt.Printf("%v", user.ID)
		/*
			user1, _ := models.UserFactory("example1@example.com", "", "password", 0)
			user2, _ := models.UserFactory("example2@example.com", "", "password", 0)
			user3, _ := models.UserFactory("example3@example.com", "", "password", 0)

			fmt.Printf("%v %v %v", user1.ID, user2.ID, user3.ID)
			/*
				post1, _ := user1.CreatePost("Post Title", "test-post-1", "lorem ipsum", "image", 0, 0)
				post2, _ := user2.CreatePost("Post Title", "test-post-2", "lorem ipsum", "image", 0, 0)
				post3, _ := user3.CreatePost("Post Title", "test-post-3", "lorem ipsum", "image", 0, 0)

				// User 1 votes on all 3 posts (including their own). Expected payout is $5 each
				user1.CreatePostTagVote(post1.ID, 1)
				user1.CreatePostTagVote(post2.ID, 1)
				user1.CreatePostTagVote(post3.ID, 1)

				// User 2 votes on only user 3's post. Expected payout is $20
				user2.CreatePostTagVote(post3.ID, 1)

				// User 3 votes on user 1 and user 2's posts. Expected payout is $2 each
				user3.CreatePostTagVote(post1.ID, 1)
				user3.CreatePostTagVote(post2.ID, 1)

					payouts, err := CalculatePayouts()
					if err != nil {
						t.Errorf("Payout calculation failed: %v", err)
					}

					for _, payout := range payouts {
						fmt.Printf("user %v, payout: %v", payout.UserID, payout.Payout)

					}
	*/
}
