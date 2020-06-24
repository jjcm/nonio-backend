package models

import (
	"testing"
)

func TestWeCanCreateComments(t *testing.T) {
	setupTestingDB()

	author, _ := UserFactory("example@example.com", "", "password")

	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	// create comment
	comment, err := author.CommentOnPost(post, nil, "This is a dumb post")

	if err != nil {
		t.Errorf("We should have been able to create a comment. Error recieved: %s", err)
	}
	if comment.ID == 0 {
		t.Error("The created comment should have been instantiated correctly.")
	}
}

func TestWeCanStructureCommentsCorrectly(t *testing.T) {
	setupTestingDB()

	// setup users
	person, _ := UserFactory("person@example.com", "friendlyPerson", "password")
	troll, _ := UserFactory("troll@example.com", "uglyTroll", "password")
	moderator, _ := UserFactory("moderator@example.com", "uglyTroll", "password")

	// create a post and a typical comment thread
	post, _ := person.CreatePost("Post Title", "post-title", "lorem ipsum", "image")
	nastyComment, _ := troll.CommentOnPost(post, nil, "WOW, this post is dumb!")
	reply, _ := person.CommentOnPost(post, &nastyComment, "Hay, I worked really hard on this")
	nastyAgain, _ := troll.CommentOnPost(post, &reply, "Don't you mean \"Hey\"? You must be as dumb as the post you created.")
	moderatorComment, _ := moderator.CommentOnPost(post, nil, "Hey troll, play nice or go somewhere else. We may shut down commenting if you keep this up.")
	trollIsAskingForIt, _ := troll.CommentOnPost(post, &moderatorComment, "You're dumb too!")

	// at this point, we should have 5 comments on the post and the tree should look like this:
	// Post Comments:
	//   - nastyComment
	//     - reply
	//       - nastyAgain
	//   - moderatorComment
	//     - trollIsAskingForIt
	expectedParentChildRelations := map[int]Comment{
		moderatorComment.ID: trollIsAskingForIt,
		nastyComment.ID:     reply,
		reply.ID:            nastyAgain,
	}
	for id, c := range expectedParentChildRelations {
		if c.ParentID != id {
			t.Errorf("Expected comment relation failed. ParentID: %v, Child's ParentID: %v", id, c.ParentID)
		}
	}
}
