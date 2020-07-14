package models

import (
	"testing"
)

func TestWeCanCreateASubscription(t *testing.T) {
	setupTestingDB()

	// create a user and a tag
	user, _ := UserFactory("example@example.com", "", "password")
	tag, _ := TagFactory("funny", user)

	// use our subscription factory to make a new subscription
	subscription, err := SubscriptionFactory(tag, user)
	if err != nil {
		t.Errorf("Subscription creation should have worked. Error recieved: %v", err)
	}
	if subscription.TagName != "Funny" {
		t.Errorf("The tag name should have been 'funny'.")
	}
	if subscription.User.ID != user.ID {
		t.Errorf("The user associated with this subscription should be hydrated automatically")
	}
}
