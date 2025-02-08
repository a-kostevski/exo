package note_test

import (
	"path/filepath"
	"testing"

	"github.com/a-kostevski/exo/pkg/note"
	"github.com/a-kostevski/exo/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseNoteFactory_CreateNote_Success(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)

	factory := note.NewBaseNoteFactory(note.NoteType("test"), cfg, dtm, dl, dfs)
	n, err := factory.CreateNote("Factory Note",
		note.WithSubDir("factory"),
		note.WithFileName("factory_note.md"),
		note.WithTemplateName("default"),
		note.WithContent("Factory Content"),
	)
	require.NoError(t, err)
	require.NotNil(t, n)

	expectedPath := filepath.Join(cfg.Dir.DataHome, "factory", "factory_note.md")
	assert.Equal(t, expectedPath, n.Path())
	assert.Equal(t, "Factory Note", n.Title())
	assert.Equal(t, "Factory Content", n.Content())
}

func TestBaseNoteFactory_CreateNote_Failure(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)
	factory := note.NewBaseNoteFactory(note.NoteType("test"), cfg, dtm, dl, dfs)

	// Missing mandatory options.
	_, err := factory.CreateNote("Incomplete Note")
	require.Error(t, err)
}

func TestBaseNoteFactory_NoteType(t *testing.T) {
	tmpDir := t.TempDir()
	cfg, dtm, dl, dfs, _ := testutil.NewDummyDeps(tmpDir)
	noteType := note.NoteType("custom")
	factory := note.NewBaseNoteFactory(noteType, cfg, dtm, dl, dfs)
	assert.Equal(t, noteType, factory.NoteType())
}
