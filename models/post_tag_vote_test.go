package models

import (
	"fmt"
	"testing"
)

func TestWeCanFindPostTagVoteByUK(t *testing.T) {
	setupTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
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
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
		return
	}

	votes, err := item.GetVotesByPostUser(item.PostID, item.VoterID)
	if err != nil {
		t.Errorf("Get votes: %v", err)
	}
	if len(votes) != 1 {
		t.Errorf("We should have been able to find this PostTagVote voter ID")
	}
}

func TestWeCanGetUntalliedPostTagVotesForAUser(t *testing.T) {
	setupTestingDB()

	// create the PostTagVote first
	item := &PostTagVote{
		PostID:  1,
		TagID:   1,
		VoterID: 1,
	}
	if err := item.CreatePostTagVote(); err != nil {
		t.Errorf("PostTagVote creation should have worked: %v", err)
		return
	}

	item = &PostTagVote{
		PostID:  2,
		TagID:   2,
		VoterID: 1,
	}
	item.CreatePostTagVote()

	item = &PostTagVote{
		PostID:  3,
		TagID:   3,
		VoterID: 2,
	}
	item.CreatePostTagVote()

	votes, err := item.GetUntalliedVotesByUser(1)
	if err != nil {
		t.Errorf("Get votes: %v", err)
	}
	if len(votes) != 2 {
		t.Errorf(fmt.Sprintf("%v", len(votes)))
	}
}
