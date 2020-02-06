package mylog

import (
	"fmt"
	"gin-blog/pkg/file"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

var (
	F                  *os.File
	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}
)

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func Setup() {
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = file.MustOpen(fileName, filePath)
	if err != nil {
		log.Fatalln(err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func setLogPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}

	logger.SetPrefix(logPrefix)
}

func Debug(v ...interface{}) {
	setLogPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...interface{}) {
	setLogPrefix(INFO)
	logger.Println(v)
}

func Warning(v ...interface{}) {
	setLogPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...interface{}) {
	setLogPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...interface{}) {
	setLogPrefix(FATAL)
	logger.Fatalln(v)
}
