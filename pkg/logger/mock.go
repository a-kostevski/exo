package logger

import (
	"context"
	"testing"
)

// Mock logger for testing
// MockLogger implements Logger interface for testing
type MockLogger struct {
	t *testing.T // optional, for test assertions
}

func NewMockLogger(t *testing.T) *MockLogger {
	return &MockLogger{t: t}
}

// Basic logging methods
func (m *MockLogger) Debug(msg string, fields ...Field) {}
func (m *MockLogger) Info(msg string, fields ...Field)  {}
func (m *MockLogger) Warn(msg string, fields ...Field)  {}
func (m *MockLogger) Error(msg string, fields ...Field) {}
func (m *MockLogger) Fatal(msg string, fields ...Field) {
	if m.t != nil {
		m.t.Fatal(msg)
	}
}

// Format logging methods
func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Infof(format string, args ...interface{})  {}
func (m *MockLogger) Warnf(format string, args ...interface{})  {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	if m.t != nil {
		m.t.Fatalf(format, args...)
	}
}

func (m *MockLogger) With(fields ...Field) Logger {
	return m
}

func (m *MockLogger) WithContext(ctx context.Context) Logger {
	return m
}
