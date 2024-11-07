package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	// 31是红色的代码
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
	// 34是蓝色的代码
	infoLog = log.New(os.Stdout, "\033[34m[info ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers = []*log.Logger{errorLog, infoLog}
	mu      sync.Mutex
)

// log method
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// log levels
const (
	InfoLevel = iota // 用于声明const中的连续整数值，通常用于枚举类型的实现
	ErrorLevel
	DisabledLevel
)

// 设置日志打印的等级
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout) // 将log输出重定位到控制台的标准输出上
	}

	if ErrorLevel < level { // 如果log的级别高于ErrorLevel则将ErrorLog重定向到io.Discard
		errorLog.SetOutput(io.Discard)
	}
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
