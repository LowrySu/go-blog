package store

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Post struct {
	ID         int
	Title      string `binding:"required,min=3,max=50"`
	Content    string `binding:"required,min=5,max=5000"`
	CreatedAt  time.Time
	ModifiedAt time.Time
	UserID     int `json:"-"`
}

func AddPost(user *User, post *Post) error {
	post.UserID = user.ID
	_, err := db.Model(post).Returning("*").Insert()
	if err != nil {
		log.Panic().Err(err).Msg("Error inserting new post")
	}
	return err
}
