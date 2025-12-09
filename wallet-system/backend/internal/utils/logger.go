package utils

import (
	"log"
	"os"
)

// Logger handles logging
type Logger struct {
	level string
}

// NewLogger creates a new logger
func NewLogger(level string) *Logger {
	return &Logger{level: level}
}

// Info logs info level messages
func (l *Logger) Info(msg string, args ...interface{}) {
	log.Printf("[INFO] "+msg+"\n", args...)
}

// Error logs error level messages
func (l *Logger) Error(msg string, args ...interface{}) {
	log.Printf("[ERROR] "+msg+"\n", args...)
}

// Debug logs debug level messages
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level == "debug" {
		log.Printf("[DEBUG] "+msg+"\n", args...)
	}
}

// Warn logs warning level messages
func (l *Logger) Warn(msg string, args ...interface{}) {
	log.Printf("[WARN] "+msg+"\n", args...)
}

// Fatal logs fatal level messages and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	log.Fatalf("[FATAL] "+msg+"\n", args...)
	os.Exit(1)
}
