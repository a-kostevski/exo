package templates_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/templates"
	"github.com/a-kostevski/exo/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallDefaultTemplates_Forced(t *testing.T) {
	// Use a temporary directory as the target for installation.
	tmpDir := t.TempDir()

	// create a default template store using the embedded defaults.
	defaultStore := templates.NewEmbedTemplateStore(templates.DefaultTemplatesFS, templates.DefaultTemplateBaseDir)

	cfg := templates.TemplateConfig{
		TemplateDir:       tmpDir, // Target directory.
		TemplateExtension: ".md",
		FilePermissions:   0644,
		Logger:            testutil.NewDummyLogger(),
		FS:                testutil.NewDummyFS(),
	}
	opts := templates.InstallOptions{
		TargetDir: tmpDir,
		Force:     true, // Force installation (no interactive prompt).
		Reader:    &testutil.DummyInputReader{Response: "overwrite"},
	}
	err := templates.InstallDefaultTemplates(cfg, opts, defaultStore)
	require.NoError(t, err)

	// Verify that each file from the default store is installed.
	defFiles, err := defaultStore.ListTemplates()
	require.NoError(t, err)
	for _, file := range defFiles {
		destPath := filepath.Join(tmpDir, file)
		_, err := os.Stat(destPath)
		assert.NoError(t, err, "expected template %s to be installed", file)

		expectedContent, err := defaultStore.ReadTemplate(file)
		require.NoError(t, err)
		actualContent, err := os.ReadFile(destPath)
		require.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(string(expectedContent)), strings.TrimSpace(string(actualContent)))
	}
}

func TestInstallDefaultTemplates_Interactive(t *testing.T) {
	// Use a temporary directory as target.
	tmpDir := t.TempDir()
	defaultStore := templates.NewEmbedTemplateStore(templates.DefaultTemplatesFS, templates.DefaultTemplateBaseDir)

	cfg := templates.TemplateConfig{
		TemplateDir:       tmpDir,
		TemplateExtension: ".md",
		FilePermissions:   0644,
		Logger:            testutil.NewDummyLogger(),
		FS:                fs.NewOSFileSystem(),
	}

	// First, install templates forcefully.
	optsForced := templates.InstallOptions{
		TargetDir: tmpDir,
		Force:     true,
		Reader:    &testutil.DummyInputReader{Response: "overwrite"},
	}
	err := templates.InstallDefaultTemplates(cfg, optsForced, defaultStore)
	require.NoError(t, err)

	// Now simulate interactive installation by setting Force=false and responding "n" (skip).
	optsInteractive := templates.InstallOptions{
		TargetDir: tmpDir,
		Force:     false,
		Reader:    &testutil.DummyInputReader{Response: "n"},
	}
	err = templates.InstallDefaultTemplates(cfg, optsInteractive, defaultStore)
	require.NoError(t, err)

	// Verify that installed files remain unchanged.
	defFiles, err := defaultStore.ListTemplates()
	require.NoError(t, err)
	for _, file := range defFiles {
		destPath := filepath.Join(tmpDir, file)
		_, err := os.Stat(destPath)
		assert.NoError(t, err, "expected template %s to exist", file)
	}
}

func TestcreateBackup(t *testing.T) {
	tmpDir := t.TempDir()
	originalPath := filepath.Join(tmpDir, "sample.md")
	content := []byte("original content")
	err := os.WriteFile(originalPath, content, 0644)
	require.NoError(t, err)

	err = templates.CreateBackup(originalPath)
	require.NoError(t, err)

	// The original file should no longer exist.
	_, err = os.Stat(originalPath)
	assert.True(t, os.IsNotExist(err))

	backupPath := originalPath + templates.BackupExtension
	_, err = os.Stat(backupPath)
	require.NoError(t, err)

	backupContent, err := os.ReadFile(backupPath)
	require.NoError(t, err)
	assert.Equal(t, content, backupContent)
}

func TestcreateBackup_UniqueNames(t *testing.T) {
	tmpDir := t.TempDir()
	originalPath := filepath.Join(tmpDir, "sample.md")
	content := []byte("original")
	err := os.WriteFile(originalPath, content, 0644)
	require.NoError(t, err)

	err = templates.CreateBackup(originalPath)
	require.NoError(t, err)

	backupPath1 := originalPath + templates.BackupExtension
	_, err = os.Stat(backupPath1)
	require.NoError(t, err)

	err = os.WriteFile(originalPath, content, 0644)
	require.NoError(t, err)
	err = templates.CreateBackup(originalPath)
	require.NoError(t, err)

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	var backups []string
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), "sample.md") && strings.Contains(e.Name(), templates.BackupExtension) {
			backups = append(backups, e.Name())
		}
	}
	assert.GreaterOrEqual(t, len(backups), 2)
}
