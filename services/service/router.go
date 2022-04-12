package service

import (
	"go-blog/services/store"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setRouter() *gin.Engine {
	// 创建默认的gin路由器, 并且已经附加了Logger 个 Recovery中间件
	router := gin.Default()

	// Enables automatic redirection if the current route can’t be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = true

	// 创建API路由组
	api := router.Group("/api")
	api.Use(customErrors)
	{
		api.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "world"})
		})
		api.POST("/signup", gin.Bind(store.User{}), signUp)
		api.POST("/signin", gin.Bind(store.User{}), signIn)
	}

	// 博客路由组
	authorized := api.Group("/")
	authorized.Use(authorization) // 使用中间件认证用户
	{
		authorized.GET("/posts", indexPosts)
		authorized.POST("/posts", createPost)
		authorized.PUT("/posts", updatePost)
		authorized.DELETE("/posts/:id", deletePost)
	}

	// 404
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	return router

}
