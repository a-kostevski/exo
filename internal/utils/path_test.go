package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandPath(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set up a test home directory
	tmpDir := t.TempDir()
	require.NoError(t, os.Setenv("HOME", tmpDir))

	tests := []struct {
		name     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "expand home directory",
			path:     "~/test",
			expected: filepath.Join(tmpDir, "test"),
			wantErr:  false,
		},
		{
			name:     "path without tilde",
			path:     "/absolute/path",
			expected: "/absolute/path",
			wantErr:  false,
		},
		{
			name:     "relative path",
			path:     "relative/path",
			expected: "relative/path",
			wantErr:  false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "tilde in middle of path",
			path:     "/path/~/test",
			expected: "/path/~/test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expanded := ExpandPath(tt.path)
			assert.Equal(t, tt.expected, expanded)
		})
	}
}

func TestResolvePath(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set up a test home directory
	tmpDir := t.TempDir()
	require.NoError(t, os.Setenv("HOME", tmpDir))

	tests := []struct {
		name     string
		base     string
		path     string
		expected string
		wantErr  bool
	}{
		{
			name:     "absolute path",
			base:     tmpDir,
			path:     "/absolute/path",
			expected: "/absolute/path",
			wantErr:  false,
		},
		{
			name:     "relative path",
			base:     tmpDir,
			path:     "relative/path",
			expected: filepath.Join(tmpDir, "relative/path"),
			wantErr:  false,
		},
		{
			name:     "path with tilde",
			base:     tmpDir,
			path:     "~/test",
			expected: filepath.Join(tmpDir, "test"),
			wantErr:  false,
		},
		{
			name:     "empty path",
			base:     tmpDir,
			path:     "",
			expected: tmpDir,
			wantErr:  false,
		},
		{
			name:     "empty base",
			base:     "",
			path:     "test",
			expected: "test",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// First expand any tildes in the path
			expandedPath := ExpandPath(tt.path)
			resolved := ResolvePath(tt.base, expandedPath)
			assert.Equal(t, tt.expected, resolved)
		})
	}
}

func TestGetXDGConfigHome(t *testing.T) {
	// Save original environment variables
	originalHome := os.Getenv("HOME")
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)
	}()

	// Set up a test home directory
	tmpDir := t.TempDir()
	require.NoError(t, os.Setenv("HOME", tmpDir))

	tests := []struct {
		name           string
		xdgConfigHome  string
		expectedSuffix string
	}{
		{
			name:           "XDG_CONFIG_HOME set",
			xdgConfigHome:  filepath.Join(tmpDir, "custom/config"),
			expectedSuffix: "custom/config",
		},
		{
			name:           "XDG_CONFIG_HOME not set",
			xdgConfigHome:  "",
			expectedSuffix: ".config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.xdgConfigHome != "" {
				require.NoError(t, os.Setenv("XDG_CONFIG_HOME", tt.xdgConfigHome))
			} else {
				require.NoError(t, os.Unsetenv("XDG_CONFIG_HOME"))
			}

			configHome := GetXDGConfigHome()
			if tt.xdgConfigHome != "" {
				assert.Equal(t, tt.xdgConfigHome, configHome)
			} else {
				assert.Equal(t, filepath.Join(tmpDir, tt.expectedSuffix), configHome)
			}
		})
	}
}

func TestGetXDGDataHome(t *testing.T) {
	// Save original environment variables
	originalHome := os.Getenv("HOME")
	originalXDGDataHome := os.Getenv("XDG_DATA_HOME")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("XDG_DATA_HOME", originalXDGDataHome)
	}()

	// Set up a test home directory
	tmpDir := t.TempDir()
	require.NoError(t, os.Setenv("HOME", tmpDir))

	tests := []struct {
		name           string
		xdgDataHome    string
		expectedSuffix string
	}{
		{
			name:           "XDG_DATA_HOME set",
			xdgDataHome:    filepath.Join(tmpDir, "custom/data"),
			expectedSuffix: "custom/data",
		},
		{
			name:           "XDG_DATA_HOME not set",
			xdgDataHome:    "",
			expectedSuffix: ".local/share",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.xdgDataHome != "" {
				require.NoError(t, os.Setenv("XDG_DATA_HOME", tt.xdgDataHome))
			} else {
				require.NoError(t, os.Unsetenv("XDG_DATA_HOME"))
			}

			dataHome := GetXDGDataHome()
			if tt.xdgDataHome != "" {
				assert.Equal(t, tt.xdgDataHome, dataHome)
			} else {
				assert.Equal(t, filepath.Join(tmpDir, tt.expectedSuffix), dataHome)
			}
		})
	}
}
