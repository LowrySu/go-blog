package main

import (
	"go-blog/services/cli"
	"go-blog/services/conf"
	"go-blog/services/service"
)

func main() {
	// 读取环境变量
	cli.Parse()

	// 启动web服务
	service.Start(conf.NewConfig())
}
