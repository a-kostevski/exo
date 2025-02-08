package note_test

import (
	"path/filepath"
	"testing"

	"github.com/a-kostevski/exo/pkg/note"
	"github.com/a-kostevski/exo/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBaseNote_Success(t *testing.T) {
	// Use t.TempDir() to create a temporary directory.
	tmpDir := t.TempDir()

	// Get dummy dependencies.
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	n, err := note.NewBaseNote("Test Note", cfg, dtm, dl, dfs,
		note.WithSubDir("notes"),
		note.WithFileName("test.md"),
		note.WithTemplateName("default"),
		note.WithContent("Initial Content"),
	)
	require.NoError(t, err)
	assert.NotNil(t, n)

	expectedPath := filepath.Join(cfg.Dir.DataHome, "notes", "test.md")
	assert.Equal(t, expectedPath, n.Path())
	assert.Equal(t, "Initial Content", n.Content())
}

func TestNewBaseNote_Failure_MissingOptions(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)
	_, err := note.NewBaseNote("Test Note", cfg, dtm, dl, dfs)
	require.Error(t, err)
}
