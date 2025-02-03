package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Field is a key-value pair for structured logging
type Field struct {
	Key   string
	Value interface{}
}

// Level represents the logging level
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
	FatalLevel Level = "fatal"
)

// Format represents the log output format
type Format string

const (
	JSONFormat Format = "json"
	TextFormat Format = "text"
)

// OutputType represents where logs should be written
type OutputType string

const (
	StdoutOutput  OutputType = "stdout"
	StderrOutput  OutputType = "stderr"
	FileOutput    OutputType = "file"
	DiscardOutput OutputType = "discard"
)

// Config holds logger configuration
type Config struct {
	Level      Level      `mapstructure:"level"`
	Format     Format     `mapstructure:"format"`
	Output     OutputType `mapstructure:"output"`
	File       string     `mapstructure:"file"`
	MaxSize    int        `mapstructure:"max_size"`    // Maximum size in megabytes before log rotation
	MaxBackups int        `mapstructure:"max_backups"` // Maximum number of old log files to retain
	MaxAge     int        `mapstructure:"max_age"`     // Maximum number of days to retain old log files
	Compress   bool       `mapstructure:"compress"`    // Compress rotated files
}

// Logger interface defines logging behavior
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}

type zapLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

var (
	defaultLogger Logger
	once          sync.Once
	mu            sync.RWMutex
)

func Default() Logger {
	if defaultLogger == nil {
		// Initialize with basic configuration
		defaultLogger = &zapLogger{
			logger: zap.NewExample(),
			ctx:    context.Background(),
		}
	}
	return defaultLogger
}

// Initialize sets up the logger with the given configuration
func Initialize(cfg Config) error {
	var err error
	once.Do(func() {
		defaultLogger, err = NewLogger(cfg)
	})
	return err
}

// Reinitialize allows changing the logger configuration after initialization
func Reinitialize(cfg Config) error {
	mu.Lock()
	defer mu.Unlock()

	logger, err := NewLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to create new logger: %w", err)
	}

	// If the previous logger was using a file, close it
	if l, ok := defaultLogger.(*zapLogger); ok && l != nil {
		_ = l.logger.Sync()
	}

	defaultLogger = logger
	once = sync.Once{}
	return nil
}

// NewLogger creates a new logger instance
func NewLogger(cfg Config) (Logger, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = getZapLevel(cfg.Level)
	zapCfg.Encoding = string(cfg.Format)

	// Configure output
	var output zapcore.WriteSyncer
	switch cfg.Output {
	case StdoutOutput:
		output = zapcore.AddSync(os.Stdout)
	case StderrOutput:
		output = zapcore.AddSync(os.Stderr)
	case FileOutput:
		if cfg.File == "" {
			return nil, fmt.Errorf("file output specified but no file path provided")
		}
		output = zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.File,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
	case DiscardOutput:
		output = zapcore.AddSync(io.Discard)
	default:
		return nil, fmt.Errorf("invalid output type: %s", cfg.Output)
	}

	// Create custom encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create core
	var encoder zapcore.Encoder
	if cfg.Format == JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(encoder, output, zapCfg.Level)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &zapLogger{logger: logger, ctx: context.Background()}, nil
}

func validateConfig(cfg Config) error {
	switch cfg.Level {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel:
	default:
		return fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	switch cfg.Format {
	case JSONFormat, TextFormat:
	default:
		return fmt.Errorf("invalid log format: %s", cfg.Format)
	}

	switch cfg.Output {
	case StdoutOutput, StderrOutput, FileOutput, DiscardOutput:
	default:
		return fmt.Errorf("invalid output type: %s", cfg.Output)
	}

	if cfg.Output == FileOutput && cfg.File == "" {
		return fmt.Errorf("file output specified but no file path provided")
	}

	return nil
}

func getZapLevel(level Level) zap.AtomicLevel {
	switch level {
	case DebugLevel:
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case InfoLevel:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WarnLevel:
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case ErrorLevel:
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case FatalLevel:
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}

// Implementation of Logger interface methods
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, convertFields(fields...)...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, convertFields(fields...)...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, convertFields(fields...)...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, convertFields(fields...)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, convertFields(fields...)...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(convertFields(fields...)...),
		ctx:    l.ctx,
	}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	return &zapLogger{
		logger: l.logger,
		ctx:    ctx,
	}
}

// Helper functions for the default logger
func Debug(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Fatal(msg, fields...)
	}
}

func Debugf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatalf(format, args...)
	}
}

// convertFields converts our Field type to zap.Field
func convertFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}
