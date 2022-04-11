package cli

import (
	"flag"
	"fmt"
	"go-blog/services/logging"
	"os"
)

func usage() {
	fmt.Print(`This program runs backend server.
  
  Usage:
  
  pharos [arguments]
  
  Supported arguments:
  
  `)
	flag.PrintDefaults()
	os.Exit(1)
}

func Parse() {
	flag.Usage = usage

	// 获取环境变量 env
	env := flag.String("env", "dev",
		`Sets run environment. Possible values are "dev" and "prod"`)
	flag.Parse() // 解析参数

	// 修改日志输出类型
	logging.ConfigureLogger(*env)
	if *env == "prod" {
		logging.SetGinLogToFile()
	}

	fmt.Println("env = ", *env)
}
