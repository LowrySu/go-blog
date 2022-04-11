package service

import (
	"go-blog/services/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createPost(ctx *gin.Context) {
	// 创建博客

	// 获取请求信息中的博客信息
	post := new(store.Post)
	if err := ctx.Bind(post); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取当前用户
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 添加博客
	if err := store.AddPost(user, post); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Post created successfully.",
		"data": post,
	})
}
