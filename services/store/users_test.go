package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	// 清空数据库
	testSetUp()
	// 创建用户
	user, err := addTestUser()
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.NotEmpty(t, user.Salt)
	assert.NotEmpty(t, user.HashedPassword)
}

func TestAddUserWithExistUserName(t *testing.T) {
	// 清空数据库
	testSetUp()
	// 创建一个用户
	user, err := addTestUser()
	assert.NoError(t, err) // 没有报错
	assert.Equal(t, 1, user.ID)
	// 再穿件一个同名用户
	user, err = addTestUser()
	assert.Error(t, err) // 报错
	assert.Equal(t, "Username already exists.", err.Error())
}
