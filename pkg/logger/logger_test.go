package logger_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureOutput temporarily redirects output (stdout or stderr) while f is executed,
// and then returns the captured output as a string.
func captureOutput(output *os.File, f func()) string {
	// Save the original file descriptor.
	orig := output
	// Create a pipe.
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Redirect output.
	if output == os.Stdout {
		os.Stdout = w
	} else {
		os.Stderr = w
	}

	// Execute the function.
	f()
	w.Close()

	// Read the captured output.
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		panic(err)
	}

	// Restore the original output.
	if output == os.Stdout {
		os.Stdout = orig
	} else {
		os.Stderr = orig
	}

	return buf.String()
}

func TestInfo(t *testing.T) {
	log := logger.NewLogger()
	output := captureOutput(os.Stdout, func() {
		log.Info("Test info", logger.Field{Key: "user", Value: "alice"})
	})

	// Verify that the output contains the info level tag, the message, and the field.
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "Test info")
	assert.Contains(t, output, "user=alice")
}

func TestInfof(t *testing.T) {
	log := logger.NewLogger()
	output := captureOutput(os.Stdout, func() {
		log.Infof("Infof: number %d", 42)
	})

	// Verify that the output contains the info level tag and the formatted message.
	assert.Contains(t, output, "[INFO]")
	assert.Contains(t, output, "Infof: number 42")
}

func TestError(t *testing.T) {
	log := logger.NewLogger()
	output := captureOutput(os.Stderr, func() {
		log.Error("Test error", logger.Field{Key: "code", Value: 500})
	})

	// Verify that the output contains the error level tag, the message, and the field.
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "Test error")
	assert.Contains(t, output, "code=500")
}

func TestErrorf(t *testing.T) {
	log := logger.NewLogger()
	output := captureOutput(os.Stderr, func() {
		log.Errorf("Errorf: %s", "something went wrong")
	})

	// Verify that the output contains the error level tag and the formatted message.
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "Errorf: something went wrong")
}

// Optionally, check that the timestamp appears to be a valid RFC3339 string.
func TestTimestampFormat(t *testing.T) {
	log := logger.NewLogger()

	output := captureOutput(os.Stdout, func() {
		log.Info("Timestamp test")
	})
	// Extract the timestamp between "[INFO] " and " - "
	// Example output: "[INFO] 2025-02-08T15:04:05Z07:00 - Timestamp test"
	parts := strings.SplitN(output, " - ", 2)
	require.Len(t, parts, 2)
	tsPart := strings.TrimPrefix(parts[0], "[INFO] ")
	_, err := time.Parse(time.RFC3339, strings.TrimSpace(tsPart))
	assert.NoError(t, err, "Timestamp should be in RFC3339 format")
}
