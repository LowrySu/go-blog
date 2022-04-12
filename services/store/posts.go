package store

import (
	"time"

	"github.com/go-pg/pg/v10/orm"
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

func FetchUserPosts(user *User) error {
	// 输入user变量, 在数据库搜索user_id = user.id 的post, 并赋值到user变量的posts属性里
	err := db.Model(user).WherePK().
		Relation(
			"Posts", func(q *orm.Query) (*orm.Query, error) {
				return q.Order("id ASC"), nil
			}).Select()
	if err != nil {
		log.Error().Err(err).Msg("Error fetch user's posts")
	}
	return err
}

func FetchPost(id int) (*Post, error) {
	// 通过post的id 获取post
	post := new(Post)
	post.ID = id
	err := db.Model(post).WherePK().Select()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching post")
		return nil, err
	}
	return post, nil
}

func UpdatePost(post *Post) error {
	// 通过id找到需要修改的post并修改其属性
	_, err := db.Model(post).WherePK().UpdateNotZero()
	if err != nil {
		log.Error().Err(err).Msg("Error updating post")
	}
	return err
}

func DeletePost(post *Post) error {
	// 通过id找到需要修改的post并删除该post
	_, err := db.Model(post).WherePK().Delete()
	if err != nil {
		log.Error().Err(err).Msg("Error deleting post")
	}
	return err
}
