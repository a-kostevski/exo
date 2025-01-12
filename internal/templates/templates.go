package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateManager defines the interface for managing note templates.
// It provides functionality for loading, processing, and listing templates
// stored in the configured template directory.
type TemplateManager interface {
	// ProcessTemplate applies the named template to the given data and returns the result.
	// If the template doesn't exist, it returns an error.
	ProcessTemplate(name string, data interface{}) (string, error)

	// LoadTemplate loads the content of a template by name.
	// The template is expected to be a markdown file in the template directory.
	LoadTemplate(name string) (string, error)

	// ListTemplates returns a list of available template names without extensions.
	ListTemplates() ([]string, error)
}

// DefaultTemplateManager implements the TemplateManager interface
type DefaultTemplateManager struct {
	templateDir string
	templates   map[string]*template.Template
}

const (
	templateExtension = ".md"
	defaultPerms      = 0644
)

// NewTemplateManager creates a new template manager with the given template directory
func NewTemplateManager(templateDir string) (*DefaultTemplateManager, error) {
	if err := validateTemplatePath(templateDir); err != nil {
		return nil, err
	}

	// Create template directory if it doesn't exist
	if err := os.MkdirAll(templateDir, defaultPerms); err != nil {
		return nil, fmt.Errorf("failed to create template directory: %w", err)
	}

	// Create a new template manager
	tm := &DefaultTemplateManager{
		templateDir: templateDir,
		templates:   make(map[string]*template.Template),
	}

	// Load default templates from embedded files
	if err := tm.loadDefaultTemplates(); err != nil {
		return nil, err
	}

	// Load custom templates
	if err := tm.loadCustomTemplates(); err != nil {
		return nil, err
	}

	return tm, nil
}

// ProcessTemplate processes a template with the given data and returns the result
func (tm *DefaultTemplateManager) ProcessTemplate(name string, data interface{}) (string, error) {
	if err := validateTemplateName(name); err != nil {
		return "", err
	}

	// Try to load from file first
	if content, err := tm.LoadTemplate(name); err == nil {
		if err := tm.parseTemplate(name, content); err != nil {
			return "", err
		}
	}

	// Get template from cache or load default
	tmpl, ok := tm.templates[name]
	if !ok {
		return "", fmt.Errorf("template %q not found", name)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// LoadTemplate loads a template from the template directory
func (tm *DefaultTemplateManager) LoadTemplate(name string) (string, error) {
	if err := validateTemplateName(name); err != nil {
		return "", err
	}

	path := filepath.Join(tm.templateDir, name+templateExtension)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template %q not found", name)
		}
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	return string(content), nil
}

// ListTemplates returns a list of available templates
func (tm *DefaultTemplateManager) ListTemplates() ([]string, error) {
	templates := make(map[string]bool)

	// Add all templates from cache
	for name := range tm.templates {
		templates[name] = true
	}

	// Read template directory
	entries, err := os.ReadDir(tm.templateDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read template directory: %w", err)
	}

	// Add templates from directory
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && filepath.Ext(entry.Name()) == templateExtension {
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

	return result, nil
}

// Helper functions

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
	entries, err := os.ReadDir(tm.templateDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read template directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == templateExtension {
			name := getTemplateName(entry.Name())
			content, err := os.ReadFile(filepath.Join(tm.templateDir, entry.Name()))
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
	return filename[:len(filename)-len(templateExtension)]
}
