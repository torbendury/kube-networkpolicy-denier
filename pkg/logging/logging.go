// Logging package. This package wraps the standard log package and provides global logging functions.
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

// Info logs a message to stdout. The message is prefixed with "INFO" and the current date and time.
func Info(msg string) {
	logger.Println(msg)
}

// Error logs a message to stderr. The message is prefixed with "ERROR" and the current date and time.
func Error(msg string) {
	errorLogger.Println(msg)
}
