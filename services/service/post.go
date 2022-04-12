package service

import (
	"go-blog/services/store"
	"net/http"
	"strconv"
	"time"

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

func indexPosts(ctx *gin.Context) {
	// 查看当前用户的所有博客

	// 获取当前用户
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	// 去数据库查询该用户的所有博客并赋值到用户的Posts属性里
	if err := store.FetchUserPosts(user); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Posts fetched successfully.",
		"data": user.Posts,
	})

}

func updatePost(ctx *gin.Context) {
	// 修改博客

	// 查看请求信息是否匹配post的属性(是否有Title和Content属性)
	jsonPost := new(store.Post)
	if err := ctx.Bind(jsonPost); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取当前用户
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取数据库中需要修改的博客
	dbPost, err := store.FetchPost(jsonPost.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 验证当前用户操作的是否是自己的博客
	if dbPost.UserID != user.ID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Not authorized.",
		})
		return
	}

	// 修改时间
	jsonPost.ModifiedAt = time.Now()

	// 修改数据库中的数据
	if err := store.UpdatePost(jsonPost); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Post updated successfully.",
		"data": jsonPost,
	})
}

func deletePost(ctx *gin.Context) {
	// 删除博客

	// 提取url路径中的id属性
	paramID := ctx.Param("id")

	// 把id转化为数字
	id, err := strconv.Atoi(paramID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Not valid ID",
		})
	}

	// 获取当前用户
	user, err := currentUser(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	// 使用id获取数据库中的博客
	post, err := store.FetchPost(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	// 验证当前用户是否有权限删除该博客
	if user.ID != post.UserID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
	}

	// 删除数据库中的该博客
	if err := store.DeletePost(post); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "Post deleted successfully.",
	})
}
