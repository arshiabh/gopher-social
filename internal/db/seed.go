package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/arshiabh/gopher-social/internal/store"
)

func Seed(store *store.Storage, db *sql.DB) {
	ctx := context.Background()
	users := generateUser(100)
	posts := generatePost(200, users)
	comments := generateComments(500, users, posts)

	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			tx.Rollback()
			log.Fatal(err)
			return
		}
	}
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Fatal(err)
			return
		}
	}
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Fatal(err)
			return
		}
	}

}

// ** dont forget to add add init
func generateUser(num int) []*store.User {
	users := make([]*store.User, num)
	for i := 0; i < num; i++ {
		users[i] = &store.User{
			ID:       int64(i),
			Username: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("email%d@gmail.com", i),
			RoleID:   1,
		}
		users[i].Password.Set("1234")
	}
	return users
}

func generatePost(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < 200; i++ {
		//** change the type right after to id and use it
		user := users[rand.Intn(100)].ID
		posts[i] = &store.Post{
			ID:      int64(i),
			UserID:  user,
			Title:   fmt.Sprintf("title %d", i),
			Content: fmt.Sprintf("content for %d post", i),
			Tags:    []string{"#mbappe", "#vini"},
		}
	}
	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comments {
	comments := make([]*store.Comments, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(100)].ID
		post := posts[rand.Intn(100)].ID
		comments[i] = &store.Comments{
			ID:      int64(i),
			PostID:  post,
			UserID:  user,
			Content: fmt.Sprintf("comment number %d", i),
		}
	}
	return comments
}
