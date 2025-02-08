package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	// Set a temporary HOME for the test.
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	pathWithTilde := "~/folder/file.txt"
	expanded := fs.ExpandPath(pathWithTilde)
	expected := filepath.Join(tmpHome, "folder", "file.txt")
	assert.Equal(t, expected, expanded)

	// When no tilde, path is returned unchanged.
	absolutePath := "/abs/path"
	assert.Equal(t, absolutePath, fs.ExpandPath(absolutePath))
}

func TestResolvePath(t *testing.T) {
	base := "/base/path"
	relative := "relative/file.txt"
	resolved := fs.ResolvePath(base, relative)
	expected := filepath.Join(base, relative)
	assert.Equal(t, expected, resolved)

	// When the given path is already absolute, it is returned as-is.
	abs := "/another/abs/file.txt"
	resolved = fs.ResolvePath(base, abs)
	assert.Equal(t, abs, resolved)
}

func TestSanitizePath(t *testing.T) {
	// Set a temporary HOME.
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", origHome)

	// Test a relative path.
	relPath := "folder//subfolder/../file.txt"
	sanitized := fs.SanitizePath(relPath, tmpHome)
	expected := filepath.Join(tmpHome, "folder", "file.txt")
	assert.Equal(t, expected, sanitized)

	// Test a path with tilde.
	tildePath := "~/folder/file.txt"
	sanitized = fs.SanitizePath(tildePath, tmpHome)

	expected = filepath.Join(tmpHome, "folder", "file.txt")
	assert.Equal(t, expected, sanitized)
}
