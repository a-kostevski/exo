package fs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureDirectoryExists(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a file path within a directory that does not yet exist.
	testFile := filepath.Join(tmpDir, "nonexistent", "subdir", "file.txt")
	osfs := fs.NewOSFileSystem()

	err := osfs.EnsureDirectoryExists(testFile)
	require.NoError(t, err)

	// Check that the parent directory exists.
	parentDir := filepath.Dir(testFile)
	info, err := os.Stat(parentDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestWriteAndReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	osfs := fs.NewOSFileSystem()

	filePath := filepath.Join(tmpDir, "test.txt")
	content := []byte("Hello, world!")
	err := osfs.WriteFile(filePath, content)
	require.NoError(t, err)

	// Verify that the file exists.
	assert.True(t, osfs.FileExists(filePath))

	// Read the file and verify its content.
	readContent, err := osfs.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, content, readContent)
}

func TestAppendToFile(t *testing.T) {
	tmpDir := t.TempDir()
	osfs := fs.NewOSFileSystem()

	filePath := filepath.Join(tmpDir, "append.txt")
	initialContent := "Line 1"
	err := osfs.WriteFile(filePath, []byte(initialContent))
	require.NoError(t, err)

	appendContent := "Line 2"
	err = osfs.AppendToFile(filePath, appendContent)
	require.NoError(t, err)

	// Expected file content: initialContent + newline + appended line + newline.
	expected := initialContent + appendContent + "\n"
	readContent, err := osfs.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, expected, string(readContent))
}

func TestDeleteFile(t *testing.T) {
	tmpDir := t.TempDir()
	osfs := fs.NewOSFileSystem()

	filePath := filepath.Join(tmpDir, "delete.txt")
	err := osfs.WriteFile(filePath, []byte("to be deleted"))
	require.NoError(t, err)
	assert.True(t, osfs.FileExists(filePath))

	err = osfs.DeleteFile(filePath)
	require.NoError(t, err)
	assert.False(t, osfs.FileExists(filePath))
}

// TestOpenInEditor simulates opening a file in an editor by using a dummy editor script.
func TestOpenInEditor(t *testing.T) {
	tmpDir := t.TempDir()
	osfs := fs.NewOSFileSystem()

	// Create a dummy file that will be "opened."
	filePath := filepath.Join(tmpDir, "open.txt")
	content := "dummy content"
	err := osfs.WriteFile(filePath, []byte(content))
	require.NoError(t, err)

	// Create a dummy editor script that writes "opened" to a marker file.
	markerPath := filepath.Join(tmpDir, "marker.txt")
	dummyEditor := filepath.Join(tmpDir, "dummy_editor.sh")
	script := `#!/bin/sh
echo "opened" > "` + markerPath + `"
exit 0
`
	err = os.WriteFile(dummyEditor, []byte(script), 0755)
	require.NoError(t, err)

	// Call OpenInEditor using our dummy editor.
	err = osfs.OpenInEditor(filePath, dummyEditor)
	require.NoError(t, err)

	// Verify that the marker file exists and contains the expected text.
	markerContent, err := os.ReadFile(markerPath)
	require.NoError(t, err)
	assert.Equal(t, "opened\n", string(markerContent))
}

func TestReadDir_Success(t *testing.T) {
	// Create a temporary directory.
	tmpDir := t.TempDir()

	// Create some test files and directories.
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	subDir := filepath.Join(tmpDir, "subdir")
	require.NoError(t, os.MkdirAll(subDir, 0755))
	require.NoError(t, os.WriteFile(file1, []byte("content1"), 0644))
	require.NoError(t, os.WriteFile(file2, []byte("content2"), 0644))

	// Create an instance of OSFileSystem.
	fsys := fs.NewOSFileSystem()

	// Call ReadDir on the temporary directory.
	entries, err := fsys.ReadDir(tmpDir)
	require.NoError(t, err)

	// Extract names of returned entries.
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	// Verify that the created files and directory appear.
	assert.Contains(t, names, "file1.txt")
	assert.Contains(t, names, "file2.txt")
	assert.Contains(t, names, "subdir")
}

func TestReadDir_NonExistent(t *testing.T) {
	fsys := fs.NewOSFileSystem()
	_, err := fsys.ReadDir("nonexistent_directory_abcxyz")
	require.Error(t, err)
}
