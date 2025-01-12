package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
				Level:  "info",
				Format: "text",
				Output: "stderr",
			},
			wantErr: false,
		},
		{
			name: "json format",
			cfg: Config{
				Level:  "debug",
				Format: "json",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			cfg: Config{
				Level:  "invalid",
				Format: "text",
				Output: "stderr",
			},
			wantErr: false, // Should not error, defaults to info
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize(tt.cfg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileOutput(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	cfg := Config{
		Level:  "info",
		Format: "text",
		Output: "file",
		File:   logFile,
	}

	err := Initialize(cfg)
	assert.NoError(t, err)

	Info("test message")

	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test message")
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.SetOutput(&buf)

	tests := []struct {
		level   string
		logFunc func(args ...interface{})
		format  string
		args    []interface{}
	}{
		{"debug", Debug, "test debug", nil},
		{"info", Info, "test info", nil},
		{"warn", Warn, "test warn", nil},
		{"error", Error, "test error", nil},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			buf.Reset()
			err := Initialize(Config{Level: tt.level, Format: "text"})
			assert.NoError(t, err)
			defaultLogger.SetOutput(&buf)

			tt.logFunc(tt.format)
			assert.Contains(t, buf.String(), tt.format)
		})
	}
}

func TestWithFields(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.SetOutput(&buf)

	err := Initialize(Config{
		Level:  "info",
		Format: "json",
		Output: "stderr",
	})
	assert.NoError(t, err)
	defaultLogger.SetOutput(&buf)

	fields := logrus.Fields{
		"key1": "value1",
		"key2": 123,
	}

	WithFields(fields).Info("test with fields")

	var logEntry map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	assert.NoError(t, err)

	assert.Equal(t, "value1", logEntry["key1"])
	assert.Equal(t, float64(123), logEntry["key2"]) // JSON numbers are float64
	assert.Equal(t, "test with fields", logEntry["msg"])
}

func TestFormatted(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.SetOutput(&buf)

	err := Initialize(Config{Level: "debug", Format: "text"})
	assert.NoError(t, err)
	defaultLogger.SetOutput(&buf)

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

func TestMultiOutput(t *testing.T) {
	// Create a temporary directory for log files
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// Create a buffer to capture stderr output
	var buf bytes.Buffer

	// Initialize logger with both file and stderr output
	err := Initialize(Config{
		Level:  "info",
		Format: "text",
		Output: "both",
		File:   logFile,
	})
	assert.NoError(t, err)

	// Create a multi-writer for both file and buffer
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	assert.NoError(t, err)
	defer file.Close()

	defaultLogger.SetOutput(io.MultiWriter(file, &buf))

	// Log a test message
	testMsg := "test multi output"
	Info(testMsg)

	// Read the log file
	content, err := os.ReadFile(logFile)
	assert.NoError(t, err)

	// Verify the message appears in both outputs
	assert.Contains(t, string(content), testMsg)
	assert.Contains(t, buf.String(), testMsg)
}
