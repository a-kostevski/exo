package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/a-kostevski/exo/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTest(t *testing.T) (string, func()) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Initialize logger with minimal configuration
	err := logger.Initialize(logger.Config{
		Level:  "panic",
		Format: "text",
		Output: "discard",
	})
	require.NoError(t, err)

	// Save original environment variables
	origHome := os.Getenv("HOME")
	origDataHome := os.Getenv("EXO_DATA_HOME")
	origXDGData := os.Getenv("XDG_DATA_HOME")
	origEditor := os.Getenv("EDITOR")

	// Set HOME to temp directory to avoid reading real config
	os.Setenv("HOME", tmpDir)
	os.Setenv("EXO_DATA_HOME", "")
	os.Setenv("XDG_DATA_HOME", "")
	os.Setenv("EDITOR", "")

	// Reset singleton instance
	resetInstance()

	// Return cleanup function
	cleanup := func() {
		os.Setenv("HOME", origHome)
		os.Setenv("EXO_DATA_HOME", origDataHome)
		os.Setenv("XDG_DATA_HOME", origXDGData)
		os.Setenv("EDITOR", origEditor)
		resetInstance()
	}

	return tmpDir, cleanup
}

func TestGetDataHome(t *testing.T) {
	tmpDir, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name        string
		dataHome    string
		xdgData     string
		expected    string
		description string
	}{
		{
			name:        "EXO_DATA_HOME set",
			dataHome:    "/custom/data/home",
			xdgData:     "/should/not/use/this",
			expected:    "/custom/data/home",
			description: "Should use EXO_DATA_HOME when set",
		},
		{
			name:        "EXO_DATA_HOME with tilde",
			dataHome:    "~/custom/data/home",
			xdgData:     "/should/not/use/this",
			expected:    "~/custom/data/home",
			description: "Should preserve tilde in EXO_DATA_HOME",
		},
		{
			name:        "XDG_DATA_HOME set",
			dataHome:    "",
			xdgData:     "/xdg/data",
			expected:    filepath.Join("/xdg/data", "exo"),
			description: "Should use XDG_DATA_HOME/exo when EXO_DATA_HOME not set",
		},
		{
			name:        "XDG_DATA_HOME with tilde",
			dataHome:    "",
			xdgData:     "~/xdg/data",
			expected:    filepath.Join("~/xdg/data", "exo"),
			description: "Should preserve tilde in XDG_DATA_HOME",
		},
		{
			name:        "no env vars set",
			dataHome:    "",
			xdgData:     "",
			expected:    filepath.Join(tmpDir, ".local/share", "exo"),
			description: "Should use default XDG path when no env vars set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("EXO_DATA_HOME", tt.dataHome)
			os.Setenv("XDG_DATA_HOME", tt.xdgData)

			result := getDataHome(tmpDir)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	_, cleanup := setupTest(t)
	defer cleanup()

	tests := []struct {
		name        string
		config      Config
		expectError bool
		description string
	}{
		{
			name: "valid config",
			config: Config{
				Editor:      defaultEditor,
				DataHome:    "/data/home",
				TemplateDir: "/template/dir",
				PeriodicDir: "/periodic/dir",
				ZettelDir:   "/zettel/dir",
				Log: logger.Config{
					Level:  defaultLogLevel,
					Format: defaultLogFormat,
					Output: defaultLogOutput,
				},
			},
			expectError: false,
			description: "Valid configuration should not return error",
		},
		{
			name: "empty editor",
			config: Config{
				DataHome:    "/data/home",
				TemplateDir: "/template/dir",
				PeriodicDir: "/periodic/dir",
				ZettelDir:   "/zettel/dir",
			},
			expectError: true,
			description: "Empty editor should return error",
		},
		{
			name: "empty data home",
			config: Config{
				Editor:      defaultEditor,
				TemplateDir: "/template/dir",
				PeriodicDir: "/periodic/dir",
				ZettelDir:   "/zettel/dir",
			},
			expectError: true,
			description: "Empty data_home should return error",
		},
		{
			name: "empty template dir",
			config: Config{
				Editor:      defaultEditor,
				DataHome:    "/data/home",
				PeriodicDir: "/periodic/dir",
				ZettelDir:   "/zettel/dir",
			},
			expectError: true,
			description: "Empty template_dir should return error",
		},
		{
			name: "empty periodic dir",
			config: Config{
				Editor:      defaultEditor,
				DataHome:    "/data/home",
				TemplateDir: "/template/dir",
				ZettelDir:   "/zettel/dir",
			},
			expectError: true,
			description: "Empty periodic_dir should return error",
		},
		{
			name: "empty zettel dir",
			config: Config{
				Editor:      defaultEditor,
				DataHome:    "/data/home",
				TemplateDir: "/template/dir",
				PeriodicDir: "/periodic/dir",
			},
			expectError: true,
			description: "Empty zettel_dir should return error",
		},
		{
			name: "relative paths",
			config: Config{
				Editor:      defaultEditor,
				DataHome:    "data/home",
				TemplateDir: "template/dir",
				PeriodicDir: "periodic/dir",
				ZettelDir:   "zettel/dir",
			},
			expectError: false,
			description: "Relative paths should be valid",
		},
		{
			name: "paths with tilde",
			config: Config{
				Editor:      defaultEditor,
				DataHome:    "~/data/home",
				TemplateDir: "~/template/dir",
				PeriodicDir: "~/periodic/dir",
				ZettelDir:   "~/zettel/dir",
			},
			expectError: false,
			description: "Paths with tilde should be valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	t.Run("initialize with non-existent config file", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		err := Initialize("")
		require.NoError(t, err, "Initialize should succeed with default values when config file doesn't exist")

		cfg, err := Get()
		require.NoError(t, err, "Get should return the initialized config")
		assert.Equal(t, defaultEditor, cfg.Editor, "Editor should have default value")
		assert.NotEmpty(t, cfg.DataHome, "DataHome should not be empty")
	})

	t.Run("initialize with invalid config file path", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		err := Initialize("/nonexistent/path/config.yaml")
		assert.Error(t, err, "Initialize should fail with invalid config file path")
	})

	t.Run("initialize with valid config file", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Create a valid config file
		configDir := filepath.Join(tmpDir, ".config", "exo")
		require.NoError(t, os.MkdirAll(configDir, 0755))

		configFile := filepath.Join(configDir, "config.yaml")
		configContent := []byte(`
editor: code
data_home: /custom/data
template_dir: /custom/templates
periodic_dir: /custom/periodic
zettel_dir: /custom/zettel
log:
  level: debug
  format: json
  output: stdout
`)
		require.NoError(t, os.WriteFile(configFile, configContent, 0644))

		err := Initialize(configFile)
		require.NoError(t, err, "Initialize should succeed with valid config file")

		cfg, err := Get()
		require.NoError(t, err, "Get should return the initialized config")
		assert.Equal(t, "code", cfg.Editor, "Editor should be loaded from config")
		assert.Equal(t, "/custom/data", cfg.DataHome, "DataHome should be loaded from config")
		assert.Equal(t, "debug", cfg.Log.Level, "Log level should be loaded from config")
	})

	t.Run("initialize with invalid config content", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Create an invalid config file
		configDir := filepath.Join(tmpDir, ".config", "exo")
		require.NoError(t, os.MkdirAll(configDir, 0755))

		configFile := filepath.Join(configDir, "config.yaml")
		configContent := []byte(`
editor: [invalid, yaml]
`)
		require.NoError(t, os.WriteFile(configFile, configContent, 0644))

		err := Initialize(configFile)
		assert.Error(t, err, "Initialize should fail with invalid config content")
	})

	t.Run("initialize twice", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		// First initialization
		err1 := Initialize("")
		require.NoError(t, err1, "First initialization should succeed")

		// Get the first instance
		cfg1, err := Get()
		require.NoError(t, err, "Get should return the first initialized config")

		// Second initialization
		err2 := Initialize("")
		require.NoError(t, err2, "Second initialization should succeed but not change instance")

		// Get the second instance
		cfg2, err := Get()
		require.NoError(t, err, "Get should return the same config")

		// Compare instances
		assert.Equal(t, cfg1, cfg2, "Multiple initializations should return the same instance")
	})

	t.Run("initialize with tilde paths", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Create a valid config file with tilde paths
		configDir := filepath.Join(tmpDir, ".config", "exo")
		require.NoError(t, os.MkdirAll(configDir, 0755))

		configFile := filepath.Join(configDir, "config.yaml")
		configContent := []byte(`
editor: code
data_home: ~/data
template_dir: ~/.config/exo/templates
periodic_dir: ~/notes/periodic
zettel_dir: ~/notes/zettel
log:
  level: debug
  format: json
  output: stdout
`)
		require.NoError(t, os.WriteFile(configFile, configContent, 0644))

		err := Initialize(configFile)
		require.NoError(t, err, "Initialize should succeed with tilde paths")

		cfg, err := Get()
		require.NoError(t, err, "Get should return the initialized config")

		// Check that paths are expanded
		home := tmpDir // In our test environment, HOME is set to tmpDir
		assert.Equal(t, filepath.Join(home, "data"), cfg.DataHome, "DataHome should have expanded tilde")
		assert.Equal(t, filepath.Join(home, ".config/exo/templates"), cfg.TemplateDir, "TemplateDir should have expanded tilde")
		assert.Equal(t, filepath.Join(home, "notes/periodic"), cfg.PeriodicDir, "PeriodicDir should have expanded tilde")
		assert.Equal(t, filepath.Join(home, "notes/zettel"), cfg.ZettelDir, "ZettelDir should have expanded tilde")
	})
}

func TestGetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		dataHome string
	}{
		{
			name:     "absolute_paths",
			dataHome: "/test/data",
		},
		{
			name:     "relative_paths",
			dataHome: "test/data",
		},
		{
			name:     "paths_with_tilde",
			dataHome: "~/test/data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defaults := getDefaults(tt.dataHome)

			// Check basic values
			assert.Equal(t, defaultEditor, defaults["editor"])
			assert.Equal(t, tt.dataHome, defaults["data_home"])
			assert.Equal(t, filepath.Join(tt.dataHome, "templates"), defaults["template_dir"])
			assert.Equal(t, filepath.Join(tt.dataHome, "periodic"), defaults["periodic_dir"])
			assert.Equal(t, filepath.Join(tt.dataHome, "zettel"), defaults["zettel_dir"])

			// Check log configuration
			logConfig := defaults["log"].(map[string]interface{})
			assert.Equal(t, defaultLogLevel, logConfig["level"])
			assert.Equal(t, defaultLogFormat, logConfig["format"])
			assert.Equal(t, defaultLogOutput, logConfig["output"])
		})
	}
}

func TestGet(t *testing.T) {
	t.Run("get before initialize", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		cfg, err := Get()
		assert.Error(t, err, "Get should return error when config is not initialized")
		assert.Nil(t, cfg, "Config should be nil when not initialized")
	})

	t.Run("get after initialize", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		err := Initialize("")
		require.NoError(t, err, "Initialize should succeed")

		cfg, err := Get()
		assert.NoError(t, err, "Get should succeed after initialization")
		assert.NotNil(t, cfg, "Config should not be nil after initialization")
	})

	t.Run("concurrent gets", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		err := Initialize("")
		require.NoError(t, err, "Initialize should succeed")

		// Run multiple goroutines accessing the config
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				cfg, err := Get()
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestMustGet(t *testing.T) {
	t.Run("must get before initialize", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		assert.Panics(t, func() {
			MustGet()
		}, "MustGet should panic when config is not initialized")
	})

	t.Run("must get after initialize", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		err := Initialize("")
		require.NoError(t, err, "Initialize should succeed")

		assert.NotPanics(t, func() {
			cfg := MustGet()
			assert.NotNil(t, cfg, "Config should not be nil after initialization")
		}, "MustGet should not panic after initialization")
	})

	t.Run("concurrent must gets", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		err := Initialize("")
		require.NoError(t, err, "Initialize should succeed")

		// Run multiple goroutines accessing the config
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				assert.NotPanics(t, func() {
					cfg := MustGet()
					assert.NotNil(t, cfg)
				})
				done <- true
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}
