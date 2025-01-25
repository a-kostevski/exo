package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/a-kostevski/exo/internal/logger"
	"github.com/a-kostevski/exo/internal/utils"
	"github.com/spf13/viper"
)

// Environment variables
const (
	envDataHome      = "EXO_DATA_HOME"
	envXDGData       = "XDG_DATA_HOME"
	envXDGCache      = "XDG_CACHE_HOME"
	defaultXDGData   = ".local/share"
	defaultXDGCache  = ".cache"
	defaultXDGConfig = ".config"
)

// Default configuration values
const (
	defaultEditor    = "nvim"
	defaultLogLevel  = string(logger.InfoLevel)
	defaultLogFormat = string(logger.TextFormat)
	defaultLogOutput = "stdout"
)

// Config represents the main configuration structure
type Config struct {
	General GeneralConfig `mapstructure:"general"`
	Dir     DirConfig     `mapstructure:"dir"`
	Log     logger.Config `mapstructure:"log"`
}

type GeneralConfig struct {
	Editor string `mapstructure:"editor"`
}

// DirConfig holds directory-related configuration
type DirConfig struct {
	DataHome    string `mapstructure:"data_home"`
	TemplateDir string `mapstructure:"template_dir"`
	PeriodicDir string `mapstructure:"periodic_dir"`
	ZettelDir   string `mapstructure:"zettel_dir"`
	ProjectsDir string `mapstructure:"projects_dir"`
	InboxDir    string `mapstructure:"inbox_dir"`
	IdeaDir     string `mapstructure:"idea_dir"`
}

var (
	cfg  *Config
	once sync.Once
)

// Get returns the singleton config instance
func Get() (*Config, error) {
	if cfg == nil {
		return nil, errors.New("configuration not initialized")
	}
	return cfg, nil
}

func MustGet() *Config {
	if cfg == nil {
		panic("configuration not initialized")
	}
	return cfg
}

// Initialize sets up the configuration once
func Initialize(configPath string) error {
	var initErr error
	once.Do(func() {
		var loadedCfg *Config
		loadedCfg, initErr = load(configPath)
		if initErr == nil {
			initErr = loadedCfg.Validate()
			if initErr == nil {
				cfg = loadedCfg
			}
		}
	})
	return initErr
}

// Returns an error if configuration is not initialized.
func load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	v.AddConfigPath(filepath.Join(home, defaultXDGConfig))

	// Set default values
	v.SetDefault("general.editor", defaultEditor)
	v.SetDefault("log.level", defaultLogLevel)
	v.SetDefault("log.format", defaultLogFormat)
	v.SetDefault("log.output", defaultLogOutput)

	dataHome := getDataHome(home)
	v.SetDefault("dir.data_home", dataHome)
	v.SetDefault("dir.template_dir", filepath.Join(dataHome, "templates"))
	v.SetDefault("dir.periodic_dir", filepath.Join(dataHome, "periodic"))
	v.SetDefault("dir.zettel_dir", filepath.Join(dataHome, "zettel"))
	v.SetDefault("dir.projects_dir", filepath.Join(dataHome, "projects"))
	v.SetDefault("dir.inbox_dir", filepath.Join(dataHome, "0-inbox"))
	v.SetDefault("dir.idea_dir", filepath.Join(dataHome, "ideas"))

	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand paths
	cfg.Dir.DataHome = sanitizePath(cfg.Dir.DataHome, home)
	cfg.Dir.TemplateDir = sanitizePath(cfg.Dir.TemplateDir, home)
	cfg.Dir.PeriodicDir = sanitizePath(cfg.Dir.PeriodicDir, home)
	cfg.Dir.ZettelDir = sanitizePath(cfg.Dir.ZettelDir, home)
	cfg.Dir.ProjectsDir = sanitizePath(cfg.Dir.ProjectsDir, home)
	cfg.Dir.InboxDir = sanitizePath(cfg.Dir.InboxDir, home)
	cfg.Dir.IdeaDir = sanitizePath(cfg.Dir.IdeaDir, home)

	// Apply environment variable overrides
	if editor := os.Getenv("EDITOR"); editor != "" {
		cfg.General.Editor = editor
	}

	return &cfg, nil
}

func defaultConfig() map[string]interface{} {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	dataHome := getDataHome(home)
	return map[string]interface{}{
		"general.editor":   defaultEditor,
		"log.level":        defaultLogLevel,
		"log.format":       defaultLogFormat,
		"log.output":       defaultLogOutput,
		"dir.data_home":    dataHome,
		"dir.template_dir": filepath.Join(dataHome, "templates"),
		"dir.periodic_dir": filepath.Join(dataHome, "periodic"),
		"dir.zettel_dir":   filepath.Join(dataHome, "zettel"),
		"dir.projects_dir": filepath.Join(dataHome, "projects"),
		"dir.inbox_dir":    filepath.Join(dataHome, "0-inbox"),
		"dir.idea_dir":     filepath.Join(dataHome, "ideas"),
	}
}

// getDataHome determines the appropriate data home directory based on environment
// variables and system defaults. It follows the XDG Base Directory Specification
// with additional support for EXO-specific overrides.
//
// Priority order:
// 1. EXO_DATA_HOME environment variable
// 2. XDG_DATA_HOME environment variable + "/exo"
// 3. $HOME/.local/share/exo
func getDataHome(home string) string {
	if dataHome := os.Getenv(envDataHome); dataHome != "" {
		return sanitizePath(dataHome, home)
	}

	xdgData := utils.GetXDGDataHome()
	if xdgData == "" {
		xdgData = filepath.Join(home, defaultXDGData)
	} else {
		xdgData = sanitizePath(xdgData, home)
	}

	return filepath.Join(xdgData, "exo")
}

// sanitizePath cleans and normalizes the provided path
func sanitizePath(path, home string) string {
	expanded := utils.ExpandPath(path)
	cleaned := filepath.Clean(expanded)
	if !filepath.IsAbs(cleaned) {
		cleaned = filepath.Join(home, cleaned)
	}
	return cleaned
}

// ensureDirectories creates all necessary directories for the application.
// It returns an error if any directory cannot be created.
func ensureDirectories(paths ...string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	for _, path := range paths {
		absPath := sanitizePath(path, home)
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", absPath, err)
		}
		logger.Info("Created directory", logger.Field{Key: "path", Value: absPath})
	}
	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	logger.Info("Validating configuration")

	if c.General.Editor == "" {
		return fmt.Errorf("editor cannot be empty")
	}

	if c.Dir.DataHome == "" {
		return fmt.Errorf("data_home cannot be empty")
	}

	if c.Dir.TemplateDir == "" {
		return fmt.Errorf("template_dir cannot be empty")
	}

	if c.Dir.PeriodicDir == "" {
		return fmt.Errorf("periodic_dir cannot be empty")
	}

	if c.Dir.ZettelDir == "" {
		return fmt.Errorf("zettel_dir cannot be empty")
	}

	logger.Info("Configuration validation passed")
	return nil
}

func (c *Config) String() string {
	var sb strings.Builder
	sb.WriteString("Configuration:\n")
	sb.WriteString("-------------\n\n")

	// General section
	sb.WriteString("General:\n")
	sb.WriteString(fmt.Sprintf("  editor:        %s\n\n", c.General.Editor))

	// Directories section
	sb.WriteString("Directories:\n")
	sb.WriteString(fmt.Sprintf("  data_home:     %s\n", c.Dir.DataHome))
	sb.WriteString(fmt.Sprintf("  template_dir:  %s\n", c.Dir.TemplateDir))
	sb.WriteString(fmt.Sprintf("  periodic_dir:  %s\n", c.Dir.PeriodicDir))
	sb.WriteString(fmt.Sprintf("  zettel_dir:    %s\n", c.Dir.ZettelDir))
	sb.WriteString(fmt.Sprintf("  projects_dir:  %s\n", c.Dir.ProjectsDir))
	sb.WriteString(fmt.Sprintf("  inbox_dir:     %s\n", c.Dir.InboxDir))
	sb.WriteString(fmt.Sprintf("  idea_dir:      %s\n\n", c.Dir.IdeaDir))

	// Logging section
	sb.WriteString("Logging:\n")
	sb.WriteString(fmt.Sprintf("  level:         %s\n", c.Log.Level))
	sb.WriteString(fmt.Sprintf("  format:        %s\n", c.Log.Format))
	sb.WriteString(fmt.Sprintf("  output:        %s\n", c.Log.Output))

	return sb.String()
}

// Save persists the current configuration to file
func (c *Config) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(home, ".config", "exo")
	configPath := filepath.Join(configDir, "config.yaml")

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(configPath)

	v.Set("general", c.General)
	v.Set("dir", c.Dir)
	v.Set("log", c.Log)

	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
