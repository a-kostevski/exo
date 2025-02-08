package templates

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// -------------------------
// Embedded Default Templates
// -------------------------

// DefaultTemplatesFS holds the embedded default templates.
//
//go:embed default/*
var DefaultTemplatesFS embed.FS

// DefaultTemplateBaseDir is the base directory for default templates in the embedded FS.
const (
	DefaultTemplateBaseDir = "default"
	BackupExtension        = ".bak"
	defaultDirPerms        = 0755 // Owner can read/write/execute, others can read/execute
)

// -------------------------
// Default Template Store (for Installation)
// -------------------------

// DefaultTemplateStore defines an interface to obtain built-in default templates.
type DefaultTemplateStore interface {
	// ReadTemplate returns the content of the template identified by name.
	ReadTemplate(name string) ([]byte, error)
	// ListTemplates returns a list of template filenames (with extension) in the default store.
	ListTemplates() ([]string, error)
}

type embedTemplateStore struct {
	efs     embed.FS
	baseDir string
}

// NewEmbedTemplateStore creates a new DefaultTemplateStore using an embed.FS and a base directory.
func NewEmbedTemplateStore(efs embed.FS, baseDir string) DefaultTemplateStore {
	return &embedTemplateStore{efs: efs, baseDir: baseDir}
}

func (e *embedTemplateStore) ReadTemplate(name string) ([]byte, error) {
	path := filepath.Join(e.baseDir, name)
	return e.efs.ReadFile(path)
}

func (e *embedTemplateStore) ListTemplates() ([]string, error) {
	entries, err := e.efs.ReadDir(e.baseDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

// -------------------------
// Installation API
// -------------------------

// InstallOptions holds configuration for installing default templates.
type InstallOptions struct {
	TargetDir string      // Directory where default templates will be installed.
	Force     bool        // If true, always overwrite.
	Reader    InputReader // Used to prompt user if necessary.
}

// InputReader defines an interface for reading interactive input.
type InputReader interface {
	ReadResponse() (string, error)
}

// DefaultInputReader implements InputReader using standard input.
type DefaultInputReader struct{}

func (r *DefaultInputReader) ReadResponse() (string, error) {
	var response string
	_, err := fmt.Scanln(&response)
	return response, err
}

// InstallDefaultTemplates installs built-in templates from the default template store
// into the target directory (usually the custom TemplateDir).
func InstallDefaultTemplates(cfg TemplateConfig, opts InstallOptions, defaultStore DefaultTemplateStore) error {
	if strings.TrimSpace(opts.TargetDir) == "" {
		return fmt.Errorf("target directory cannot be empty")
	}
	// Ensure target directory exists.
	if err := os.MkdirAll(opts.TargetDir, defaultDirPerms); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}
	if defaultStore == nil {
		return fmt.Errorf("default templates source is not configured")
	}
	templateFiles, err := defaultStore.ListTemplates()
	if err != nil {
		return fmt.Errorf("failed to list default templates: %w", err)
	}
	for _, file := range templateFiles {
		content, err := defaultStore.ReadTemplate(file)
		if err != nil {
			return fmt.Errorf("failed to read default template %s: %w", file, err)
		}
		destPath := filepath.Join(opts.TargetDir, file)
		// If file exists and not forced, prompt the user.
		if _, err := os.Stat(destPath); err == nil {
			if !opts.Force {
				if opts.Reader == nil {
					return fmt.Errorf("file %s exists; set Force to true to overwrite", file)
				}
				fmt.Printf("File %s exists. Overwrite? [y/n]: ", file)
				resp, err := opts.Reader.ReadResponse()
				if err != nil {
					return fmt.Errorf("failed to read user response: %w", err)
				}
				if strings.ToLower(strings.TrimSpace(resp)) != "y" {
					// Skip installation for this file.
					continue
				}
			}
			// Create backup.
			if err := CreateBackup(destPath); err != nil {
				return fmt.Errorf("failed to create backup for %s: %w", destPath, err)
			}
		}
		// Write the file.
		if err := os.WriteFile(destPath, content, cfg.FilePermissions); err != nil {
			return fmt.Errorf("failed to write template %s: %w", file, err)
		}
	}
	return nil
}

// CreateBackup renames the existing file by appending backupExtension.
// If a backup already exists, it appends a timestamp.
func CreateBackup(path string) error {
	backupPath := path + BackupExtension
	if _, err := os.Stat(backupPath); err == nil {
		timestamp := time.Now().Format("20060102150405")
		backupPath = fmt.Sprintf("%s.%s%s", path, timestamp, BackupExtension)
	}
	return os.Rename(path, backupPath)
}
