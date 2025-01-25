package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileOperations(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("EnsureDirectoryExists creates directories", func(t *testing.T) {
		path := filepath.Join(tmpDir, "test", "nested", "file.txt")
		err := EnsureDirectoryExists(path)
		require.NoError(t, err)

		// Check that the directory was created
		dirInfo, err := os.Stat(filepath.Dir(path))
		require.NoError(t, err)
		assert.True(t, dirInfo.IsDir())
	})

	t.Run("FileExists checks file existence", func(t *testing.T) {
		nonExistentPath := filepath.Join(tmpDir, "nonexistent.txt")
		assert.False(t, FileExists(nonExistentPath))

		existingPath := filepath.Join(tmpDir, "existing.txt")
		require.NoError(t, WriteFile(existingPath, []byte("test")))
		assert.True(t, FileExists(existingPath))
	})

	t.Run("CreateFile creates empty file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "empty.txt")
		err := CreateFile(path)
		require.NoError(t, err)

		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.Equal(t, int64(0), info.Size())
	})

	t.Run("WriteFile writes content", func(t *testing.T) {
		path := filepath.Join(tmpDir, "write.txt")
		content := []byte("test content")
		err := WriteFile(path, content)
		require.NoError(t, err)

		// Read and verify content
		readContent, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, content, readContent)
	})

	t.Run("ReadFile reads content", func(t *testing.T) {
		path := filepath.Join(tmpDir, "read.txt")
		content := []byte("test content")
		require.NoError(t, os.WriteFile(path, content, 0644))

		readContent, err := ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, content, readContent)
	})

	t.Run("AppendToFile appends content", func(t *testing.T) {
		path := filepath.Join(tmpDir, "append.txt")

		// First write
		err := AppendToFile(path, "line1")
		require.NoError(t, err)

		// Second write
		err = AppendToFile(path, "line2")
		require.NoError(t, err)

		// Read and verify content
		content, err := os.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "line1\nline2\n", string(content))
	})

	t.Run("DeleteFile removes file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "to_delete.txt")
		require.NoError(t, WriteFile(path, []byte("test")))
		assert.True(t, FileExists(path))

		require.NoError(t, DeleteFile(path))
		assert.False(t, FileExists(path))
	})

	t.Run("DeleteFile handles non-existent file", func(t *testing.T) {
		path := filepath.Join(tmpDir, "nonexistent.txt")
		err := DeleteFile(path)
		assert.NoError(t, err)
	})

	t.Run("operations with empty path", func(t *testing.T) {
		assert.Error(t, AppendToFile("", "content"))
		assert.False(t, FileExists(""))
	})

	t.Run("operations with invalid paths", func(t *testing.T) {
		invalidPath := filepath.Join(string([]byte{0}), "invalid.txt")
		assert.Error(t, WriteFile(invalidPath, []byte("test")))
		assert.Error(t, CreateFile(invalidPath))
		_, err := ReadFile(invalidPath)
		assert.Error(t, err)
	})
}

func TestEnsureDirectories(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Set up a test home directory
	tmpDir := t.TempDir()
	require.NoError(t, os.Setenv("HOME", tmpDir))

	t.Run("create multiple directories", func(t *testing.T) {
		paths := []string{
			filepath.Join(tmpDir, "dir1"),
			filepath.Join(tmpDir, "dir2", "subdir"),
			filepath.Join(tmpDir, "dir3", "subdir", "subsubdir"),
		}

		err := EnsureDirectories(paths...)
		require.NoError(t, err)

		// Verify directories were created
		for _, path := range paths {
			info, err := os.Stat(path)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
			assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
		}
	})

	t.Run("create directories with tilde paths", func(t *testing.T) {
		paths := []string{
			"~/test/dir1",
			"~/test/dir2/subdir",
		}

		err := EnsureDirectories(paths...)
		require.NoError(t, err)

		// Verify directories were created
		for _, path := range paths {
			expandedPath := filepath.Join(tmpDir, path[2:]) // Remove ~/ and join with tmpDir
			info, err := os.Stat(expandedPath)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
			assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
		}
	})

	t.Run("create directories with relative paths", func(t *testing.T) {
		paths := []string{
			"./test/dir1",
			"test/dir2/subdir",
		}

		err := EnsureDirectories(paths...)
		require.NoError(t, err)

		// Verify directories were created
		for _, path := range paths {
			fullPath := filepath.Join(tmpDir, filepath.Clean(path))
			info, err := os.Stat(fullPath)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
			assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
		}
	})

	t.Run("ensure existing directories", func(t *testing.T) {
		path := filepath.Join(tmpDir, "existing")
		require.NoError(t, os.MkdirAll(path, 0755))

		err := EnsureDirectories(path)
		assert.NoError(t, err)

		// Verify directory still exists with correct permissions
		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
		assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
	})

	t.Run("permission denied", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("Skipping test when running as root")
		}

		// Create a read-only parent directory
		parent := filepath.Join(tmpDir, "readonly")
		require.NoError(t, os.MkdirAll(parent, 0555))

		path := filepath.Join(parent, "newdir")
		err := EnsureDirectories(path)
		assert.Error(t, err)
	})
}
