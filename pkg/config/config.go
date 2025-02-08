package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Environment variables for configuration overrides.
const (
	envDataHome = "EXO_DATA_HOME"
)

// Default configuration values.
const (
	defaultEditor    = "nvim"
	defaultLogLevel  = "info"
	defaultLogFormat = "text"
	defaultLogOutput = "stdout"
)

// Config represents the main configuration structure.
type Config struct {
	General GeneralConfig `mapstructure:"general"`
	Dir     DirConfig     `mapstructure:"dir"`
	Log     LogConfig     `mapstructure:"log"`
}

// GeneralConfig holds general configuration values.
type GeneralConfig struct {
	Editor string `mapstructure:"editor"`
}

// DirConfig holds directory-related configuration.
type DirConfig struct {
	DataHome    string `mapstructure:"data_home"`
	TemplateDir string `mapstructure:"template_dir"`
	PeriodicDir string `mapstructure:"periodic_dir"`
	ZettelDir   string `mapstructure:"zettel_dir"`
	ProjectsDir string `mapstructure:"projects_dir"`
	InboxDir    string `mapstructure:"inbox_dir"`
	IdeaDir     string `mapstructure:"idea_dir"`
}

// LogConfig holds logging configuration.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// NewConfig creates a new configuration instance.
// If configPath is non‑empty, it attempts to load configuration from that file,
// otherwise defaults (plus environment overrides) are used.
func NewConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Set default values.
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

	// If a config file is provided, read it.
	if configPath != "" {
		if _, err := os.Stat(configPath); err != nil {
			return nil, fmt.Errorf("config file not accessible: %w", err)
		}
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Otherwise, add the default config search path.
		v.AddConfigPath(filepath.Join(home, ".config", "exo"))
	}

	if err := v.ReadInConfig(); err != nil {
		// Only return error if specific config file was requested
		if configPath != "" {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand and sanitize directory paths.
	cfg.Dir.DataHome = sanitizePath(cfg.Dir.DataHome, home)
	cfg.Dir.TemplateDir = sanitizePath(cfg.Dir.TemplateDir, home)
	cfg.Dir.PeriodicDir = sanitizePath(cfg.Dir.PeriodicDir, home)
	cfg.Dir.ZettelDir = sanitizePath(cfg.Dir.ZettelDir, home)
	cfg.Dir.ProjectsDir = sanitizePath(cfg.Dir.ProjectsDir, home)
	cfg.Dir.InboxDir = sanitizePath(cfg.Dir.InboxDir, home)
	cfg.Dir.IdeaDir = sanitizePath(cfg.Dir.IdeaDir, home)

	// Apply environment variable override for editor.
	if editor := os.Getenv("EDITOR"); editor != "" {
		cfg.General.Editor = editor
	}

	// Validate configuration.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// getDataHome determines the data home directory.
// Priority: EXO_DATA_HOME environment variable, else $HOME/.local/share/exo.
func getDataHome(home string) string {
	if dataHome := os.Getenv(envDataHome); dataHome != "" {
		return sanitizePath(dataHome, home)
	}
	return filepath.Join(home, ".local", "share", "exo")
}

// sanitizePath expands the tilde and converts relative paths to absolute based on home.
func sanitizePath(path, home string) string {
	// If path starts with ~/, replace with home directory.
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	}
	path = filepath.Clean(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(home, path)
	}
	return path
}

// Validate checks that required configuration fields are non‑empty.
func (c *Config) Validate() error {
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
	return nil
}

// Save writes the configuration to $HOME/.config/exo/config.yaml.
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
	v.Set("general", c.General)
	v.Set("dir", c.Dir)
	v.Set("log", c.Log)

	if err := v.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// String returns a human‑readable representation of the configuration.
func (c *Config) String() string {
	var sb strings.Builder
	sb.WriteString("Configuration:\n")
	sb.WriteString("-------------\n\n")
	sb.WriteString("General:\n")
	sb.WriteString(fmt.Sprintf("  editor:        %s\n\n", c.General.Editor))
	sb.WriteString("Directories:\n")
	sb.WriteString(fmt.Sprintf("  data_home:     %s\n", c.Dir.DataHome))
	sb.WriteString(fmt.Sprintf("  template_dir:  %s\n", c.Dir.TemplateDir))
	sb.WriteString(fmt.Sprintf("  periodic_dir:  %s\n", c.Dir.PeriodicDir))
	sb.WriteString(fmt.Sprintf("  zettel_dir:    %s\n", c.Dir.ZettelDir))
	sb.WriteString(fmt.Sprintf("  projects_dir:  %s\n", c.Dir.ProjectsDir))
	sb.WriteString(fmt.Sprintf("  inbox_dir:     %s\n", c.Dir.InboxDir))
	sb.WriteString(fmt.Sprintf("  idea_dir:      %s\n\n", c.Dir.IdeaDir))
	sb.WriteString("Logging:\n")
	sb.WriteString(fmt.Sprintf("  level:         %s\n", c.Log.Level))
	sb.WriteString(fmt.Sprintf("  format:        %s\n", c.Log.Format))
	sb.WriteString(fmt.Sprintf("  output:        %s\n", c.Log.Output))
	return sb.String()
}

// package config
//
// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"
//
// 	"github.com/spf13/viper"
// )
//
// // Environment variables for configuration overrides.
// const (
// 	envDataHome = "EXO_DATA_HOME"
// )
//
// // Default configuration values.
// const (
// 	defaultEditor    = "nvim"
// 	defaultLogLevel  = "info"
// 	defaultLogFormat = "text"
// 	defaultLogOutput = "stdout"
// )
//
// // Config represents the main configuration structure.
// type Config struct {
// 	General GeneralConfig `mapstructure:"general"`
// 	Dir     DirConfig     `mapstructure:"dir"`
// 	Log     LogConfig     `mapstructure:"log"`
// }
//
// // GeneralConfig holds general configuration values.
// type GeneralConfig struct {
// 	Editor string `mapstructure:"editor"`
// }
//
// // DirConfig holds directory-related configuration.
// type DirConfig struct {
// 	DataHome    string `mapstructure:"data_home"`
// 	TemplateDir string `mapstructure:"template_dir"`
// 	PeriodicDir string `mapstructure:"periodic_dir"`
// 	ZettelDir   string `mapstructure:"zettel_dir"`
// 	ProjectsDir string `mapstructure:"projects_dir"`
// 	InboxDir    string `mapstructure:"inbox_dir"`
// 	IdeaDir     string `mapstructure:"idea_dir"`
// }
//
// // LogConfig holds logging configuration.
// type LogConfig struct {
// 	Level  string `mapstructure:"level"`
// 	Format string `mapstructure:"format"`
// 	Output string `mapstructure:"output"`
// }
//
// // NewConfig creates a new configuration instance.
// // If configPath is non‑empty, it attempts to load configuration from that file,
// // otherwise defaults (plus environment overrides) are used.
// func NewConfig(configPath string) (*Config, error) {
// 	v := viper.New()
// 	v.SetConfigName("config")
// 	v.SetConfigType("yaml")
//
// 	if configPath != "" {
// 		if _, err := os.Stat(configPath); err != nil {
// 			return nil, fmt.Errorf("config file not accessible: %w", err)
// 		}
// 		v.SetConfigFile(configPath)
// 		if err := v.ReadInConfig(); err != nil {
// 			return nil, fmt.Errorf("failed to read config file: %w", err)
// 		}
// 	}
//
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get user home directory: %w", err)
// 	}
// 	v.AddConfigPath(filepath.Join(home, ".config", "exo"))
//
// 	// Set default values.
// 	v.SetDefault("general.editor", defaultEditor)
// 	v.SetDefault("log.level", defaultLogLevel)
// 	v.SetDefault("log.format", defaultLogFormat)
// 	v.SetDefault("log.output", defaultLogOutput)
//
// 	dataHome := getDataHome(home)
// 	v.SetDefault("dir.data_home", dataHome)
// 	v.SetDefault("dir.template_dir", filepath.Join(dataHome, "templates"))
// 	v.SetDefault("dir.periodic_dir", filepath.Join(dataHome, "periodic"))
// 	v.SetDefault("dir.zettel_dir", filepath.Join(dataHome, "zettel"))
// 	v.SetDefault("dir.projects_dir", filepath.Join(dataHome, "projects"))
// 	v.SetDefault("dir.inbox_dir", filepath.Join(dataHome, "0-inbox"))
// 	v.SetDefault("dir.idea_dir", filepath.Join(dataHome, "ideas"))
//
// 	var cfg Config
// 	if err := v.Unmarshal(&cfg); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
// 	}
// 	// Expand and sanitize directory paths.
// 	cfg.Dir.DataHome = sanitizePath(cfg.Dir.DataHome, home)
// 	cfg.Dir.TemplateDir = sanitizePath(cfg.Dir.TemplateDir, home)
// 	cfg.Dir.PeriodicDir = sanitizePath(cfg.Dir.PeriodicDir, home)
// 	cfg.Dir.ZettelDir = sanitizePath(cfg.Dir.ZettelDir, home)
// 	cfg.Dir.ProjectsDir = sanitizePath(cfg.Dir.ProjectsDir, home)
// 	cfg.Dir.InboxDir = sanitizePath(cfg.Dir.InboxDir, home)
// 	cfg.Dir.IdeaDir = sanitizePath(cfg.Dir.IdeaDir, home)
//
// 	// Apply environment variable override for editor.
// 	if editor := os.Getenv("EDITOR"); editor != "" {
// 		cfg.General.Editor = editor
// 	}
//
// 	// Validate configuration.
// 	if err := cfg.Validate(); err != nil {
// 		return nil, err
// 	}
//
// 	return &cfg, nil
// }
//
// // getDataHome determines the data home directory.
// // Priority: EXO_DATA_HOME environment variable, else $HOME/.local/share/exo.
// func getDataHome(home string) string {
// 	if dataHome := os.Getenv(envDataHome); dataHome != "" {
// 		return sanitizePath(dataHome, home)
// 	}
// 	// Return direct join without going through sanitizePath since this is already absolute
// 	return filepath.Join(home, ".local", "share", "exo")
// }
//
// // sanitizePath expands the tilde and converts relative paths to absolute based on home.
// func sanitizePath(path, home string) string {
// 	// If path starts with ~/, replace with home directory
// 	if strings.HasPrefix(path, "~/") {
// 		path = filepath.Join(home, path[2:])
// 	}
//
// 	// Clean the path
// 	path = filepath.Clean(path)
//
// 	// If path is not absolute, make it absolute relative to home
// 	if !filepath.IsAbs(path) {
// 		path = filepath.Join(home, path)
// 	}
//
// 	return path
// }
//
// // Validate checks that required configuration fields are non‑empty.
// func (c *Config) Validate() error {
// 	if c.General.Editor == "" {
// 		return fmt.Errorf("editor cannot be empty")
// 	}
// 	if c.Dir.DataHome == "" {
// 		return fmt.Errorf("data_home cannot be empty")
// 	}
// 	if c.Dir.TemplateDir == "" {
// 		return fmt.Errorf("template_dir cannot be empty")
// 	}
// 	if c.Dir.PeriodicDir == "" {
// 		return fmt.Errorf("periodic_dir cannot be empty")
// 	}
// 	if c.Dir.ZettelDir == "" {
// 		return fmt.Errorf("zettel_dir cannot be empty")
// 	}
// 	return nil
// }
//
// // Save writes the configuration to $HOME/.config/exo/config.yaml.
// func (c *Config) Save() error {
// 	home, err := os.UserHomeDir()
// 	if err != nil {
// 		return fmt.Errorf("failed to get user home directory: %w", err)
// 	}
//
// 	configDir := filepath.Join(home, ".config", "exo")
// 	configPath := filepath.Join(configDir, "config.yaml")
//
// 	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
// 		return fmt.Errorf("failed to create config directory: %w", err)
// 	}
//
// 	v := viper.New()
// 	v.SetConfigType("yaml")
// 	v.Set("general", c.General)
// 	v.Set("dir", c.Dir)
// 	v.Set("log", c.Log)
//
// 	if err := v.WriteConfigAs(configPath); err != nil {
// 		return fmt.Errorf("failed to write config file: %w", err)
// 	}
//
// 	return nil
// }
//
// // String returns a human‑readable representation of the configuration.
// func (c *Config) String() string {
// 	var sb strings.Builder
// 	sb.WriteString("Configuration:\n")
// 	sb.WriteString("-------------\n\n")
// 	sb.WriteString("General:\n")
// 	sb.WriteString(fmt.Sprintf("  editor:        %s\n\n", c.General.Editor))
// 	sb.WriteString("Directories:\n")
// 	sb.WriteString(fmt.Sprintf("  data_home:     %s\n", c.Dir.DataHome))
// 	sb.WriteString(fmt.Sprintf("  template_dir:  %s\n", c.Dir.TemplateDir))
// 	sb.WriteString(fmt.Sprintf("  periodic_dir:  %s\n", c.Dir.PeriodicDir))
// 	sb.WriteString(fmt.Sprintf("  zettel_dir:    %s\n", c.Dir.ZettelDir))
// 	sb.WriteString(fmt.Sprintf("  projects_dir:  %s\n", c.Dir.ProjectsDir))
// 	sb.WriteString(fmt.Sprintf("  inbox_dir:     %s\n", c.Dir.InboxDir))
// 	sb.WriteString(fmt.Sprintf("  idea_dir:      %s\n\n", c.Dir.IdeaDir))
// 	sb.WriteString("Logging:\n")
// 	sb.WriteString(fmt.Sprintf("  level:         %s\n", c.Log.Level))
// 	sb.WriteString(fmt.Sprintf("  format:        %s\n", c.Log.Format))
// 	sb.WriteString(fmt.Sprintf("  output:        %s\n", c.Log.Output))
// 	return sb.String()
// }
