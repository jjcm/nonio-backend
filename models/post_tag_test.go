package models

import "testing"

func TestWeCanFindAPostTagByItsID(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create the PostTag first
	item := &PostTag{
		PostID: 1,
		TagID:  1,
	}
	if err := item.CreatePostTag(); err != nil {
		t.Errorf("PostTag creation should have worked: %v", err)
	}

	p := &PostTag{}
	p.FindByID(1)

	if p.ID == 0 {
		t.Errorf("We should have been able to find this PostTag by it's ID")
	}
}

func TestWeCanFindAPostTagByPostTagIds(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create the PostTag first
	item := &PostTag{
		PostID: 1,
		TagID:  1,
	}
	if err := item.CreatePostTag(); err != nil {
		t.Errorf("PostTag creation should have worked: %v", err)
	}

	p := &PostTag{}
	p.FindByPostTagIds(1, 1)
	if p.ID == 0 {
		t.Errorf("We should have been able to find this PostTag by PostID and TagID")
	}
}

func TestWeCanCreatePostTag(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	p := &PostTag{
		PostID: 1,
		TagID:  1,
	}
	if err := p.CreatePostTag(); err != nil {
		t.Errorf("PostTag creation should have worked: %v", err)
	}
}
