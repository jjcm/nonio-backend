package finance

import (
	"testing"
)

func TestWeCanAllocatePayouts(t *testing.T) {
	setupTestingDB()

	user1, _ := UserFactory("example1@example.com", "", "password", 0)
	user2, _ := UserFactory("example2@example.com", "", "password", 0)
	user3, _ := UserFactory("example3@example.com", "", "password", 0)

}
