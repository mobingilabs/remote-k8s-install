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
	funcN, file, line, ok := runtime.Caller(2)
	if ok {
		logger.Printf("func:%s,file:%s,line:%s", funcN, file, line)
	}

	logger.Println(a...)
}

func Errorf(format string, a ...interface{}) {
	funcN, file, line, ok := runtime.Caller(2)
	if ok {
		logger.Printf("func:%s,file:%s,line:%s", funcN, file, line)
	}

	logger.Printf(format, a...)
}
