package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
)

// TemplateManager defines the interface for processing templates.
type TemplateManager interface {
	// ProcessTemplate loads a template from the custom directory, parses it,
	// executes it with the given data, and returns the resulting string.
	ProcessTemplate(name string, data interface{}) (string, error)
	// ListTemplates returns the names (without extension) of templates available in the custom directory.
	ListTemplates() ([]string, error)
}

// TemplateConfig holds configuration for template processing.
type TemplateConfig struct {
	TemplateDir       string        // Custom directory from which to load templates.
	TemplateExtension string        // e.g. ".md"
	FilePermissions   os.FileMode   // For writing files.
	Logger            logger.Logger // Logger to use.
	FS                fs.FileSystem // Abstract file system for file operations.
}

// defaultTemplateManager implements TemplateManager.
type defaultTemplateManager struct {
	config TemplateConfig
}

// NewTemplateManager creates a new TemplateManager instance using dependency injection.
func NewTemplateManager(cfg TemplateConfig) (TemplateManager, error) {
	if strings.TrimSpace(cfg.TemplateDir) == "" {
		return nil, fmt.Errorf("template directory is required")
	}
	if strings.TrimSpace(cfg.TemplateExtension) == "" {
		cfg.TemplateExtension = ".md"
	}
	if cfg.FilePermissions == 0 {
		cfg.FilePermissions = 0644
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	if cfg.FS == nil {
		return nil, fmt.Errorf("file system is required")
	}
	return &defaultTemplateManager{config: cfg}, nil
}

// ProcessTemplate loads and executes a template from the custom directory.
func (tm *defaultTemplateManager) ProcessTemplate(name string, data interface{}) (string, error) {
	path := filepath.Join(tm.config.TemplateDir, name+tm.config.TemplateExtension)
	content, err := tm.config.FS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %w", name, err)
	}
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		tm.config.Logger.Error("failed to parse template",
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "error", Value: err})
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		tm.config.Logger.Error("failed to execute template",
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "error", Value: err})
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// ListTemplates lists the names (without extension) of templates in the custom directory.
func (tm *defaultTemplateManager) ListTemplates() ([]string, error) {
	entries, err := tm.config.FS.ReadDir(tm.config.TemplateDir)
	if err != nil {
		tm.config.Logger.Error("failed to read template directory",
			logger.Field{Key: "dir", Value: tm.config.TemplateDir},
			logger.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to read template directory: %w", err)
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == tm.config.TemplateExtension {
			name := entry.Name()[0 : len(entry.Name())-len(tm.config.TemplateExtension)]
			names = append(names, name)
		}
	}
	return names, nil
}
