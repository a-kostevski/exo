package templates

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/a-kostevski/exo/pkg/logger"
)

type TemplateManager interface {
	ProcessTemplate(name string, data interface{}) (string, error)
	ProcessTemplateWithContext(ctx context.Context, name string, data interface{}) (string, error)

	LoadTemplate(name string) (string, error)
	ListTemplates() ([]string, error)
}

type TemplateConfig struct {
	TemplateDir       string
	TemplateExtension string
	Logger            logger.Logger
	FilePermissions   os.FileMode
	FS                fs.FS
}

// DefaultTemplateManager implements the TemplateManager interface
type DefaultTemplateManager struct {
	config    TemplateConfig
	templates map[string]*template.Template
	log       logger.Logger
}

type TemplateManagerOption func(*DefaultTemplateManager) error

// NewTemplateManager creates a new template manager with the given template directory
func NewTemplateManager(config TemplateConfig, opts ...TemplateManagerOption) (*DefaultTemplateManager, error) {
	if config.Logger == nil {
		return nil, errors.New("logger is required")
	}

	if config.TemplateDir == "" {
		return nil, errors.New("template directory is required")
	}

	if config.FilePermissions == 0 {
		config.FilePermissions = 0644
	}
	if config.TemplateExtension == "" {
		config.TemplateExtension = ".md"
	}

	tm := &DefaultTemplateManager{
		config:    config,
		templates: make(map[string]*template.Template),
		log:       config.Logger,
	}

	for _, opt := range opts {
		if err := opt(tm); err != nil {
			tm.log.Error("Failed to apply template manager option",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	if err := tm.initialize(); err != nil {
		tm.log.Error("Failed to initialize template manager",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})
		return nil, err
	}

	tm.log.Info("Successfully created template manager",
		logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})
	return tm, nil
}

func (tm *DefaultTemplateManager) initialize() error {
	tm.log.Debug("Initializing template manager",
		logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})

	if err := os.MkdirAll(tm.config.TemplateDir, tm.config.FilePermissions); err != nil {
		tm.log.Error("Failed to create template directory",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	if err := tm.loadTemplates(); err != nil {
		tm.log.Error("Failed to load templates",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})
		return err
	}

	tm.log.Info("Successfully initialized template manager",
		logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})
	return nil
}

func (tm *DefaultTemplateManager) loadTemplates() error {
	return nil
}

// ProcessTemplate processes a template with the given data and returns the result
func (tm *DefaultTemplateManager) ProcessTemplate(name string, data interface{}) (string, error) {
	tm.log.Debug("Processing template",
		logger.Field{Key: "name", Value: name})

	if err := validateTemplateName(name); err != nil {
		tm.log.Error("Invalid template name",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name})
		return "", err
	}

	// Try to load from file first
	if content, err := tm.LoadTemplate(name); err == nil {
		tm.log.Debug("Parsing template from file",
			logger.Field{Key: "name", Value: name})
		if err := tm.parseTemplate(name, content); err != nil {
			tm.log.Error("Failed to parse template",
				logger.Field{Key: "error", Value: err},
				logger.Field{Key: "name", Value: name})
			return "", err
		}
	}

	// Get template from cache or load default
	tmpl, ok := tm.templates[name]
	if !ok {
		tm.log.Error("Template not found",
			logger.Field{Key: "name", Value: name})
		return "", fmt.Errorf("template %q not found", name)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		tm.log.Error("Failed to execute template",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name})
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	tm.log.Info("Successfully processed template",
		logger.Field{Key: "name", Value: name})
	return buf.String(), nil
}

func (tm *DefaultTemplateManager) ProcessTemplateWithContext(ctx context.Context, name string, data interface{}) (string, error) {
	tm.log.Debug("Processing template with context",
		logger.Field{Key: "name", Value: name})

	// Create a channel for the result
	resultCh := make(chan struct {
		result string
		err    error
	})

	// Process template in a goroutine
	go func() {
		result, err := tm.ProcessTemplate(name, data)
		resultCh <- struct {
			result string
			err    error
		}{result, err}
	}()

	// Wait for either context cancellation or template processing completion
	select {
	case <-ctx.Done():
		tm.log.Warn("Template processing cancelled by context",
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "error", Value: ctx.Err()})
		return "", ctx.Err()
	case res := <-resultCh:
		if res.err != nil {
			tm.log.Error("Failed to process template with context",
				logger.Field{Key: "error", Value: res.err},
				logger.Field{Key: "name", Value: name})
		} else {
			tm.log.Info("Successfully processed template with context",
				logger.Field{Key: "name", Value: name})
		}
		return res.result, res.err
	}
}

// LoadTemplate loads a template from the template directory
func (tm *DefaultTemplateManager) LoadTemplate(name string) (string, error) {
	tm.log.Debug("Loading template",
		logger.Field{Key: "name", Value: name})

	if err := validateTemplateName(name); err != nil {
		tm.log.Error("Invalid template name",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "name", Value: name})
		return "", err
	}

	path := filepath.Join(tm.config.TemplateDir, name+tm.config.TemplateExtension)
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			tm.log.Error("Template file not found",
				logger.Field{Key: "path", Value: path},
				logger.Field{Key: "name", Value: name})
			return "", fmt.Errorf("template %q not found", name)
		}
		tm.log.Error("Failed to read template file",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "path", Value: path})
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	tm.log.Info("Successfully loaded template",
		logger.Field{Key: "name", Value: name},
		logger.Field{Key: "path", Value: path})
	return string(content), nil
}

func (tm *DefaultTemplateManager) ListTemplates() ([]string, error) {
	tm.log.Debug("Listing templates")

	templates := make(map[string]bool)

	// Add all templates from cache
	for name := range tm.templates {
		templates[name] = true
	}

	// Read template directory
	entries, err := os.ReadDir(tm.config.TemplateDir)
	if err != nil && !os.IsNotExist(err) {
		tm.log.Error("Failed to read template directory",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "tm.config.TemplateDir", Value: tm.config.TemplateDir})
		return nil, fmt.Errorf("failed to read template directory: %w", err)
	}

	// Add templates from directory
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == tm.config.TemplateExtension {
				name := getTemplateName(entry.Name())
				templates[name] = true
			}
		}
	}

	// Convert map to slice
	var result []string
	for name := range templates {
		result = append(result, name)
	}

	tm.log.Info("Successfully listed templates",
		logger.Field{Key: "count", Value: len(result)})
	return result, nil
}

func WithTemplateExtension(ext string) TemplateManagerOption {
	return func(tm *DefaultTemplateManager) error {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		tm.config.TemplateExtension = ext
		return nil
	}
}

func validateTemplatePath(path string) error {
	if path == "" {
		return fmt.Errorf("template directory cannot be empty")
	}
	return nil
}

func validateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}
	return nil
}

func (tm *DefaultTemplateManager) parseTemplate(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	tm.templates[name] = tmpl
	return nil
}

func (tm *DefaultTemplateManager) loadDefaultTemplates() error {
	entries, err := templatesFS.ReadDir("default")
	if err != nil {
		return fmt.Errorf("failed to read default templates: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := getTemplateName(entry.Name())
		content, err := templatesFS.ReadFile(filepath.Join("default", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read default template %q: %w", name, err)
		}

		if err := tm.parseTemplate(name, string(content)); err != nil {
			return err
		}
	}

	return nil
}

func (tm *DefaultTemplateManager) loadCustomTemplates() error {
	entries, err := os.ReadDir(tm.config.TemplateDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read template directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == tm.config.TemplateExtension {
			name := getTemplateName(entry.Name())
			content, err := os.ReadFile(filepath.Join(tm.config.TemplateDir, entry.Name()))
			if err != nil {
				return fmt.Errorf("failed to read template %q: %w", name, err)
			}

			if err := tm.parseTemplate(name, string(content)); err != nil {
				return err
			}
		}
	}

	return nil
}

func getTemplateName(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}
