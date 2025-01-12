package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/a-kostevski/exo/internal/fs"
	"github.com/a-kostevski/exo/internal/logger"
	"github.com/spf13/viper"
)

// Environment variables
const (
	envDataHome     = "EXO_DATA_HOME"
	envXDGData      = "XDG_DATA_HOME"
	envXDGCache     = "XDG_CACHE_HOME"
	defaultXDGData  = ".local/share"
	defaultXDGCache = ".cache"
)

// Default configuration values
const (
	defaultEditor    = "vim"
	defaultLogLevel  = "info"
	defaultLogFormat = "text"
	defaultLogOutput = "stderr"
)

// Log output options
const (
	LogOutputStdout = "stdout"
	LogOutputStderr = "stderr"
	LogOutputFile   = "file"
	LogOutputBoth   = "both"
)

// Log format options
const (
	LogFormatText = "text"
	LogFormatJSON = "json"
)

// Log level options
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

var (
	instance *Config
	once     sync.Once
	mu       sync.RWMutex
)

// For testing purposes only
func resetInstance() {
	mu.Lock()
	defer mu.Unlock()
	instance = nil
	once = sync.Once{}
}

// Config holds the application configuration
type Config struct {
	// Editor is the command to open notes
	Editor string `mapstructure:"editor"`

	// Directory paths
	DataHome    string `mapstructure:"data_home"`
	TemplateDir string `mapstructure:"template_dir"`
	PeriodicDir string `mapstructure:"periodic_dir"`
	ZettelDir   string `mapstructure:"zettel_dir"`
	ProjectsDir string `mapstructure:"projects_dir"`
	InboxDir    string `mapstructure:"inbox_dir"`

	// Logging configuration
	Log logger.Config `mapstructure:"log"`
}

// Get returns the global configuration instance.
// It will return an error if the configuration hasn't been loaded yet.
func Get() (*Config, error) {
	mu.RLock()
	defer mu.RUnlock()

	if instance == nil {
		logger.Error("Configuration not initialized")
		return nil, fmt.Errorf("configuration not initialized")
	}
	logger.Debug("Retrieved configuration instance")
	return instance, nil
}

// MustGet returns the global configuration instance.
// It will panic if the configuration hasn't been loaded yet.
func MustGet() *Config {
	cfg, err := Get()
	if err != nil {
		panic(err)
	}
	return cfg
}

// Initialize loads the configuration from the specified path
func Initialize(path string) error {
	// Initialize logger first
	err := logger.Initialize(logger.Config{
		Level:  "panic",
		Format: "text",
		Output: "discard",
	})
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger.Infof("Initializing configuration with file: %s", path)
	var initErr error
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()

		var cfg *Config
		cfg, initErr = load(path)
		if initErr != nil {
			logger.Errorf("Failed to load configuration: %v", initErr)
			return
		}
		instance = cfg
		logger.Info("Configuration initialized successfully")
	})
	return initErr
}

// getDataHome returns the data home directory based on environment variables
func getDataHome(home string) string {
	// Check EXO_DATA_HOME first
	if dataHome := os.Getenv(envDataHome); dataHome != "" {
		logger.Debugf("Using EXO_DATA_HOME: %s", dataHome)
		return dataHome
	}

	// Check XDG_DATA_HOME next
	xdgData := os.Getenv(envXDGData)
	if xdgData == "" {
		xdgData = filepath.Join(home, defaultXDGData)
		logger.Debugf("Using default XDG data home: %s", xdgData)
	} else {
		logger.Debugf("Using XDG_DATA_HOME: %s", xdgData)
	}

	// Return XDG_DATA_HOME/exo
	dataHome := filepath.Join(xdgData, "exo")
	logger.Debugf("Final data home directory: %s", dataHome)
	return dataHome
}

// load loads the configuration from the specified file
func load(cfgFile string) (*Config, error) {
	logger.Debug("Loading configuration")
	v := viper.New()
	v.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err != nil {
		logger.Errorf("Failed to get user home directory: %v", err)
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	logger.Debugf("User home directory: %s", home)

	// Get data home directory
	dataHome := getDataHome(home)

	// Set default values
	defaults := getDefaults(dataHome)
	for key, value := range defaults {
		v.SetDefault(key, value)
		logger.Debugf("Setting default value for %s: %v", key, value)
	}

	// If config file is specified, use it
	if cfgFile != "" {
		logger.Infof("Using specified config file: %s", cfgFile)
		v.SetConfigFile(cfgFile)
	} else {
		// Search in default locations
		configPath := filepath.Join(home, ".config", "exo")
		logger.Debugf("Searching for config in: %s", configPath)
		v.AddConfigPath(configPath)
		v.SetConfigName("config")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logger.Errorf("Failed to read config file: %v", err)
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, use defaults
		logger.Info("No config file found, using defaults")
		logConfig := defaults["log"].(map[string]interface{})
		cfg := &Config{
			Editor:      defaults["editor"].(string),
			DataHome:    defaults["data_home"].(string),
			TemplateDir: defaults["template_dir"].(string),
			PeriodicDir: defaults["periodic_dir"].(string),
			ZettelDir:   defaults["zettel_dir"].(string),
			Log: logger.Config{
				Level:  logConfig["level"].(string),
				Format: logConfig["format"].(string),
				Output: logConfig["output"].(string),
				File:   logConfig["file"].(string),
			},
		}
		return cfg, nil
	}
	logger.Info("Successfully read config file")

	// Unmarshal config
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		logger.Errorf("Failed to unmarshal config: %v", err)
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	logger.Debug("Successfully unmarshaled config")

	// Expand tildes in paths
	cfg.DataHome = fs.ExpandPath(cfg.DataHome)
	cfg.TemplateDir = fs.ExpandPath(cfg.TemplateDir)
	cfg.PeriodicDir = fs.ExpandPath(cfg.PeriodicDir)
	cfg.ZettelDir = fs.ExpandPath(cfg.ZettelDir)
	logger.Debug("Expanded all path variables")

	// Validate config
	if err := cfg.Validate(); err != nil {
		logger.Errorf("Invalid configuration: %v", err)
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	logger.Info("Configuration validation successful")

	return cfg, nil
}

// getDefaults returns a map of default configuration values
func getDefaults(dataHome string) map[string]interface{} {
	// Get cache directory
	cacheDir := os.Getenv(envXDGCache)
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Warnf("Failed to get user home directory: %v, using current directory for cache", err)
			cacheDir = defaultXDGCache
		} else {
			cacheDir = filepath.Join(home, defaultXDGCache)
		}
	}
	logFile := filepath.Join(cacheDir, "exo", "exo.log")
	logger.Debugf("Using log file path: %s", logFile)

	return map[string]interface{}{
		"editor":       defaultEditor,
		"data_home":    dataHome,
		"template_dir": filepath.Join(dataHome, "templates"),
		"periodic_dir": filepath.Join(dataHome, "periodic"),
		"zettel_dir":   filepath.Join(dataHome, "zettel"),
		"projects_dir": filepath.Join(dataHome, "projects"),
		"inbox_dir":    filepath.Join(dataHome, "0-inbox"),
		"log": map[string]interface{}{
			"level":  defaultLogLevel,
			"format": defaultLogFormat,
			"output": defaultLogOutput,
			"file":   logFile,
		},
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	logger.Debug("Validating configuration")

	if c.Editor == "" {
		return fmt.Errorf("editor cannot be empty")
	}

	if c.DataHome == "" {
		return fmt.Errorf("data_home cannot be empty")
	}

	if c.TemplateDir == "" {
		return fmt.Errorf("template_dir cannot be empty")
	}

	if c.PeriodicDir == "" {
		return fmt.Errorf("periodic_dir cannot be empty")
	}

	if c.ZettelDir == "" {
		return fmt.Errorf("zettel_dir cannot be empty")
	}

	logger.Debug("Configuration validation passed")
	return nil
}

// String returns a pretty-printed string representation of the configuration
func (c *Config) String() string {
	var sb strings.Builder
	sb.WriteString("Configuration:\n")
	sb.WriteString("-------------\n")

	// Use reflection to get all fields
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	// Track the current section
	var currentSection string

	// Process each field
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Get the mapstructure tag which represents the config key
		configKey := fieldType.Tag.Get("mapstructure")
		if configKey == "" {
			configKey = fieldType.Name
		}

		// Determine the section based on the field name or type
		var section string
		if fieldType.Name == "Log" {
			section = "Logging"
		} else if strings.HasSuffix(fieldType.Name, "Dir") || fieldType.Name == "DataHome" {
			section = "Directories"
		} else {
			section = "General"
		}

		// Print section header if we're entering a new section
		if section != currentSection {
			if currentSection != "" {
				sb.WriteString("\n")
			}
			sb.WriteString(section + ":\n")
			currentSection = section
		}

		// Handle the Log struct specially
		if fieldType.Name == "Log" {
			logVal := reflect.ValueOf(c.Log)
			logType := logVal.Type()

			// Print each log field
			for j := 0; j < logVal.NumField(); j++ {
				logField := logVal.Field(j)
				logFieldType := logType.Field(j)
				logKey := logFieldType.Tag.Get("mapstructure")
				if logKey == "" {
					logKey = logFieldType.Name
				}

				// Only print File field if it's not empty
				if logFieldType.Name == "File" && logField.String() == "" {
					continue
				}

				sb.WriteString(fmt.Sprintf("  %-13s %v\n", logKey+":", logField.Interface()))
			}
			continue
		}

		// Print regular fields with proper alignment
		sb.WriteString(fmt.Sprintf("  %-13s %v\n", configKey+":", field.Interface()))
	}

	return sb.String()
}

// Save persists the current configuration to file
func Save() error {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		logger.Error("Cannot save: configuration not initialized")
		return fmt.Errorf("configuration not initialized")
	}

	logger.Debug("Starting configuration save")

	// Get user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		logger.Errorf("Failed to get user home directory: %v", err)
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	logger.Debugf("Using home directory: %s", home)

	// Ensure config directory exists
	configDir := filepath.Join(home, ".config", "exo")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		logger.Errorf("Failed to create config directory %s: %v", configDir, err)
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	logger.Debugf("Using config directory: %s", configDir)

	// Create viper instance
	v := viper.New()
	v.SetConfigType("yaml")
	configFile := filepath.Join(configDir, "config.yaml")
	v.SetConfigFile(configFile)
	logger.Debugf("Config file path: %s", configFile)

	// Build configuration map
	config := map[string]interface{}{
		"editor":       instance.Editor,
		"data_home":    instance.DataHome,
		"template_dir": instance.TemplateDir,
		"periodic_dir": instance.PeriodicDir,
		"zettel_dir":   instance.ZettelDir,
		"log": map[string]interface{}{
			"level":  instance.Log.Level,
			"format": instance.Log.Format,
			"output": instance.Log.Output,
		},
	}
	logger.Debugf("Current configuration state: %+v", config)

	// Set all configuration values
	for key, value := range config {
		v.Set(key, value)
		logger.Debugf("Setting config value: %s = %v", key, value)
	}

	// Write config to file, handling the case where it might already exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		logger.Debug("Creating new config file")
		if err := v.WriteConfig(); err != nil {
			logger.Errorf("Failed to write new config file: %v", err)
			return fmt.Errorf("failed to write new config file: %w", err)
		}
		logger.Info("Created new configuration file")
	} else {
		logger.Debug("Updating existing config file")
		if err := v.SafeWriteConfig(); err != nil {
			// If the error is because the file exists, try to overwrite it
			if err := v.WriteConfig(); err != nil {
				logger.Errorf("Failed to update existing config file: %v", err)
				return fmt.Errorf("failed to update existing config file: %w", err)
			}
		}
		logger.Info("Updated existing configuration file")
	}

	logger.Info("Configuration saved successfully")
	return nil
}
