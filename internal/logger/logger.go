package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

// Field is a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Logger is a simple logging interface.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
}

var (
	// defaultLogger is the default logger instance
	defaultLogger     *logrus.Logger
	defaultLoggerOnce sync.Once
)

// Config holds the logger configuration
type Config struct {
	// Level is the logging level (debug, info, warn, error)
	Level string
	// Format is the output format (text, json)
	Format string
	// Output is where the logs will be written to
	Output string
	// File is the log file path (if Output is "file")
	File string
}

// Initialize sets up the logger with the given configuration
func Initialize(cfg Config) error {
	defaultLoggerOnce.Do(func() {
		defaultLogger = logrus.New()
	})

	// Set logging level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	defaultLogger.SetLevel(level)

	// Set output format
	switch cfg.Format {
	case "json":
		defaultLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		defaultLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output destination
	switch cfg.Output {
	case "file":
		if cfg.File != "" {
			fileWriter, err := setupFileOutput(cfg.File)
			if err != nil {
				return fmt.Errorf("failed to setup file output: %w", err)
			}
			defaultLogger.SetOutput(fileWriter)
		}
	case "stdout":
		defaultLogger.SetOutput(os.Stdout)
	case "stderr":
		defaultLogger.SetOutput(os.Stderr)
	case "both":
		if cfg.File != "" {
			// Create a multi-writer for both file and stderr
			fileWriter, err := setupFileOutput(cfg.File)
			if err != nil {
				return fmt.Errorf("failed to setup file output: %w", err)
			}
			defaultLogger.SetOutput(io.MultiWriter(os.Stderr, fileWriter))
		} else {
			defaultLogger.SetOutput(os.Stderr)
		}
	default:
		defaultLogger.SetOutput(os.Stderr)
	}

	return nil
}

// setupFileOutput creates the log directory and opens the log file
func setupFileOutput(filePath string) (io.Writer, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return f, nil
}

// WithFields creates a new entry with the given fields
func WithFields(fields map[string]interface{}) *logrus.Entry {
	return defaultLogger.WithFields(logrus.Fields(fields))
}

// SetOutput sets the logger output
func SetOutput(out io.Writer) {
	defaultLogger.SetOutput(out)
}

// Debug logs a message at level Debug
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Debugf logs a formatted message at level Debug
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Info logs a message at level Info
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Infof logs a formatted message at level Info
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warn logs a message at level Warn
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Warnf logs a formatted message at level Warn
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Error logs a message at level Error
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Errorf logs a formatted message at level Error
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Fatal logs a message at level Fatal then the process will exit with status set to 1
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Fatalf logs a formatted message at level Fatal then the process will exit with status set to 1
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}
