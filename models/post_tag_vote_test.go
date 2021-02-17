package models

import "testing"

func TestWeCanFindPostTagVoteByUK(t *testing.T) {
	setupTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
		Tallied: false,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
		return
	}

	p := &PostTagVote{}
	p.FindByUK(item.PostID, item.TagID, item.VoterID)

	if p.ID == 0 {
		t.Errorf("We should have been able to find this PostTagVote by it's ID")
	}
}

func TestWeCanCreatePostTagVote(t *testing.T) {
	setupTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
		Tallied: false,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
	}
}

func TestWeCanGetPostTagVotesByPostUser(t *testing.T) {
	setupTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
		Tallied: false,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
		return
	}

	votes, err := item.GetVotesByPostUser(item.PostID, item.VoterID)
	if err != nil {
		t.Errorf("Get votes: %v", err)
	}
	if len(votes) == 0 {
		t.Errorf("We should have been able to find this PostTagVote by post id and voter id")
	}
}

func TestWeCanGetUntalliedPostTagVotes(t *testing.T) {
	setupTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
		Tallied: false,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
		return
	}
	// create the PostTagVote first
	item = &PostTagVote{
		PostID:  2,
		TagID:   2,
		VoterID: 2,
		Tallied: false,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
		return
	}
	// create the PostTagVote first
	item = &PostTagVote{
		PostID:  2,
		TagID:   2,
		VoterID: 2,
		Tallied: false,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
		return
	}

	votes, err := item.GetVotesByPostUser(item.PostID, item.VoterID)
	if err != nil {
		t.Errorf("Get votes: %v", err)
	}
	if len(votes) == 0 {
		t.Errorf("We should have been able to find this PostTagVote by post id and voter id")
	}
}
