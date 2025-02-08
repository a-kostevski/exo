package zettel_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/a-kostevski/exo/pkg/note"
	"github.com/a-kostevski/exo/pkg/testutil"
	"github.com/a-kostevski/exo/pkg/zettel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewZettelNote_Success verifies that a Zettel note is created properly
// using the defaults (subdirectory "zettel", filename derived from title, etc.)
func TestNewZettelNote_Success(t *testing.T) {
	// Create a temporary directory for DataHome.
	tmpDir := t.TempDir()
	// Get dummy dependencies.
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	title := "TestZettel"

	// Create the Zettel note with an initial content.
	zNote, err := zettel.NewZettelNote(title, cfg, dtm, dl, dfs, note.WithContent("Initial Zettel Content"))
	require.NoError(t, err)
	require.NotNil(t, zNote)

	// Expected file path is: DataHome/zettel/<title>.md
	expectedPath := filepath.Join(cfg.Dir.DataHome, "0-inbox", title+".md")
	assert.Equal(t, expectedPath, zNote.Path())
	assert.Equal(t, title, zNote.Title())
	assert.Equal(t, "Initial Zettel Content", zNote.Content())

	// Check that the string representation contains both title and tag.
	str := zNote.String()
	assert.Contains(t, str, title)
}

// TestZettelNote_Validate tests that the Validate method enforces the Zettel-specific rule
// that a tag must be provided.
func TestZettelNote_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	// Valid note with non-empty tag.
	zNote, err := zettel.NewZettelNote("ValidNote", cfg, dtm, dl, dfs)
	require.NoError(t, err)
	err = zNote.Validate()
	require.NoError(t, err)
}

// TestZettelNote_UpdateContent tests that UpdateContent changes the content.
func TestZettelNote_UpdateContent(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	zNote, err := zettel.NewZettelNote("UpdateContentNote", cfg, dtm, dl, dfs,
		note.WithContent("Old Content"),
	)
	require.NoError(t, err)
	assert.Equal(t, "Old Content", zNote.Content())

	err = zNote.SetContent("New Content")
	require.NoError(t, err)
	assert.Equal(t, "New Content", zNote.Content())
}

// TestZettelNote_Save tests that saving a Zettel note writes its content to disk.
func TestZettelNote_Save(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	zNote, err := zettel.NewZettelNote("SaveNote", cfg, dtm, dl, dfs,
		note.WithContent("Content to Save"),
	)
	require.NoError(t, err)

	// Save the note. Using DummyFS (which wraps OS calls) will create the file.
	err = zNote.Save()
	require.NoError(t, err)

	// Verify that the file exists and has the expected content.
	content, err := os.ReadFile(zNote.Path())
	require.NoError(t, err)
	assert.Equal(t, "Content to Save", string(content))
}

// TestZettelNote_String tests that the String method returns a string containing
// the title and tag.
func TestZettelNote_String(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	zNote, err := zettel.NewZettelNote("StringNote", cfg, dtm, dl, dfs)
	require.NoError(t, err)

	str := zNote.String()
	assert.Contains(t, str, "StringNote")
}

// TestZettelNote_Timestamps ensures that the created and modified timestamps are set appropriately.
func TestZettelNote_Timestamps(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	start := time.Now()
	zNote, err := zettel.NewZettelNote("TimeNote", cfg, dtm, dl, dfs)
	require.NoError(t, err)

	// Check that the created and modified times are within one second of the note creation.
	assert.WithinDuration(t, start, zNote.Created(), time.Second)
	assert.WithinDuration(t, start, zNote.Modified(), time.Second)
}
