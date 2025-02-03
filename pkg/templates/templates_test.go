package templates

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestConfig(t *testing.T) TemplateConfig {
	return TemplateConfig{
		TemplateDir:       t.TempDir(),
		TemplateExtension: ".md",
		Logger:            logger.NewMockLogger(t),
		FilePermissions:   0644,
	}
}

func TestNewTemplateManager(t *testing.T) {
	tests := []struct {
		name    string
		config  TemplateConfig
		opts    []TemplateManagerOption
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid config",
			config:  setupTestConfig(t),
			wantErr: false,
		},
		{
			name: "with custom extension",
			config: func() TemplateConfig {
				cfg := setupTestConfig(t)
				cfg.TemplateExtension = ".txt"
				return cfg
			}(),
			wantErr: false,
		},
		{
			name: "missing logger",
			config: TemplateConfig{
				TemplateDir: t.TempDir(),
			},
			wantErr: true,
			errMsg:  "logger is required",
		},
		{
			name: "empty directory",
			config: TemplateConfig{
				Logger: logger.NewMockLogger(t),
				// TemplateDir is intentionally empty
			},
			wantErr: true,
			errMsg:  "template directory is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm, err := NewTemplateManager(tt.config, tt.opts...)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tm)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}
			require.NoError(t, err)
			require.NotNil(t, tm)
			assert.Equal(t, tt.config.TemplateDir, tm.config.TemplateDir)
		})
	}
}

func TestTemplateManager_ProcessTemplate(t *testing.T) {
	config := setupTestConfig(t)
	tm, err := NewTemplateManager(config)
	require.NoError(t, err)

	// Create a custom template
	customTemplate := "# Custom Note {{.Title}}\n\nContent: {{.Content}}"
	err = os.WriteFile(filepath.Join(config.TemplateDir, "custom.md"), []byte(customTemplate), config.FilePermissions)
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
	config := setupTestConfig(t)
	tm, err := NewTemplateManager(config)
	require.NoError(t, err)

	testTemplate := "# Test Template"
	err = os.WriteFile(filepath.Join(config.TemplateDir, "test.md"), []byte(testTemplate), config.FilePermissions)
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
	config := setupTestConfig(t)
	tm, err := NewTemplateManager(config)
	require.NoError(t, err)

	templates := map[string]string{
		"test1": "# Test Template 1",
		"test2": "# Test Template 2",
	}
	for name, content := range templates {
		err = os.WriteFile(
			filepath.Join(config.TemplateDir, name+config.TemplateExtension),
			[]byte(content),
			config.FilePermissions,
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
	config := setupTestConfig(t)
	tm, err := NewTemplateManager(config)
	require.NoError(t, err)

	customTemplate := "# Custom Template {{.Title}}"
	err = os.WriteFile(filepath.Join(config.TemplateDir, "note.md"), []byte(customTemplate), config.FilePermissions)
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
