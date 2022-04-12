package store

import (
	"github.com/gin-gonic/gin"
)

func testSetUp() {
	gin.SetMode(gin.TestMode)
	// 清空测试数据库
	ResetTestDatabase()
}

func addTestUser() (*User, error) {
	user := &User{
		Username: "batman",
		Password: "secret123",
	}

	err := AddUser(user)
	return user, err
}

func addTestPost(user *User) (*Post, error) {

	post := &Post{
		Title:   "test post title",
		Content: "test post content",
	}
	err := AddPost(user, post)
	return post, err
}
