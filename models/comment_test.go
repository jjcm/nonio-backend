package models

import "testing"

func TestWeCanCreateComments(t *testing.T) {
	// setup
	setupTestingDB()
	defer teardownTestingDB()
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	post, _ := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	// create comment
	comment, err := author.CommentOnPost(post, nil, "text", "This is a dumb post", "")

	if err != nil {
		t.Errorf("We should have been able to create a comment. Error recieved: %s", err)
	}
	if comment.ID == 0 {
		t.Error("The created comment should have been instantiated correctly.")
	}
}

func TestWeCanStructureCommentsCorrectly(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// setup users
	CreateUser("person@example.com", "friendlyPerson", "password")
	CreateUser("troll@example.com", "uglyTroll", "password")
	CreateUser("moderator@example.com", "uglyTroll", "password")
	person := User{}
	person.FindByEmail("person@example.com")
	troll := User{}
	troll.FindByEmail("troll@example.com")
	moderator := User{}
	moderator.FindByEmail("moderator@example.com")

	// create a post and a typical comment thread
	post, _ := person.CreatePost("Post Title", "post-title", "lorem ipsum", "image")
	nastyComment, _ := troll.CommentOnPost(post, nil, "text", "WOW, this post is dumb!", "")
	reply, _ := person.CommentOnPost(post, &nastyComment, "text", "Hay, I worked really hard on this", "")
	nastyAgain, _ := troll.CommentOnPost(post, &reply, "text", "Don't you mean \"Hey\"? You must be as dumb as the post you created.", "")
	moderatorComment, _ := moderator.CommentOnPost(post, nil, "text", "Hey troll, play nice or go somewhere else. We may shut down commenting if you keep this up.", "")
	trollIsAskingForIt, _ := troll.CommentOnPost(post, &moderatorComment, "text", "You're dumb too!", "")

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
