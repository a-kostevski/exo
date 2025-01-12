package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupInstallTest(t *testing.T) string {
	return t.TempDir()
}

func TestInstallDefault(t *testing.T) {
	tmpDir := setupInstallTest(t)

	tests := []struct {
		name    string
		opts    InstallOptions
		wantErr bool
		setup   func(string) error
	}{
		{
			name: "install to empty directory",
			opts: InstallOptions{
				TemplateDir: tmpDir,
				Force:       false,
				Reader:      NewMockInputReader([]string{"s"}),
			},
			wantErr: false,
		},
		{
			name: "install to empty path",
			opts: InstallOptions{
				TemplateDir: "",
				Force:       false,
				Reader:      NewMockInputReader([]string{"s"}),
			},
			wantErr: true,
		},
		{
			name: "force install with existing files",
			opts: InstallOptions{
				TemplateDir: tmpDir,
				Force:       true,
				Reader:      NewMockInputReader([]string{"s"}),
			},
			setup: func(dir string) error {
				testFile := filepath.Join(dir, "note.md")
				return os.WriteFile(testFile, []byte("existing content"), defaultPerms)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				err := tt.setup(tt.opts.TemplateDir)
				require.NoError(t, err)
			}

			err := InstallDefaultWithOptions(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify files were installed
			entries, err := templatesFS.ReadDir(defaultTemplateDir)
			require.NoError(t, err)

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				// Check if file exists in destination
				destPath := filepath.Join(tt.opts.TemplateDir, entry.Name())
				_, err := os.Stat(destPath)
				assert.NoError(t, err, "template file should exist: %s", entry.Name())

				// Only check backup for the specific test file we created
				if tt.opts.Force && entry.Name() == "note.md" {
					backupPath := destPath + backupExtension
					_, err := os.Stat(backupPath)
					assert.NoError(t, err, "backup file should exist: %s", backupPath)

					// Verify backup content
					content, err := os.ReadFile(backupPath)
					assert.NoError(t, err)
					assert.Equal(t, "existing content", string(content), "backup should contain original content")
				}
			}
		})
	}
}

func TestFindExistingFiles(t *testing.T) {
	tmpDir := setupInstallTest(t)

	// Create some test files
	existingFiles := []string{"test1.md", "test2.md"}
	for _, file := range existingFiles {
		err := os.WriteFile(filepath.Join(tmpDir, file), []byte("test"), defaultPerms)
		require.NoError(t, err)
	}

	// Create test entries
	entries, err := templatesFS.ReadDir(defaultTemplateDir)
	require.NoError(t, err)

	// Test finding existing files
	found, err := findExistingFiles(tmpDir, entries)
	assert.NoError(t, err)

	// Verify results
	for _, file := range found {
		assert.FileExists(t, filepath.Join(tmpDir, file))
	}
}

func TestCreateBackup(t *testing.T) {
	tmpDir := setupInstallTest(t)
	testFile := filepath.Join(tmpDir, "test.md")
	testContent := "test content"

	// Create test file
	err := os.WriteFile(testFile, []byte(testContent), defaultPerms)
	require.NoError(t, err)

	// Create backup
	err = createBackup(testFile)
	assert.NoError(t, err)

	// Verify backup was created
	backupFile := testFile + backupExtension
	assert.FileExists(t, backupFile)

	// Verify backup content
	content, err := os.ReadFile(backupFile)
	assert.NoError(t, err)
	assert.Equal(t, testContent, string(content))

	// Original file should be gone
	_, err = os.Stat(testFile)
	assert.True(t, os.IsNotExist(err))
}

func TestInstallTemplate(t *testing.T) {
	tmpDir := setupInstallTest(t)

	// Get a template entry from embedded files
	entries, err := templatesFS.ReadDir(defaultTemplateDir)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	var templateEntry os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() {
			templateEntry = entry
			break
		}
	}
	require.NotNil(t, templateEntry)

	tests := []struct {
		name    string
		opts    InstallOptions
		setup   func(string) error
		wantErr bool
	}{
		{
			name: "install new template",
			opts: InstallOptions{
				Force:       false,
				TemplateDir: tmpDir,
			},
			wantErr: false,
		},
		{
			name: "force install over existing",
			opts: InstallOptions{
				Force:       true,
				TemplateDir: tmpDir,
			},
			setup: func(dir string) error {
				return os.WriteFile(
					filepath.Join(dir, templateEntry.Name()),
					[]byte("existing content"),
					defaultPerms,
				)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				err := tt.setup(tt.opts.TemplateDir)
				require.NoError(t, err)
			}

			err := installTemplate(tt.opts, templateEntry)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify file was installed
			destPath := filepath.Join(tt.opts.TemplateDir, templateEntry.Name())
			assert.FileExists(t, destPath)

			if tt.opts.Force {
				// Verify backup was created
				backupPath := destPath + backupExtension
				assert.FileExists(t, backupPath)
			}
		})
	}
}

// MockPromptResponder helps test functions that require user input
type MockPromptResponder struct {
	responses []string
	current   int
}

func NewMockPromptResponder(responses []string) *MockPromptResponder {
	return &MockPromptResponder{
		responses: responses,
		current:   0,
	}
}

func (m *MockPromptResponder) GetResponse() string {
	if m.current >= len(m.responses) {
		return ""
	}
	response := m.responses[m.current]
	m.current++
	return response
}

// MockInputReader implements InputReader for testing
type MockInputReader struct {
	responses []string
	current   int
}

func NewMockInputReader(responses []string) *MockInputReader {
	return &MockInputReader{
		responses: responses,
		current:   0,
	}
}

func (r *MockInputReader) ReadResponse() (string, error) {
	if r.current >= len(r.responses) {
		return "", fmt.Errorf("no more responses")
	}
	response := r.responses[r.current]
	r.current++
	return response, nil
}

func TestPromptForAction(t *testing.T) {
	tests := []struct {
		name      string
		responses []string
		want      string
		wantErr   bool
	}{
		{
			name:      "skip all",
			responses: []string{"s"},
			want:      "skip",
			wantErr:   false,
		},
		{
			name:      "overwrite all",
			responses: []string{"o"},
			want:      "overwrite",
			wantErr:   false,
		},
		{
			name:      "choose per file",
			responses: []string{"c"},
			want:      "choose",
			wantErr:   false,
		},
		{
			name:      "invalid response",
			responses: []string{"x"},
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewMockInputReader(tt.responses)
			got, err := promptForAction(reader)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.want, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPromptForFile(t *testing.T) {
	tests := []struct {
		name      string
		file      string
		responses []string
		want      bool
		wantErr   bool
	}{
		{
			name:      "confirm overwrite",
			file:      "test.md",
			responses: []string{"y"},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "confirm overwrite uppercase",
			file:      "test.md",
			responses: []string{"Y"},
			want:      true,
			wantErr:   false,
		},
		{
			name:      "skip overwrite",
			file:      "test.md",
			responses: []string{"n"},
			want:      false,
			wantErr:   false,
		},
		{
			name:      "skip overwrite empty",
			file:      "test.md",
			responses: []string{""},
			want:      false,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := NewMockInputReader(tt.responses)
			got, err := promptForFile(reader, tt.file)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
