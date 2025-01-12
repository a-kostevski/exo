package templates

import (
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) string {
	return t.TempDir()
}

func TestNewTemplateManager(t *testing.T) {
	tests := []struct {
		name        string
		templateDir string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid directory",
			templateDir: setupTest(t),
			wantErr:     false,
		},
		{
			name:        "empty directory",
			templateDir: "",
			wantErr:     true,
			errMsg:      "template directory cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := NewTemplateManager(tt.templateDir)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tm)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, tm)
			assert.Equal(t, tt.templateDir, tm.templateDir)

			// Verify default templates are loaded
			templates, err := tm.ListTemplates()
			assert.NoError(t, err)
			assert.NotEmpty(t, templates)
		})
	}
}

func TestDefaultTemplateManager_ProcessTemplate(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create a custom template
	customTemplate := "# Custom Note {{.Title}}\n\nContent: {{.Content}}"
	err = os.WriteFile(filepath.Join(tmpDir, "custom.md"), []byte(customTemplate), defaultPerms)
	require.NoError(t, err)

	tests := []struct {
		name     string
		template string
		data     interface{}
		want     string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "process custom template",
			template: "custom",
			data: struct {
				Title   string
				Content string
			}{
				Title:   "Test Custom",
				Content: "Test Content",
			},
			want:    "# Custom Note Test Custom\n\nContent: Test Content",
			wantErr: false,
		},
		{
			name:     "non-existent template",
			template: "nonexistent",
			data:     struct{}{},
			want:     "",
			wantErr:  true,
			errMsg:   "template \"nonexistent\" not found",
		},
		{
			name:     "empty template name",
			template: "",
			data:     struct{}{},
			want:     "",
			wantErr:  true,
			errMsg:   "template name cannot be empty",
		},
		{
			name:     "invalid template data",
			template: "custom",
			data:     "invalid",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tm.ProcessTemplate(tt.template, tt.data)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultTemplateManager_LoadTemplate(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create a test template
	testTemplate := "# Test Template"
	err = os.WriteFile(filepath.Join(tmpDir, "test.md"), []byte(testTemplate), defaultPerms)
	require.NoError(t, err)

	tests := []struct {
		name     string
		template string
		want     string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "load existing template",
			template: "test",
			want:     testTemplate,
			wantErr:  false,
		},
		{
			name:     "load non-existent template",
			template: "nonexistent",
			want:     "",
			wantErr:  true,
			errMsg:   "template \"nonexistent\" not found",
		},
		{
			name:     "empty template name",
			template: "",
			want:     "",
			wantErr:  true,
			errMsg:   "template name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tm.LoadTemplate(tt.template)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultTemplateManager_ListTemplates(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create some test templates
	templates := map[string]string{
		"test1": "# Test Template 1",
		"test2": "# Test Template 2",
	}
	for name, content := range templates {
		err = os.WriteFile(filepath.Join(tmpDir, name+templateExtension), []byte(content), defaultPerms)
		require.NoError(t, err)
	}

	// Test listing templates
	t.Run("list all templates", func(t *testing.T) {
		got, err := tm.ListTemplates()
		assert.NoError(t, err)

		// Check that custom templates are included
		for name := range templates {
			assert.Contains(t, got, name)
		}

		// Check that we have more templates than just our custom ones
		// (default templates should be included)
		assert.Greater(t, len(got), len(templates))
	})

	// Test with non-existent directory
	t.Run("non-existent directory", func(t *testing.T) {
		tm, err := NewTemplateManager(filepath.Join(tmpDir, "nonexistent"))
		require.NoError(t, err)

		got, err := tm.ListTemplates()
		assert.NoError(t, err)
		assert.NotEmpty(t, got) // Should still have default templates
	})
}

func TestTemplateOverrides(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create a custom template that overrides a default one
	customTemplate := "# Custom Template {{.Title}}"
	err = os.WriteFile(filepath.Join(tmpDir, "note.md"), []byte(customTemplate), defaultPerms)
	require.NoError(t, err)

	data := struct {
		Title string
	}{
		Title: "Test Override",
	}

	// Process the template and verify it uses the custom version
	got, err := tm.ProcessTemplate("note", data)
	assert.NoError(t, err)
	assert.Equal(t, "# Custom Template Test Override", got)
}

// Test helper functions

func TestValidateTemplatePath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid path",
			path:    "/valid/path",
			wantErr: false,
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
			errMsg:  "template directory cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplatePath(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name     string
		template string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid name",
			template: "valid",
			wantErr:  false,
		},
		{
			name:     "empty name",
			template: "",
			wantErr:  true,
			errMsg:   "template name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplateName(tt.template)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestGetTemplateName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "standard template",
			filename: "test.md",
			want:     "test",
		},
		{
			name:     "template with dots",
			filename: "test.note.md",
			want:     "test.note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTemplateName(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseTemplate(t *testing.T) {
	tm := &DefaultTemplateManager{
		templates: make(map[string]*template.Template),
	}

	tests := []struct {
		name     string
		template string
		content  string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "valid template",
			template: "test",
			content:  "# {{.Title}}",
			wantErr:  false,
		},
		{
			name:     "invalid template",
			template: "invalid",
			content:  "# {{.Title}",
			wantErr:  true,
			errMsg:   "failed to parse template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tm.parseTemplate(tt.template, tt.content)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.Contains(t, tm.templates, tt.template)
		})
	}
}
