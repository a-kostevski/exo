package logger

import (
	"fmt"
	"os"
	"time"
)

// Field represents a key/value pair for structured logging.
type Field struct {
	Key   string
	Value interface{}
}

// Logger is our minimal logging interface with two levels (Info and Error)
// plus formatted versions of these methods.
type Logger interface {
	// Info logs an informational message.
	Info(msg string, fields ...Field)
	// Error logs an error message.
	Error(msg string, fields ...Field)
	// Infof logs a formatted informational message.
	Infof(format string, args ...interface{})
	// Errorf logs a formatted error message.
	Errorf(format string, args ...interface{})
}

// simpleLogger is a basic implementation of Logger.
type simpleLogger struct{}

// NewLogger creates a new instance of a Logger.
func NewLogger() Logger {
	return &simpleLogger{}
}

// Info logs an informational message to stdout.
func (l *simpleLogger) Info(msg string, fields ...Field) {
	timestamp := time.Now().Format(time.RFC3339)
	line := fmt.Sprintf("[INFO] %s - %s", timestamp, msg)
	if len(fields) > 0 {
		line += " " + formatFields(fields)
	}
	fmt.Fprintln(os.Stdout, line)
}

// Error logs an error message to stderr.
func (l *simpleLogger) Error(msg string, fields ...Field) {
	timestamp := time.Now().Format(time.RFC3339)
	line := fmt.Sprintf("[ERROR] %s - %s", timestamp, msg)
	if len(fields) > 0 {
		line += " " + formatFields(fields)
	}
	fmt.Fprintln(os.Stderr, line)
}

// Infof logs a formatted informational message.
func (l *simpleLogger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message.
func (l *simpleLogger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// formatFields converts a slice of Field into a formatted string.
func formatFields(fields []Field) string {
	var s string
	for _, f := range fields {
		s += fmt.Sprintf("%s=%v ", f.Key, f.Value)
	}
	return s
}
