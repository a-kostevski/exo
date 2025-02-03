package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "default config",
			cfg: Config{
				Level:  InfoLevel,
				Format: TextFormat,
				Output: StderrOutput,
			},
			wantErr: false,
		},
		{
			name: "json format",
			cfg: Config{
				Level:  DebugLevel,
				Format: JSONFormat,
				Output: StdoutOutput,
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			cfg: Config{
				Level:  Level("invalid"),
				Format: TextFormat,
				Output: StderrOutput,
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			cfg: Config{
				Level:  InfoLevel,
				Format: Format("invalid"),
				Output: StderrOutput,
			},
			wantErr: true,
		},
		{
			name: "invalid output",
			cfg: Config{
				Level:  InfoLevel,
				Format: TextFormat,
				Output: OutputType("invalid"),
			},
			wantErr: true,
		},
		{
			name: "file output without path",
			cfg: Config{
				Level:  InfoLevel,
				Format: TextFormat,
				Output: FileOutput,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset logger before each test
			defaultLogger = nil
			once = sync.Once{}

			err := Initialize(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, defaultLogger)
			}
		})
	}
}

func TestFileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := Config{
		Level:  InfoLevel,
		Format: TextFormat,
		Output: FileOutput,
		File:   logFile,
	}

	// Reset logger before test
	defaultLogger = nil
	once = sync.Once{}

	err := Initialize(cfg)
	require.NoError(t, err)

	Info("test message")

	// Ensure logger is flushed
	if l, ok := defaultLogger.(*zapLogger); ok {
		require.NoError(t, l.logger.Sync())
	}

	content, err := os.ReadFile(logFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test message")
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		messages []struct {
			logFunc   func(msg string, fields ...Field)
			message   string
			shouldLog bool
		}
	}{
		{
			name:  "info level",
			level: InfoLevel,
			messages: []struct {
				logFunc   func(msg string, fields ...Field)
				message   string
				shouldLog bool
			}{
				{Debug, "debug message", false},
				{Info, "info message", true},
				{Warn, "warn message", true},
				{Error, "error message", true},
			},
		},
		{
			name:  "debug level",
			level: DebugLevel,
			messages: []struct {
				logFunc   func(msg string, fields ...Field)
				message   string
				shouldLog bool
			}{
				{Debug, "debug message", true},
				{Info, "info message", true},
				{Warn, "warn message", true},
				{Error, "error message", true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset logger before each test
			defaultLogger = nil
			once = sync.Once{}

			var buf bytes.Buffer
			err := Initialize(Config{
				Level:  tt.level,
				Format: TextFormat,
				Output: StderrOutput,
			})
			require.NoError(t, err)

			// Replace stderr with our buffer
			if l, ok := defaultLogger.(*zapLogger); ok {
				l.logger = l.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
					return zapcore.NewCore(
						zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
						zapcore.AddSync(&buf),
						getZapLevel(tt.level),
					)
				}))
			}

			for _, msg := range tt.messages {
				buf.Reset()
				msg.logFunc(msg.message)
				if msg.shouldLog {
					assert.Contains(t, buf.String(), msg.message)
				} else {
					assert.NotContains(t, buf.String(), msg.message)
				}
			}
		})
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer

	// Reset logger before test
	defaultLogger = nil
	once = sync.Once{}

	err := Initialize(Config{
		Level:  InfoLevel,
		Format: JSONFormat,
		Output: StderrOutput,
	})
	require.NoError(t, err)

	// Replace stderr with our buffer
	if l, ok := defaultLogger.(*zapLogger); ok {
		l.logger = l.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(&buf),
				getZapLevel(InfoLevel),
			)
		}))
	}

	fields := []Field{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: 123},
	}

	defaultLogger.With(fields...).Info("test with fields")

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	assert.Equal(t, "value1", logEntry["key1"])
	assert.Equal(t, float64(123), logEntry["key2"]) // JSON numbers are float64
	assert.Equal(t, "test with fields", logEntry["msg"])
}

func TestFormatted(t *testing.T) {
	var buf bytes.Buffer

	// Reset logger before test
	defaultLogger = nil
	once = sync.Once{}

	err := Initialize(Config{
		Level:  DebugLevel,
		Format: TextFormat,
		Output: StderrOutput,
	})
	require.NoError(t, err)

	// Replace stderr with our buffer
	if l, ok := defaultLogger.(*zapLogger); ok {
		l.logger = l.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.AddSync(&buf),
				getZapLevel(DebugLevel),
			)
		}))
	}

	tests := []struct {
		name     string
		logFunc  func(format string, args ...interface{})
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "Debugf",
			logFunc:  Debugf,
			format:   "debug %s",
			args:     []interface{}{"message"},
			expected: "debug message",
		},
		{
			name:     "Infof",
			logFunc:  Infof,
			format:   "info %s",
			args:     []interface{}{"message"},
			expected: "info message",
		},
		{
			name:     "Warnf",
			logFunc:  Warnf,
			format:   "warn %s",
			args:     []interface{}{"message"},
			expected: "warn message",
		},
		{
			name:     "Errorf",
			logFunc:  Errorf,
			format:   "error %s",
			args:     []interface{}{"message"},
			expected: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.format, tt.args...)
			assert.Contains(t, buf.String(), tt.expected)
		})
	}
}

func TestWithContext(t *testing.T) {
	// Reset logger before test
	defaultLogger = nil
	once = sync.Once{}

	err := Initialize(Config{
		Level:  InfoLevel,
		Format: TextFormat,
		Output: StderrOutput,
	})
	require.NoError(t, err)

	ctx := context.WithValue(context.Background(), "test-key", "test-value")
	logger := defaultLogger.WithContext(ctx)

	// Verify that the context is stored
	if l, ok := logger.(*zapLogger); ok {
		assert.Equal(t, "test-value", l.ctx.Value("test-key"))
	}
}

func TestReinitialize(t *testing.T) {
	// Reset logger before test
	defaultLogger = nil
	once = sync.Once{}

	// Initial configuration
	err := Initialize(Config{
		Level:  InfoLevel,
		Format: TextFormat,
		Output: StderrOutput,
	})
	require.NoError(t, err)

	// Get the initial logger instance
	initial := defaultLogger

	// Reinitialize with new configuration
	err = Reinitialize(Config{
		Level:  DebugLevel,
		Format: JSONFormat,
		Output: StdoutOutput,
	})
	require.NoError(t, err)

	// Verify that we have a new logger instance
	assert.NotEqual(t, initial, defaultLogger)

	// Verify that the new configuration is applied by checking if debug messages are logged
	var buf bytes.Buffer
	if l, ok := defaultLogger.(*zapLogger); ok {
		l.logger = l.logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(&buf),
				getZapLevel(DebugLevel),
			)
		}))
	}

	Debug("test debug message")
	assert.Contains(t, buf.String(), "test debug message")
}

func TestMultipleInitialize(t *testing.T) {
	// Reset logger before test
	defaultLogger = nil
	once = sync.Once{}

	// First initialization
	err1 := Initialize(Config{
		Level:  InfoLevel,
		Format: TextFormat,
		Output: StderrOutput,
	})
	require.NoError(t, err1)
	initial := defaultLogger

	// Second initialization (should not change the logger)
	err2 := Initialize(Config{
		Level:  DebugLevel,
		Format: JSONFormat,
		Output: StdoutOutput,
	})
	require.NoError(t, err2)

	// Verify that the logger instance hasn't changed
	assert.Equal(t, initial, defaultLogger)
}
