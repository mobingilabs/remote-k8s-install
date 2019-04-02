package log

import (
	"log"
	"os"
	"runtime"
)

// now we need log promptly
//var logger *zap.SugaredLogger
var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Lshortfile)
}

func Info(a ...interface{}) {
	logger.Println(a...)
}

func Infof(format string, a ...interface{}) {
	logger.Printf(format, a...)
}

func Error(a ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		logger.Printf("file:%s,line:%d", file, line)
	}

	logger.Println(a...)
}

func Errorf(format string, a ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		logger.Printf("file:%s,line:%d", file, line)
	}

	logger.Printf(format, a...)
}

func Panic(a ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		logger.Printf("file:%s,line:%d", file, line)
	}

	logger.Println(a...)
	panic("")
}

func Panicf(format string, a ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		logger.Printf("file:%s,line:%d", file, line)
	}

	logger.Printf(format, a...)
	panic("")
}
