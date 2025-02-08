package testutil

import (
	"os"
	"path/filepath"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/a-kostevski/exo/pkg/fs"
	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/templates"
)

// DummyInputReader implements InputReader for testing; it always returns a fixed response.
type DummyInputReader struct {
	Response string
}

func (d *DummyInputReader) ReadResponse() (string, error) {
	return d.Response, nil
}

// DummyTemplateManager is a simple dummy implementation for testing purposes.
type DummyTemplateManager struct{}

// ProcessTemplate returns a fixed string if a "Title" field is provided.
func (dtm *DummyTemplateManager) ProcessTemplate(name string, data interface{}) (string, error) {
	if m, ok := data.(map[string]interface{}); ok {
		if title, ok := m["Title"].(string); ok {
			return "Template: " + title, nil
		}
	}
	return "Template: unknown", nil
}

func (dtm *DummyTemplateManager) ProcessTemplateWithContext(ctx interface{}, name string, data interface{}) (string, error) {
	return dtm.ProcessTemplate(name, data)
}

func (dtm *DummyTemplateManager) LoadTemplate(name string) (string, error) {
	return "", nil
}

func (dtm *DummyTemplateManager) ListTemplates() ([]string, error) {
	return []string{}, nil
}

// InstallDefaultTemplates implements the required method from TemplateManager interface
func (dtm *DummyTemplateManager) InstallDefaultTemplates(opts templates.InstallOptions) error {
	return nil // For testing purposes, just return success
}

// DummyLogger is a no-op logger for testing purposes.
type DummyLogger struct{}

// NewDummyLogger creates and returns a new DummyLogger.
func NewDummyLogger() logger.Logger {
	return &DummyLogger{}
}

// Info does nothing.
func (dl *DummyLogger) Info(msg string, fields ...logger.Field) {}

// Error does nothing.
func (dl *DummyLogger) Error(msg string, fields ...logger.Field) {}

// Infof does nothing.
func (dl *DummyLogger) Infof(format string, args ...interface{}) {}

// Errorf does nothing.
func (dl *DummyLogger) Errorf(format string, args ...interface{}) {}

// DummyFS is a dummy implementation of fs.FileSystem that uses basic OS calls
// but can be defined here to avoid importing production OSFileSystem.
type DummyFS struct{}

func (d *DummyFS) EnsureDirectoryExists(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

func (d *DummyFS) WriteFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}

func (d *DummyFS) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (d *DummyFS) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (d *DummyFS) DeleteFile(path string) error {
	return os.Remove(path)
}

func (d *DummyFS) AppendToFile(path, content string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content + "\n")
	return err
}

func (d *DummyFS) OpenInEditor(path, editor string) error {
	// For testing, simply simulate success.
	return nil
}

func (d *DummyFS) ReadDir(path string) ([]os.DirEntry, error) {
	// Use the OS-based implementation for simplicity.
	return os.ReadDir(path)
}

// NewDummyFS returns an instance of DummyFS.
func NewDummyFS() fs.FileSystem {
	return &DummyFS{}
}

// NewDummyDeps returns dummy dependencies for testing.
func NewDummyDeps(dataHome string) (config.Config, templates.TemplateManager, logger.Logger, fs.FileSystem, func()) {
	cfg := config.Config{
		General: config.GeneralConfig{
			Editor: "dummy-editor",
		},
		Dir: config.DirConfig{
			DataHome:    dataHome,
			TemplateDir: filepath.Join(dataHome, "templates"),
			PeriodicDir: filepath.Join(dataHome, "periodic"),
			ZettelDir:   filepath.Join(dataHome, "zettel"),
			ProjectsDir: filepath.Join(dataHome, "projects"),
			InboxDir:    filepath.Join(dataHome, "0-inbox"),
			IdeaDir:     filepath.Join(dataHome, "ideas"),
		},
	}
	_ = os.MkdirAll(cfg.Dir.DataHome, 0755)
	dtm := &DummyTemplateManager{}
	dl := NewDummyLogger()
	dfs := NewDummyFS()

	cleanup := func() {
		// No cleanup needed; t.TempDir() handles removal.
	}
	return cfg, dtm, dl, dfs, cleanup
}

// NewDummyDepsTemplates returns a TemplateConfig for testing.
func NewDummyDepsTemplates(dataHome string) (templates.TemplateConfig, func()) {
	cfg := templates.TemplateConfig{
		TemplateDir:       filepath.Join(dataHome, "templates"),
		TemplateExtension: ".md",
		FilePermissions:   0644,
		Logger:            &DummyLogger{},
	}
	_ = os.MkdirAll(cfg.TemplateDir, 0755)
	cleanup := func() {
		// Cleanup is handled by t.TempDir() in tests.
	}
	return cfg, cleanup
}
