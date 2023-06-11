package logger

import (
	"log"
	"os"
	"sync"
)

var (
	globalLogger    *log.Logger
	muxGlobalLogger sync.RWMutex
)

func getGlobalLogger() *log.Logger {
	if globalLogger == nil {
		initGlobalLogger()
	}

	muxGlobalLogger.RLock()
	defer muxGlobalLogger.RUnlock()

	return globalLogger
}

func initGlobalLogger() {
	muxGlobalLogger.Lock()
	defer muxGlobalLogger.Unlock()

	globalLogger = NewLogger()
}

func NewLogger() *log.Logger {
	file, err := os.OpenFile("./partsy.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(message string) {
	logger := getGlobalLogger()
	logger.Printf("[DBG] %s\n", message)
}

func Info(message string) {
	logger := getGlobalLogger()
	logger.Printf("[INFO] %s\n", message)
}

func Warn(message string) {
	logger := getGlobalLogger()
	logger.Printf("[WARN] %s\n", message)
}

func Error(err error, message string) {
	logger := getGlobalLogger()
	logger.Printf("[ERR] %s: %s\n", message, err)
}

func Critical(message string) {
	logger := getGlobalLogger()
	logger.Printf("[CRIT] %s\n", message)
}

func Fatal(err error, message string) {
	logger := getGlobalLogger()
	logger.Printf("[FATAL] %s: %s\n", message, err)
}
