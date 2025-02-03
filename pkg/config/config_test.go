package config

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/a-kostevski/exo/pkg/logger"
	"github.com/a-kostevski/exo/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resetInstance is a helper function for tests to reset the singleton instance
func resetInstance() {
	cfg = nil
	// Create a new sync.Once
	once = sync.Once{}
}

func setupTest(t *testing.T) (string, func()) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Initialize logger with minimal configuration
	err := logger.Initialize(logger.Config{
		Level:  logger.InfoLevel,
		Format: logger.TextFormat,
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
			expected:    filepath.Join(tmpDir, "custom/data/home"),
			description: "Should expand tilde in EXO_DATA_HOME",
		},
		{
			name:        "EXO_DATA_HOME with relative path",
			dataHome:    "custom/data/home",
			xdgData:     "/should/not/use/this",
			expected:    filepath.Join(tmpDir, "custom/data/home"),
			description: "Should convert relative paths to absolute",
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
			expected:    filepath.Join(tmpDir, "xdg/data", "exo"),
			description: "Should expand tilde in XDG_DATA_HOME",
		},
		{
			name:        "XDG_DATA_HOME with relative path",
			dataHome:    "",
			xdgData:     "xdg/data",
			expected:    filepath.Join(tmpDir, "xdg/data", "exo"),
			description: "Should convert relative XDG paths to absolute",
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
			os.Setenv("HOME", tmpDir)

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
				General: GeneralConfig{
					Editor: defaultEditor,
				},
				Dir: DirConfig{
					DataHome:    "/data/home",
					TemplateDir: "/template/dir",
					PeriodicDir: "/periodic/dir",
					ZettelDir:   "/zettel/dir",
				},
				Log: logger.Config{
					Level:  logger.InfoLevel,
					Format: logger.TextFormat,
					Output: defaultLogOutput,
				},
			},
			expectError: false,
			description: "Valid configuration should not return error",
		},
		{
			name: "empty editor",
			config: Config{
				Dir: DirConfig{
					DataHome:    "/data/home",
					TemplateDir: "/template/dir",
					PeriodicDir: "/periodic/dir",
					ZettelDir:   "/zettel/dir",
				},
			},
			expectError: true,
			description: "Empty editor should return error",
		},
		{
			name: "empty data home",
			config: Config{
				General: GeneralConfig{
					Editor: defaultEditor,
				},
				Dir: DirConfig{
					TemplateDir: "/template/dir",
					PeriodicDir: "/periodic/dir",
					ZettelDir:   "/zettel/dir",
				},
			},
			expectError: true,
			description: "Empty data_home should return error",
		},
		{
			name: "empty template dir",
			config: Config{
				General: GeneralConfig{
					Editor: defaultEditor,
				},
				Dir: DirConfig{
					DataHome:    "/data/home",
					PeriodicDir: "/periodic/dir",
					ZettelDir:   "/zettel/dir",
				},
			},
			expectError: true,
			description: "Empty template_dir should return error",
		},
		{
			name: "empty periodic dir",
			config: Config{
				General: GeneralConfig{
					Editor: defaultEditor,
				},
				Dir: DirConfig{
					DataHome:    "/data/home",
					TemplateDir: "/template/dir",
					ZettelDir:   "/zettel/dir",
				},
			},
			expectError: true,
			description: "Empty periodic_dir should return error",
		},
		{
			name: "empty zettel dir",
			config: Config{
				General: GeneralConfig{
					Editor: defaultEditor,
				},
				Dir: DirConfig{
					DataHome:    "/data/home",
					TemplateDir: "/template/dir",
					PeriodicDir: "/periodic/dir",
				},
			},
			expectError: true,
			description: "Empty zettel_dir should return error",
		},
		{
			name: "relative paths",
			config: Config{
				General: GeneralConfig{
					Editor: defaultEditor,
				},
				Dir: DirConfig{
					DataHome:    "data/home",
					TemplateDir: "template/dir",
					PeriodicDir: "periodic/dir",
					ZettelDir:   "zettel/dir",
				},
			},
			expectError: false,
			description: "Relative paths should be valid",
		},
		{
			name: "paths with tilde",
			config: Config{
				General: GeneralConfig{
					Editor: defaultEditor,
				},
				Dir: DirConfig{
					DataHome:    "~/data/home",
					TemplateDir: "~/data/templates",
					PeriodicDir: "~/data/periodic",
					ZettelDir:   "~/data/zettel",
				},
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

		// Reset the config instance before testing initialization
		resetInstance()

		// Initialize with empty config path should use defaults
		err := Initialize("")
		require.NoError(t, err, "Initialize should succeed with default values when config file doesn't exist")

		// Verify the config was initialized
		cfg, err := Get()
		require.NoError(t, err, "Get should return the initialized config")
		require.NotNil(t, cfg, "Config should not be nil")

		// Verify default values
		assert.Equal(t, defaultEditor, cfg.General.Editor, "Editor should have default value")
		assert.NotEmpty(t, cfg.Dir.DataHome, "DataHome should not be empty")
		assert.True(t, strings.HasSuffix(cfg.Dir.DataHome, "exo"), "DataHome should end with 'exo'")
	})

	t.Run("initialize with invalid config file path", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		// Reset the config instance
		resetInstance()

		err := Initialize("/nonexistent/path/config.yaml")
		assert.Error(t, err, "Initialize should fail with invalid config file path")

		// Verify config is not initialized
		cfg, err := Get()
		assert.Error(t, err, "Get should return error when initialization failed")
		assert.Nil(t, cfg, "Config should be nil when initialization failed")
	})

	t.Run("initialize with valid config file", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Reset the config instance
		resetInstance()

		// Create a valid config file
		configDir := filepath.Join(tmpDir, ".config", "exo")
		require.NoError(t, utils.EnsureDirectories(configDir))

		configFile := filepath.Join(configDir, "config.yaml")
		configContent := []byte(`
general:
  editor: code
dir:
  data_home: ~/data
  template_dir: ~/data/templates
  periodic_dir: ~/data/periodic
  zettel_dir: ~/data/zettel
  projects_dir: ~/data/projects
  inbox_dir: ~/data/0-inbox
  idea_dir: ~/data/ideas
log:
  level: debug
  format: json
  output: stdout
`)
		require.NoError(t, utils.WriteFile(configFile, configContent))

		err := Initialize(configFile)
		require.NoError(t, err, "Initialize should succeed with valid config file")

		cfg, err := Get()
		require.NoError(t, err, "Get should return the initialized config")
		require.NotNil(t, cfg, "Config should not be nil")

		// Verify loaded values
		assert.Equal(t, "code", cfg.General.Editor, "Editor should be loaded from config")
		assert.Equal(t, filepath.Join(tmpDir, "data"), cfg.Dir.DataHome, "DataHome should be loaded from config")
		assert.Equal(t, "debug", string(cfg.Log.Level), "Log level should be loaded from config")
	})

	t.Run("initialize with invalid config content", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		// Reset the config instance
		resetInstance()

		// Create an invalid config file
		configDir := filepath.Join(os.TempDir(), ".config", "exo")
		require.NoError(t, utils.EnsureDirectories(configDir))

		configFile := filepath.Join(configDir, "config.yaml")
		configContent := []byte(`
general:
  editor: [invalid, yaml]
`)
		require.NoError(t, utils.WriteFile(configFile, configContent))

		err := Initialize(configFile)
		assert.Error(t, err, "Initialize should fail with invalid config content")

		// Verify config is not initialized
		cfg, err := Get()
		assert.Error(t, err, "Get should return error when initialization failed")
		assert.Nil(t, cfg, "Config should be nil when initialization failed")
	})

	t.Run("initialize twice", func(t *testing.T) {
		_, cleanup := setupTest(t)
		defer cleanup()

		// Reset the config instance
		resetInstance()

		// First initialization
		err1 := Initialize("")
		require.NoError(t, err1, "First initialization should succeed")

		// Get the first instance
		cfg1, err := Get()
		require.NoError(t, err, "Get should return the first initialized config")
		require.NotNil(t, cfg1, "First config should not be nil")

		// Second initialization
		err2 := Initialize("")
		require.NoError(t, err2, "Second initialization should succeed but not change instance")

		// Get the second instance
		cfg2, err := Get()
		require.NoError(t, err, "Get should return the same config")
		require.NotNil(t, cfg2, "Second config should not be nil")

		// Compare instances
		assert.Equal(t, cfg1, cfg2, "Multiple initializations should return the same instance")
	})

	t.Run("initialize with tilde paths", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Reset the config instance
		resetInstance()

		// Create a valid config file with tilde paths
		configDir := filepath.Join(tmpDir, ".config", "exo")
		require.NoError(t, utils.EnsureDirectories(configDir))

		configFile := filepath.Join(configDir, "config.yaml")
		configContent := []byte(`
general:
  editor: code
dir:
  data_home: ~/data
  template_dir: ~/data/templates
  periodic_dir: ~/data/periodic
  zettel_dir: ~/data/zettel
  projects_dir: ~/data/projects
  inbox_dir: ~/data/0-inbox
  idea_dir: ~/data/ideas
log:
  level: debug
  format: json
  output: stdout
`)
		require.NoError(t, utils.WriteFile(configFile, configContent))

		err := Initialize(configFile)
		require.NoError(t, err, "Initialize should succeed with tilde paths")

		cfg, err := Get()
		require.NoError(t, err, "Get should return the initialized config")
		require.NotNil(t, cfg, "Config should not be nil")

		// Check that paths are expanded
		home := tmpDir // In our test environment, HOME is set to tmpDir
		assert.Equal(t, filepath.Join(home, "data"), cfg.Dir.DataHome, "DataHome should have expanded tilde")
		assert.Equal(t, filepath.Join(home, "data/templates"), cfg.Dir.TemplateDir, "TemplateDir should have expanded tilde")
		assert.Equal(t, filepath.Join(home, "data/periodic"), cfg.Dir.PeriodicDir, "PeriodicDir should have expanded tilde")
		assert.Equal(t, filepath.Join(home, "data/zettel"), cfg.Dir.ZettelDir, "ZettelDir should have expanded tilde")
	})
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
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				cfg, err := Get()
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}()
		}

		// Wait for all goroutines to complete
		wg.Wait()
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
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				assert.NotPanics(t, func() {
					cfg := MustGet()
					assert.NotNil(t, cfg)
				})
			}()
		}

		// Wait for all goroutines to complete
		wg.Wait()
	})
}

func TestConfig_Save(t *testing.T) {
	t.Run("save new config", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Initialize with default config
		err := Initialize("")
		require.NoError(t, err)

		cfg := MustGet()
		cfg.General.Editor = "test-editor"
		cfg.Log.Level = logger.DebugLevel

		// Save the config
		err = cfg.Save()
		require.NoError(t, err)

		// Verify the config file was created
		configPath := filepath.Join(tmpDir, ".config", "exo", "config.yaml")
		_, err = os.Stat(configPath)
		assert.NoError(t, err)

		// Read the saved config
		content, err := utils.ReadFile(configPath)
		require.NoError(t, err)

		// Verify content
		assert.Contains(t, string(content), "test-editor")
		assert.Contains(t, string(content), "debug")
	})

	t.Run("save with invalid permissions", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("Skipping test when running as root")
		}

		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Create config directory with read-only permissions
		configDir := filepath.Join(tmpDir, ".config", "exo")
		require.NoError(t, os.MkdirAll(filepath.Dir(configDir), 0755))
		require.NoError(t, os.Mkdir(configDir, 0555))

		err := Initialize("")
		require.NoError(t, err)

		cfg := MustGet()
		err = cfg.Save()
		assert.Error(t, err)
	})
}

func TestDefaultConfig(t *testing.T) {
	t.Run("default values are set correctly", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Set HOME to ensure consistent test environment
		os.Setenv("HOME", tmpDir)
		defaults := defaultConfig()

		assert.Equal(t, defaultEditor, defaults["general.editor"])
		assert.Equal(t, defaultLogLevel, defaults["log.level"])
		assert.Equal(t, defaultLogFormat, defaults["log.format"])
		assert.Equal(t, defaultLogOutput, defaults["log.output"])

		// Check directory defaults
		expectedDataHome := filepath.Join(filepath.Join(tmpDir, defaultXDGData), "exo")
		assert.Equal(t, expectedDataHome, defaults["dir.data_home"])
		assert.Equal(t, filepath.Join(expectedDataHome, "templates"), defaults["dir.template_dir"])
		assert.Equal(t, filepath.Join(expectedDataHome, "periodic"), defaults["dir.periodic_dir"])
		assert.Equal(t, filepath.Join(expectedDataHome, "zettel"), defaults["dir.zettel_dir"])
		assert.Equal(t, filepath.Join(expectedDataHome, "projects"), defaults["dir.projects_dir"])
		assert.Equal(t, filepath.Join(expectedDataHome, "0-inbox"), defaults["dir.inbox_dir"])
		assert.Equal(t, filepath.Join(expectedDataHome, "ideas"), defaults["dir.idea_dir"])
	})

	t.Run("data_home uses XDG_DATA_HOME", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		customXDGPath := filepath.Join(tmpDir, "custom/xdg/data")
		os.Setenv("XDG_DATA_HOME", customXDGPath)
		defaults := defaultConfig()

		assert.Equal(t, filepath.Join(customXDGPath, "exo"), defaults["dir.data_home"])
	})
}

func TestConfigValidationWithPaths(t *testing.T) {
	t.Run("validate paths with special characters", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Create test directories to ensure they exist
		testDirs := []string{
			filepath.Join(tmpDir, "data home with spaces"),
			filepath.Join(tmpDir, "template-dir-with-dashes"),
			filepath.Join(tmpDir, "periodic_dir_with_underscores"),
			filepath.Join(tmpDir, "zettel.dir.with.dots"),
		}

		for _, dir := range testDirs {
			require.NoError(t, utils.EnsureDirectories(dir))
		}

		cfg := &Config{
			General: GeneralConfig{
				Editor: defaultEditor,
			},
			Dir: DirConfig{
				DataHome:    testDirs[0],
				TemplateDir: testDirs[1],
				PeriodicDir: testDirs[2],
				ZettelDir:   testDirs[3],
			},
			Log: logger.Config{
				Level:  logger.InfoLevel,
				Format: logger.TextFormat,
				Output: defaultLogOutput,
			},
		}

		err := cfg.Validate()
		assert.NoError(t, err)

		// Verify directories exist and are accessible
		for _, dir := range testDirs {
			_, err := os.Stat(dir)
			assert.NoError(t, err)
		}
	})

	t.Run("validate paths with unicode characters", func(t *testing.T) {
		tmpDir, cleanup := setupTest(t)
		defer cleanup()

		// Create test directories with unicode names
		testDirs := []string{
			filepath.Join(tmpDir, "数据"),
			filepath.Join(tmpDir, "模板"),
			filepath.Join(tmpDir, "периодический"),
			filepath.Join(tmpDir, "zettel"),
		}

		for _, dir := range testDirs {
			require.NoError(t, utils.EnsureDirectories(dir))
		}

		cfg := &Config{
			General: GeneralConfig{
				Editor: defaultEditor,
			},
			Dir: DirConfig{
				DataHome:    testDirs[0],
				TemplateDir: testDirs[1],
				PeriodicDir: testDirs[2],
				ZettelDir:   testDirs[3],
			},
			Log: logger.Config{
				Level:  logger.InfoLevel,
				Format: logger.TextFormat,
				Output: defaultLogOutput,
			},
		}

		err := cfg.Validate()
		assert.NoError(t, err)

		// Verify unicode directories exist and are accessible
		for _, dir := range testDirs {
			_, err := os.Stat(dir)
			assert.NoError(t, err)
		}
	})
}
