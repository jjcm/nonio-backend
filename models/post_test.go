package models

import (
	"strconv"
	"testing"
)

func TestWeCanCreateAPost(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	p, err := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")
	if err != nil {
		t.Errorf("Post creation should have worked. Error recieved: %v", err)
	}
	if p.Title != "Post Title" {
		t.Errorf("The post that is returned should be the newly created post")
	}
	if p.Author.ID != author.ID {
		t.Errorf("The author associated with this post should be hydrated automatically")
	}
}

func TestWeCanIncrementScoreForPost(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	p := Post{}
	if err := p.FindByURL("post-title"); err != nil {
		t.Errorf("Find post by url: %v", err)
		return
	}
	if err := p.IncrementScore(p.ID); err != nil {
		t.Errorf("Increment score: %v", err)
		return
	}
}

func TestWeCanFindAPostByItsURL(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	p := Post{}
	p.FindByURL("post-title")

	if p.ID == 0 {
		t.Errorf("We should have been able to find this post by it's ID")
	}
}

func TestWeCanFindAPostByItsID(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	p := Post{}
	p.FindByID(1)

	if p.ID == 0 {
		t.Errorf("We should have been able to find this post by it's ID")
	}
}

func TestWhenWeMarshalAPostToJSONItHasTheShapeThatWeExpect(t *testing.T) {

	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "userName", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	p, err := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")
	if err != nil {
		t.Errorf("We should be able to create a post. Error: %v", err)
	}

	expectedJSON := `{"title":"Post Title","user":"userName","time":` + strconv.Itoa(int(p.GetCreatedAtTimestamp())) + `,"url":"post-title","content":"lorem ipsum","type":"image","score":0,"tags":[]}`

	if p.ToJSON() != expectedJSON {
		t.Errorf("JSON output wasn't what we expected it to be.\nExpected: %v\nActual:   %v", expectedJSON, p.ToJSON())
	}
}

func TestWeTagAPost(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	p := Post{}
	p.FindByURL("post-title")

	CreateTag("Tag!", author)
	tag := Tag{}
	tag.FindByTagName("Tag!")

	err := p.AddTag(tag)
	if err != nil {
		t.Errorf("We should be able to tag this freshly created post. Error: %v", err.Error())
	}
	if len(p.Tags) != 1 {
		t.Errorf("We expected that the post should now have 1 tag associated with it. Post tag count is %v", len(p.Tags))
	}
}

func TestWeCantTagAPostWithTheSameTagMoreThanOneTime(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")

	p := Post{}
	p.FindByURL("post-title")

	CreateTag("Tag!", author)
	tag := Tag{}
	tag.FindByTagName("Tag!")

	err := p.AddTag(tag)
	if err != nil {
		t.Errorf("We should be able to tag this freshly created post. Error: %v", err.Error())
	}
	if len(p.Tags) != 1 {
		t.Errorf("We expected that the post should now have 1 tag associated with it. Post tag count is %v", len(p.Tags))
	}

	err = p.AddTag(tag)
	if err == nil {
		t.Errorf("The tag has already been tagged with this tag, an error should have occured")
	}
	if len(p.Tags) != 1 {
		t.Errorf("We expected that the post should now have 1 tag associated with it. Post tag count is %v", len(p.Tags))
	}
}

func TestIfWeCreateAPostWithTheSameURLTheSystemWillGenerateAUniqueOne(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	p1, err := author.CreatePost("Post Title", "post-title", "lorem ipsum", "image")
	if err != nil {
		t.Errorf("Post creation should have worked. Error recieved: %v", err)
	}
	if p1.URL != "post-title" {
		t.Errorf("The post that is returned should be the newly created post")
	}

	// now let's create a second post with the same title
	p2, err := author.CreatePost("Post Title", "post-title", "Dolor sit amit", "image")
	if err != nil {
		t.Errorf("Post creation should have worked. Error recieved: %v", err)
	}
	if p2.URL != "post-title-2" {
		t.Errorf("The URL for the second post should be different from the original created post. Post2 url: %v", p2.URL)
	}

	// now let's create a third post with the same title
	p3, err := author.CreatePost("Post Title", "post-title", "Dolor sit amit", "image")
	if err != nil {
		t.Errorf("Post creation should have worked. Error recieved: %v", err)
	}
	if p3.URL != "post-title-3" {
		t.Errorf("The URL for the third post should be different from the original created post. Post3 url: %v", p2.URL)
	}
}

func TestIfAPostIsCreatedWithAnEmptyTypeItGetsSetToTheDefaultTypeImage(t *testing.T) {
	setupTestingDB()
	defer teardownTestingDB()

	// create an author for post
	CreateUser("example@example.com", "", "password")
	author := User{}
	author.FindByEmail("example@example.com")
	p, err := author.CreatePost("Title", "post-title", "lorem ipsum", "") // note the empty 3rd param
	if err != nil {
		t.Errorf("Post creation should have worked. Error recieved: %v", err)
	}
	if p.Type != "image" {
		t.Errorf("The post type should be the default `Image`. Type is: %v", p.Type)
	}
}
