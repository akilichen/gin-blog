package mylog

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFullPath() string {
	//prefixLogPath := getLogFilePath()
	//suffixLogPath := fmt.Sprintf("%s%s.%s", LogSavePath, time.Now().Format(TimeFormat), LogFileExt)

	fullLogPath := fmt.Sprintf("%s%s.%s", LogSavePath, time.Now().Format(TimeFormat), LogFileExt)
	return fmt.Sprintf("%s", fullLogPath)
}

func openLogFile(filepath string) *os.File {
	_, err := os.Stat(filepath) // 检查文件是否存在
	switch {
	case os.IsNotExist(err):
		MkDir()
	case os.IsPermission(err):
		log.Fatalf("Premission denied! %v", err)
	}

	handle, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to open file,%v", err)
	}
	return handle
}

func MkDir() {
	dir, _ := os.Getwd() // 获取当前路径
	//创建目录当目录不存在，包括父目录+子目录皆可
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm) // 最后一个参数为os定义常量，== 0777
	if err != nil {
		panic(err)
	}
}
