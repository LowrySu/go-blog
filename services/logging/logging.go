package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

const (
	logsDir = "logs"
	logName = "gin_production.log"
)

var logFilePath = filepath.Join(logsDir, logName)

func SetGinLogToFile() {

	// 把日志信息输出到文件内
	gin.SetMode(gin.ReleaseMode)
	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Panic().Err(err).Msg("Error opening Gin log file")
	}
	gin.DefaultWriter = io.MultiWriter(logFile)
}

func ConfigureLogger(env string) {
	// 设置日志级别
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// 判断环境,如果是dev环境则输出到std，如果是prod则输出到文件里
	switch env {
	case "dev":
		stdOutWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "12:00:00.000",
		}
		logger := zerolog.New(stdOutWriter).With().Timestamp().Logger()
		log.Logger = logger
	case "prod":
		// 如果文件夹不存在就创建文件夹
		createLogDir()
		// 把信息备份到到新的日志文件那日，清空原日志文件
		backupLastLog()
		// 获取原日志文件
		logFile := openLogFile()
		// 构建writer
		logFileWriter := zerolog.ConsoleWriter{
			Out:        logFile,
			NoColor:    true,
			TimeFormat: "15:04:05.000",
		}
		// 构建logger对象
		logger := zerolog.New(logFileWriter).With().Timestamp().Logger()
		// 设置logger
		log.Logger = logger
	default:
		fmt.Printf("Env not valid: %s\n", env)
		os.Exit(2)
	}
}

func createLogDir() {
	// 创建文件夹
	if err := os.Mkdir(logsDir, 0744); err != nil && !os.IsExist(err) {
		log.Fatal().Err(err).Msg("Unable to create logs directory.")
	}
}

func backupLastLog() {
	timeStamp := time.Now().Format("20060201_15_04_05")
	base := strings.TrimSuffix(logName, filepath.Ext(logName))
	bkpLogName := base + "_" + timeStamp + "." + filepath.Ext(logName)
	bkpLogPath := filepath.Join(logsDir, bkpLogName)

	logFile, err := ioutil.ReadFile(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Panic().Err(err).Msg("Error reading log file for backup")
	}

	if err = ioutil.WriteFile(bkpLogPath, logFile, 0644); err != nil {
		log.Panic().Err(err).Msg("Error writing backup log file")
	}
}

func openLogFile() *os.File {
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Panic().Err(err).Msg("Error while opening log file")
	}
	return logFile
}

func curentDir() string {
	path, err := os.Executable()
	if err != nil {
		log.Panic().Err(err).Msg("Can’t get current directory.")
	}
	return filepath.Dir(path)
}
