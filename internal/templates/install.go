package templates

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed default/*
var templatesFS embed.FS

const (
	defaultTemplateDir = "default"
	backupExtension    = ".bak"
	defaultPerms       = 0644 // for files
	defaultDirPerms    = 0755 // for directories
)

// InputReader defines an interface for reading user input
type InputReader interface {
	ReadResponse() (string, error)
}

// DefaultInputReader implements InputReader using standard input
type DefaultInputReader struct{}

func (r *DefaultInputReader) ReadResponse() (string, error) {
	var response string
	_, err := fmt.Scanln(&response)
	return response, err
}

// InstallOptions holds configuration for template installation
type InstallOptions struct {
	TemplateDir string
	Force       bool
	Reader      InputReader
}

// InstallDefault installs default templates to the specified directory
func InstallDefault(path string, force bool) error {
	return InstallDefaultWithOptions(InstallOptions{
		TemplateDir: path,
		Force:       force,
		Reader:      &DefaultInputReader{},
	})
}

// InstallDefaultWithOptions installs default templates with custom options
func InstallDefaultWithOptions(opts InstallOptions) error {
	if err := validateTemplatePath(opts.TemplateDir); err != nil {
		return err
	}

	// Ensure template directory exists
	if err := os.MkdirAll(opts.TemplateDir, defaultDirPerms); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	entries, err := templatesFS.ReadDir(defaultTemplateDir)
	if err != nil {
		return fmt.Errorf("failed to read default templates directory: %w", err)
	}

	templateNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			templateNames = append(templateNames, entry.Name())
		}
	}

	// First check which files already exist
	existingFiles, err := findExistingFiles(opts.TemplateDir, entries)
	if err != nil {
		return fmt.Errorf("failed to check for existing files: %w", err)
	}

	if len(existingFiles) > 0 {
		if !opts.Force {
			action, err := promptForAction(opts.Reader)
			if err != nil {
				return fmt.Errorf("failed to get user action: %w", err)
			}
			opts.Force = action == "overwrite"
		} else {
		}
	}

	return processTemplateInstallation(opts, entries)
}

func findExistingFiles(path string, entries []os.DirEntry) ([]string, error) {
	var existingFiles []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		destPath := filepath.Join(path, entry.Name())
		if _, err := os.Stat(destPath); err == nil {
			existingFiles = append(existingFiles, entry.Name())
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error checking file %s: %w", destPath, err)
		} else {
		}
	}

	return existingFiles, nil
}

func promptForAction(reader InputReader) (string, error) {
	fmt.Print("Found existing template files. Choose action \n [S]kip all\n[O]verwrite all\n[C]hoose per file\n")
	response, err := reader.ReadResponse()
	if err != nil {
		return "", fmt.Errorf("failed to read user response: %w", err)
	}

	response = strings.ToLower(strings.TrimSpace(response))
	switch response {
	case "s":
		return "skip", nil
	case "o":
		return "overwrite", nil
	case "c":
		return "choose", nil
	default:
		err := fmt.Errorf("invalid option '%s', expected one of: s, o, c", response)
		return "", err
	}
}

func processTemplateInstallation(opts InstallOptions, entries []os.DirEntry) error {
	stats := struct {
		installed int
		skipped   int
		failed    int
	}{}

	// Create a channel for errors
	errChan := make(chan error, len(entries))
	defer close(errChan)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if err := func() error {
			// Create a recovery function
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic while installing template %s: %v", entry.Name(), r)
				}
			}()

			if err := installTemplate(opts, entry); err != nil {
				stats.failed++
				return fmt.Errorf("failed to install template %s: %w", entry.Name(), err)
			}
			stats.installed++
			return nil
		}(); err != nil {
			errChan <- err
		}
	}

	// Check for any errors
	select {
	case err := <-errChan:
		return err
	default:
	}

	return nil
}

func installTemplate(opts InstallOptions, entry os.DirEntry) error {
	content, err := templatesFS.ReadFile(filepath.Join(defaultTemplateDir, entry.Name()))
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", entry.Name(), err)
	}

	destPath := filepath.Join(opts.TemplateDir, entry.Name())

	// Check if file exists
	if _, err := os.Stat(destPath); err == nil {
		if !opts.Force {
			overwrite, err := promptForFile(opts.Reader, entry.Name())
			if err != nil {
				return fmt.Errorf("failed to get user input: %w", err)
			}
			if !overwrite {
				return nil
			}
		}

		// Create backup before overwriting
		if err := createBackup(destPath); err != nil {
			return fmt.Errorf("failed to create backup for %s: %w", destPath, err)
		}
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), defaultDirPerms); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", destPath, err)
	}

	if err := os.WriteFile(destPath, content, defaultPerms); err != nil {
		return fmt.Errorf("failed to write template %s: %w", entry.Name(), err)
	}
	return nil
}

func promptForFile(reader InputReader, file string) (bool, error) {
	fmt.Printf("Overwrite %s? [y/n]: ", file)
	response, err := reader.ReadResponse()
	if err != nil {
		return false, fmt.Errorf("failed to read user response for %s: %w", file, err)
	}
	result := response == "y" || response == "Y"
	return result, nil
}

func createBackup(path string) error {
	// First, check if backup already exists
	backupPath := path + backupExtension
	if _, err := os.Stat(backupPath); err == nil {
		// Backup exists, append timestamp
		timestamp := time.Now().Format("20060102150405")
		backupPath = fmt.Sprintf("%s.%s%s", path, timestamp, backupExtension)
	}

	if err := os.Rename(path, backupPath); err != nil {
		return fmt.Errorf("failed to create backup %s: %w", backupPath, err)
	}

	return nil
}
