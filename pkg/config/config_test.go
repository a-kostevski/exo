package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/a-kostevski/exo/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_Defaults(t *testing.T) {
	// Use a temporary directory to simulate HOME.
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	os.Setenv("HOME", tmpHome)
	// Ensure EXO_DATA_HOME is not set.
	os.Unsetenv("EXO_DATA_HOME")

	cfg, err := config.NewConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify default values.
	assert.Equal(t, "nvim", cfg.General.Editor)

	// Expected data home: tmpHome/.local/share/exo
	expectedDataHome := filepath.Join(tmpHome, ".local", "share", "exo")
	assert.Equal(t, expectedDataHome, cfg.Dir.DataHome)

	// Check that other directories are set relative to data home.
	assert.Equal(t, filepath.Join(expectedDataHome, "templates"), cfg.Dir.TemplateDir)
	assert.Equal(t, filepath.Join(expectedDataHome, "periodic"), cfg.Dir.PeriodicDir)
	assert.Equal(t, filepath.Join(expectedDataHome, "zettel"), cfg.Dir.ZettelDir)
	assert.Equal(t, filepath.Join(expectedDataHome, "projects"), cfg.Dir.ProjectsDir)
	assert.Equal(t, filepath.Join(expectedDataHome, "0-inbox"), cfg.Dir.InboxDir)
	assert.Equal(t, filepath.Join(expectedDataHome, "ideas"), cfg.Dir.IdeaDir)

	// Verify logging defaults.
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "text", cfg.Log.Format)
	assert.Equal(t, "stdout", cfg.Log.Output)
}

func TestNewConfig_ConfigFile(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	os.Setenv("HOME", tmpHome)
	os.Unsetenv("EXO_DATA_HOME")
	os.Unsetenv("EDITOR")

	// Create a temporary config file.
	configDir := filepath.Join(tmpHome, ".config", "exo")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	configPath := filepath.Join(configDir, "config.yaml")
	configContent := `
general:
  editor: code
dir:
  data_home: "~/mydata"
  template_dir: "~/mydata/templates"
  periodic_dir: "~/mydata/periodic"
  zettel_dir: "~/mydata/zettel"
  projects_dir: "~/mydata/projects"
  inbox_dir: "~/mydata/0-inbox"
  idea_dir: "~/mydata/ideas"
log:
  level: debug
  format: json
  output: stderr
`
	require.NoError(t, os.WriteFile(configPath, []byte(configContent), 0644))

	cfg, err := config.NewConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "code", cfg.General.Editor)
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	expectedDataHome := filepath.Join(home, "mydata")
	assert.Equal(t, expectedDataHome, cfg.Dir.DataHome)
	assert.Equal(t, filepath.Join(expectedDataHome, "templates"), cfg.Dir.TemplateDir)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, "stderr", cfg.Log.Output)
}

func TestNewConfig_EnvOverride(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	os.Setenv("HOME", tmpHome)
	os.Unsetenv("EXO_DATA_HOME")

	os.Setenv("EDITOR", "vim")

	cfg, err := config.NewConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "vim", cfg.General.Editor)
}

func TestValidate(t *testing.T) {
	cfg := &config.Config{
		General: config.GeneralConfig{
			Editor: "",
		},
		Dir: config.DirConfig{
			DataHome:    "",
			TemplateDir: "templates",
			PeriodicDir: "periodic",
			ZettelDir:   "zettel",
		},
	}
	err := cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "editor cannot be empty")

	cfg.General.Editor = "nvim"
	err = cfg.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "data_home cannot be empty")
}

func TestSaveAndString(t *testing.T) {
	tmpHome := t.TempDir()
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)
	os.Setenv("HOME", tmpHome)
	os.Unsetenv("EXO_DATA_HOME")

	cfg, err := config.NewConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	err = cfg.Save()
	require.NoError(t, err)

	configPath := filepath.Join(tmpHome, ".config", "exo", "config.yaml")
	_, err = os.Stat(configPath)
	require.NoError(t, err)

	str := cfg.String()
	assert.Contains(t, str, "editor")
	assert.Contains(t, str, "data_home")
}
