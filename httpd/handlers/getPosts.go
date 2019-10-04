package handlers

import (
	"net/http"
	"time"

	"github.com/jjcm/soci-backend/models"
)

// GetPosts - get all the posts in the system
func GetPosts(w http.ResponseWriter, r *http.Request) {
	author := models.User{
		Name: "bob",
	}
	tagWTF := models.Tag{
		Name:  "wtf",
		Score: 10,
	}
	tagPhotography := models.Tag{
		Name:  "photography",
		Score: 10,
	}
	p1 := models.Post{
		Title:     "asdf",
		Author:    author,
		CreatedAt: time.Now(),
		URL:       "url-example-1",
		Tags: []models.Tag{
			tagWTF,
		},
	}
	p2 := models.Post{
		Title:     "cats are known to be mind controlling beings from another plant, or perhaps area 51",
		Author:    author,
		CreatedAt: time.Now(),
		URL:       "url-example-2",
		Tags: []models.Tag{
			tagWTF,
			tagPhotography,
		},
	}
	p3 := models.Post{
		Title:     "asdf2",
		Author:    author,
		CreatedAt: time.Now(),
		URL:       "url-example-3",
		Tags: []models.Tag{
			tagWTF,
		},
	}
	posts := []models.Post{}
	posts = append(posts, p1, p2, p3, p1, p2, p3)
	SendResponse(w, posts, 200)
}
