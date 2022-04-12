package service

import (
	"go-blog/services/conf"
	"go-blog/services/database"
	"go-blog/services/store"
)

const InternalServerError = "Something went wrong!"

func Start(cfg conf.Config) {
	// 建立jwt的签名者和验证者
	jwtSetup(cfg)

	// 读取数据库配置,并建立数据库链接
	store.SetDBConnection(database.NewDBOptions(cfg))

	// 获取路由
	router := setRouter()
	// 监听路由
	router.Run(":8080")
}
