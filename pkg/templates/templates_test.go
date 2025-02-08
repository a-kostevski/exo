package templates_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/templates"
	"github.com/a-kostevski/exo/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	// Write a sample template file.
	templateName := "greeting"
	templateContent := "Hello, {{.Name}}!"
	templateFile := filepath.Join(tmpDir, templateName+".md")
	err := os.WriteFile(templateFile, []byte(templateContent), 0644)
	require.NoError(t, err)

	cfg := templates.TemplateConfig{
		TemplateDir:       tmpDir,
		TemplateExtension: ".md",
		FilePermissions:   0644,
		Logger:            testutil.NewDummyLogger(),
		FS:                testutil.NewDummyFS(),
	}
	tm, err := templates.NewTemplateManager(cfg)
	require.NoError(t, err)

	data := map[string]interface{}{"Name": "Alice"}
	result, err := tm.ProcessTemplate(templateName, data)
	require.NoError(t, err)
	assert.Equal(t, "Hello, Alice!", result)
}

func TestListTemplates(t *testing.T) {
	tmpDir := t.TempDir()
	// Create two template files.
	files := []struct {
		filename string
		content  string
	}{
		{"first.md", "Content 1"},
		{"second.md", "Content 2"},
	}
	for _, f := range files {
		path := filepath.Join(tmpDir, f.filename)
		err := os.WriteFile(path, []byte(f.content), 0644)
		require.NoError(t, err)
	}
	cfg := templates.TemplateConfig{
		TemplateDir:       tmpDir,
		TemplateExtension: ".md",
		FilePermissions:   0644,
		Logger:            testutil.NewDummyLogger(),
		FS:                fs.NewOSFileSystem(),
	}
	tm, err := templates.NewTemplateManager(cfg)
	require.NoError(t, err)

	names, err := tm.ListTemplates()
	require.NoError(t, err)
	assert.Contains(t, names, "first")
	assert.Contains(t, names, "second")
	assert.Equal(t, 2, len(names))
}
