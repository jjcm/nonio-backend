package models

import (
	"testing"
)

func TestWeCanCreateASubscription(t *testing.T) {
	setupTestingDB()

	// create a user and a tag
	user, _ := UserFactory("example@example.com", "", "password", 0)
	tag, _ := TagFactory("funny", user)

	// use our subscription factory to make a new subscription
	subscription, err := user.CreateSubscription(tag)
	if err != nil {
		t.Errorf("Subscription creation should have worked. Error recieved: %v", err)
	}
	if subscription.TagID != tag.ID {
		t.Errorf("The tag should have been hydrated.")
	}
	if subscription.UserID != user.ID {
		t.Errorf("The user associated with this subscription should be hydrated automatically")
	}
}

func TestWeCanDeleteASubscription(t *testing.T) {
	setupTestingDB()

	// create a user and a tag
	user, _ := UserFactory("example@example.com", "", "password", 0)
	tag, _ := TagFactory("funny", user)

	// use our subscription factory to make a new subscription
	if _, err := user.CreateSubscription(tag); err != nil {
		t.Errorf("Subscription creation should have worked. Error recieved: %v", err)
	}
	if err := user.DeleteSubscription(tag); err != nil {
		t.Errorf("Subscription deletion should have worked. Error recieved: %v", err)
	}

	subscription := Subscription{}
	if err := subscription.FindSubscription(tag.ID, user.ID); err == nil {
		t.Errorf("The lookup should have thrown a no rows in result error. Error received: %v", err)
	}
}
