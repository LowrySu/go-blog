package service

import (
	"go-blog/services/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signUp(ctx *gin.Context) {
	// 获取属性并创建用户
	user := ctx.MustGet(gin.BindKey).(*store.User)

	// 添加用户
	if err := store.AddUser(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Sign up successfully",
		"jwt": generateJWT(user),
	})
}

func signIn(ctx *gin.Context) {

	// 新建一个用户
	user := ctx.MustGet(gin.BindKey).(*store.User)

	// 验证请求信息和数据库信息
	user, err := store.Authenticate(user.Username, user.Password)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Sign in failed."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Sign in successfully",
		"jwt": generateJWT(user),
	})
}
