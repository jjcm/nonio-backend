package models

import "testing"

func TestWeCanFindPostTagVoteByID(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
	}

	p := &PostTagVote{}
	p.FindByID(1)

	if p.ID == 0 {
		t.Errorf("We should have been able to find this PostTagVote by it's ID")
	}
}

func TestWeCanCreatePostTagVote(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
	}
}
