package templates

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

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
		opts        []TemplateManagerOption
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid directory",
			templateDir: setupTest(t),
			opts:        nil,
			wantErr:     false,
		},
		{
			name:        "with custom extension",
			templateDir: setupTest(t),
			opts: []TemplateManagerOption{
				WithTemplateExtension(".txt"),
			},
			wantErr: false,
		},
		{
			name:        "empty directory",
			templateDir: "",
			opts:        nil,
			wantErr:     true,
			errMsg:      "template directory cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := NewTemplateManager(tt.templateDir, tt.opts...)
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
			assert.Equal(t, tt.templateDir, tm.config.TemplateDir)
		})
	}
}

func TestTemplateManager_ProcessTemplate(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create a custom template
	customTemplate := "# Custom Note {{.Title}}\n\nContent: {{.Content}}"
	err = os.WriteFile(filepath.Join(tmpDir, "custom.md"), []byte(customTemplate), tm.config.FilePermissions)
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
			name:     "with context timeout",
			template: "custom",
			data:     struct{}{},
			wantErr:  true,
			errMsg:   "context deadline exceeded",
		},
		{
			name:     "non-existent template",
			template: "nonexistent",
			data:     struct{}{},
			wantErr:  true,
			errMsg:   "template \"nonexistent\" not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.name == "with context timeout" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Millisecond)
				defer cancel()
				time.Sleep(time.Millisecond * 2)
			}

			got, err := tm.ProcessTemplateWithContext(ctx, tt.template, tt.data)
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

func TestTemplateManager_LoadTemplate(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create a test template
	testTemplate := "# Test Template"
	err = os.WriteFile(filepath.Join(tmpDir, "test.md"), []byte(testTemplate), tm.config.FilePermissions)
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

func TestTemplateManager_ListTemplates(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create test templates
	templates := map[string]string{
		"test1": "# Test Template 1",
		"test2": "# Test Template 2",
	}
	for name, content := range templates {
		err = os.WriteFile(
			filepath.Join(tmpDir, name+tm.config.TemplateExtension),
			[]byte(content),
			tm.config.FilePermissions,
		)
		require.NoError(t, err)
	}

	t.Run("list all templates", func(t *testing.T) {
		got, err := tm.ListTemplates()
		assert.NoError(t, err)
		for name := range templates {
			assert.Contains(t, got, name)
		}
	})
}

func TestTemplateOverrides(t *testing.T) {
	tmpDir := setupTest(t)
	tm, err := NewTemplateManager(tmpDir)
	require.NoError(t, err)

	// Create a custom template that overrides a default one
	customTemplate := "# Custom Template {{.Title}}"
	err = os.WriteFile(filepath.Join(tmpDir, "note.md"), []byte(customTemplate), tm.config.FilePermissions)
	require.NoError(t, err)

	data := struct {
		Title string
	}{
		Title: "Test Override",
	}

	got, err := tm.ProcessTemplate("note", data)
	assert.NoError(t, err)
	assert.Equal(t, "# Custom Template Test Override", got)
}
