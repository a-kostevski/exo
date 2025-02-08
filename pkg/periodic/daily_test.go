package periodic_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/a-kostevski/exo/pkg/periodic"
	"github.com/a-kostevski/exo/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDailyNote_Initialization(t *testing.T) {
	// When a daily note is created for the first time, it should initialize its content and save the file.
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	date := time.Now().Truncate(24 * time.Hour)
	daily, err := periodic.NewDailyNote(date, cfg, dtm, dl, dfs)
	require.NoError(t, err)
	require.NotNil(t, daily)

	// The file should now exist on disk.
	assert.True(t, daily.Exists(), "Daily note should be created and exist on disk")

	// Verify that the file name is based on the date.
	expectedFile := date.Format("2006-01-02") + ".md"
	expectedPath := filepath.Join(cfg.Dir.DataHome, "day", expectedFile)
	assert.Equal(t, expectedPath, daily.Path())
}

func TestNewDailyNote_LoadExisting(t *testing.T) {
	// Create a daily note, modify its content, save it, then create another instance for the same date,
	// and verify that it loads the updated content.
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	date := time.Now().Truncate(24 * time.Hour)
	// First creation will initialize and save the note.
	daily1, err := periodic.NewDailyNote(date, cfg, dtm, dl, dfs)
	require.NoError(t, err)
	require.NotNil(t, daily1)

	// Change its content and save.
	newContent := "Updated daily note content"
	err = daily1.SetContent(newContent)
	require.NoError(t, err)
	err = daily1.Save()
	require.NoError(t, err)

	// Create another daily note for the same date. It should load the saved content.
	daily2, err := periodic.NewDailyNote(date, cfg, dtm, dl, dfs)
	require.NoError(t, err)
	require.NotNil(t, daily2)

	assert.Equal(t, newContent, daily2.Content())
}

func TestDailyNote_NavigationHelpers(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	// Create a daily note for a known date.
	date := time.Date(2025, 2, 8, 0, 0, 0, 0, time.UTC)
	daily, err := periodic.NewDailyNote(date, cfg, dtm, dl, dfs)
	require.NoError(t, err)

	// Test PreviousOrZero and NextOrZero.
	expectedPrev := date.AddDate(0, 0, -1)
	expectedNext := date.AddDate(0, 0, 1)
	prev := daily.PreviousOrZero()
	next := daily.NextOrZero()

	assert.Equal(t, expectedPrev, prev)
	assert.Equal(t, expectedNext, next)
}

func TestDailyNote_TemplateApplied(t *testing.T) {
	// This test verifies that when a daily note is initialized (i.e. when the file does not exist),
	// the template is applied and the note content is set accordingly.
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	// Create a daily note. Since the file does not exist, it should be initialized.
	date := time.Now().Truncate(24 * time.Hour)
	daily, err := periodic.NewDailyNote(date, cfg, dtm, dl, dfs)
	require.NoError(t, err)

	// Our DummyTemplateManager (used via dtm) returns "Template: unknown" because the provided
	// template data does not include a "Title" field. We therefore expect that value.
	expected := "Template: unknown"
	assert.Equal(t, expected, daily.Content())
}
