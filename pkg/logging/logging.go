package logging

import (
	"log"
	"os"
)

var (
	logger      *log.Logger
	errorLogger *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR ", log.Ldate|log.Ltime)
}

func Info(msg string) {
	logger.Println(msg)
}

func Error(msg string) {
	errorLogger.Println(msg)
}
